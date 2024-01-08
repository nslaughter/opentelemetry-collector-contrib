// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuremonitorreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver"

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/azquery"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/azuresdk"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/metadata"
)

var (
	timeGrains = map[string]int64{
		"PT1M":  60,
		"PT5M":  300,
		"PT15M": 900,
		"PT30M": 1800,
		"PT1H":  3600,
		"PT6H":  21600,
		"PT12H": 43200,
		"P1D":   86400,
	}
)

type void struct{}

type azureResourceData struct {
	resource 			*armresources.GenericResourceExpanded
	metricDefinitions	map[string]*azquery.MetricDefinition
	// TODO: could we refactor to make these getters, so we don't need to precompute?
	attributes 			map[string]*string
	tags       			map[string]*string
}

func newScraper(conf *Config, settings receiver.CreateSettings) *azureScraper {
	return &azureScraper{
		cfg:                             conf,
		settings:                        settings.TelemetrySettings,
		mb:                              metadata.NewMetricsBuilder(conf.MetricsBuilderConfig, settings),
		azIDCredentialsFunc:             azidentity.NewClientSecretCredential,
		azIDWorkloadFunc:                azidentity.NewWorkloadIdentityCredential,
		armClientFunc:                   armresources.NewClient,

		// resources:						 map[string]azureResource,
		resources:						 sync.Map
		mutex:                           &sync.Mutex{},
	}
}

type ResourceCache[T comparable, U any] struct {
	store sync.Map
}

func (rc *ResourceCache) Get(key T) (U, bool) {
	return rc.store.Load(key)
}

func (rc *ResourceCache) Set(key T, value U) {
	rc.store.Store(key, value)
}

func (rc *ResourceCache) Delete(key T) {
	rc.store.Delete(key)
}

func (rc *ResourceCache) Clean(keys []T) {
	for _, key := range keys {
		rc.Delete(key)
	}
}

func (rc *ResourceCache) Range(f func(key T, value U) bool) {
	rc.store.Range(func(key, value interface{}) bool {
		return f(key.(T), value.(U))
	})
}

func (rc *ResourceCache) KeySet() map[string]void {
	keyset := map[string]void{}
	for key := range rc.store {
		keyset[key.(string)] = void{}
	}
	return keyset
}


type azureScraper struct {
	cred azcore.TokenCredential

	// clientResources          ArmClient
	// clientMetricsDefinitions MetricsDefinitionsClientInterface
	// clientMetricsValues      MetricsValuesClient

	cfg                             *Config
	settings                        component.TelemetrySettings
	resourcesUpdated                time.Time
	mb                              *metadata.MetricsBuilder
	azIDCredentialsFunc             func(string, string, string, *azidentity.ClientSecretCredentialOptions) (*azidentity.ClientSecretCredential, error)
	azIDWorkloadFunc                func(options *azidentity.WorkloadIdentityCredentialOptions) (*azidentity.WorkloadIdentityCredential, error)
	armClientOptions                *arm.ClientOptions
	armClientFunc                   func(string, azcore.TokenCredential, *arm.ClientOptions) (*armresources.Client, error)
	// armMonitorDefinitionsClientFunc func(string, azcore.TokenCredential, *arm.ClientOptions) (*armmonitor.MetricDefinitionsClient, error)
	metricsClient 					*azquery.MetricsClient
	resourcesClientFactory 			*armresources.ClientFactory
	resourcesClient 				*armresources.Client
	resourceGroupsClient 			*armresources.ResourceGroupsClient

	// resources 						map[string]azureResourceData
	resources 						ResourceCache

	rwMu                            *sync.RWMutex
}

// TODO: see what the race detector wants to protect the resources map

// resourceUpdater is a goroutine that periodically checks the TTL of resource metadata and updates it when necessary.
func (s *azureScraper) runMetricMetadataUpdates(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.CacheResources)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.updateMetricMetadata(ctx)
			}
		}
	}()
}

func (s *azureScraper) updateResourceCache(ctx context.Context) {
	// fetch definitions
	currentKeys := s.resources.KeySet()
	rs, err := s.fetchResources(ctx)
	if err != nil {
		s.settings.Logger.Error("failed to get Azure Resources data", zap.Error(err))
		return
	}
	for _, r := range rs {
		defs, err = s.fetchDefinitions(ctx, rs)
		if err != nil {
			s.settings.Logger.Error("failed to get Azure Metrics definitions data", zap.Error(err))
			return
		}
		s.resources.Set(*r.ID, &azureResourceData{
			resource: r,
			metricDefinitions: defs,
			attributes: azuresdk.CalculateResourceAttributes(r),
			tags: r.Tags,
		})
		delete(currentKeys, *r.ID)
	}

	// evicts resources that weren't found in update
	s.resources.Clean(currentKeys)
}

func (s *azureScraper) fetchResources(ctx context.Context) ([]*armresources.GenericResourceExpanded, error) {
	filter := azuresdk.CalculateResourcesFilter(s.cfg.Services, s.cfg.ResourceGroups)
	opts := &armresources.ClientListOptions{
		Filter: to.Ptr(filter),
	}
	var results []*armresources.GenericResourceExpanded
	err := azuresdk.ProcessPager(
		ctx,
		s.resourcesClient.NewListPager(opts),
		// TODO: use our Extractor in the AzureSDK package
		func(page armresources.ClientListResponse) []*armresources.GenericResourceExpanded {
			return page.Value
		},
		func(ctx context.Context, resource *armresources.GenericResourceExpanded) error {
			results = append(results, resource)
			return nil
		},
	)
	return nil, err
}

func (s *azureScraper) fetchDefinitions(ctx context.Context, rid string) (map[string]*azquery.MetricDefinition, error) {
	var defs map[string]*azquery.MetricDefinition

	err := azuresdk.ProcessPager(
		ctx,
		s.clientMetricsDefinitions.NewListPager(rid),
		func(page azquery.MetricsClientListDefinitionsResponse) []*azquery.MetricDefinition {
			return page.Value
		},
		func(ctx context.Context, v *azquery.MetricDefinition) error {
			defs[*v.ID] = v
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return defs, nil
}

func (s *azureScraper) start(_ context.Context, _ component.Host) (err error) {
	if err = s.loadCredentials(); err != nil {
		return err
	}
	s.updateMetricMetadata(ctx)
	s.runMetricMetadataUpdates(ctx)

	return
}

func (s *azureScraper) loadCredentials() (err error) {
	switch s.cfg.Authentication {
	case servicePrincipal:
		if s.cred, err = s.azIDCredentialsFunc(s.cfg.TenantID, s.cfg.ClientID, s.cfg.ClientSecret, nil); err != nil {
			return err
		}
	case workloadIdentity:
		if s.cred, err = s.azIDWorkloadFunc(nil); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown authentication %v", s.cfg.Authentication)
	}
	return nil
}
/*
func (s *azureScraper) getResources(ctx context.Context, opt) (map[string]*armresources.GenericResourceExpanded, error) {
	pager := s.clientResources.NewListPager(nil)
	var resources []*armresources.GenericResourceExpanded
	if err := azuresdk.ProcessPager(
		ctx,
		pager,
		func(page armresources.ClientListResponse) []*armresources.GenericResourceExpanded {
			return page.Value
		},
		func(ctx context.Context, resource *armresources.GenericResourceExpanded) error {
			resources = append(resources, resource)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return resources, nil
}
*/
func (s *azureScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	// get are own copy of resources
	rs := copy(s.resources)
	var wg sync.WaitGroup
	// TODO: use a one-writer pattern to fan out the work fan in the results
	opts := &armresources.ClientListOptions{
		Filter: to.Ptr(filter),
	}
	for k, v := range rs {
		wg.Add(1)
		go func(k string, v azureResourceData) {
			defer wg.Done()
			azuresdk.ProcessPager(
				ctx,
				s.metricsClient.NewListPager(k, opts),
				func(page azquery.MetricsClientListResponse) []*azquery.Metric {
					return page.Value
				},
				func(ctx context.Context, metric *azquery.Metric) error {
					s.processMetric(ctx, metric)
					return nil
				},
			)
		}(k, v)
	}
}

func (s *azureScraper) processMetric(ctx context.Context, m *azquery.Metric) {
	for _, timeseriesElement := range m.Timeseries {
		if timeseriesElement.Data != nil {
			attributes := map[string]*string{}
			for name, value := range res.attributes {
				attributes[name] = value
			}
			for _, value := range timeseriesElement.Metadatavalues {
				name := metadataPrefix + *value.Name.Value
				attributes[name] = value.Value
			}
			if s.cfg.AppendTagsAsAttributes {
				for tagName, value := range res.tags {
					name := tagPrefix + tagName
					attributes[name] = value
				}
			}
			for _, metricValue := range timeseriesElement.Data {
				s.processTimeseriesData(resourceID, metric, metricValue, attributes)
			}
		}
	}
}

// NOTE: instead of paging all resources in a subscription using a filter
// we could use the ResourcesGroupsClient to get all ResourceGroups and then
// use that client to get all resources in a ResourceGroup.
/*
func (s *azureScraper) updateResources(ctx context.Context) {
	existingResources := map[string]void{}
	for id := range s.resources {
		existingResources[id] = void{}
	}

	filter := azuresdk.CalculateResourcesFilter(s.cfg.Services, s.cfg.ResourceGroups)
	opts := &armresources.ClientListOptions{
		Filter: to.Ptr(filter),
	}

	// pager := s.clientResources.NewListPager(opts)

	if err := azuresdk.ProcessPager(
		ctx,
		s.resourcesClient.NewListPager(opts),
		func(page armresources.ClientListResponse) []*armresources.GenericResourceExpanded {
			return page.Value
		},
		func(ctx context.Context, resource *armresources.GenericResourceExpanded) error {
			// rid fields: Parent, SubscriptionID, ResourceGroupName, Provider, Location, ResourceType, Name
			attributes, err := azuresdk.CalculateResourceAttributes(resource)
			if err != nil {
				s.settings.Logger.Error("failed to calculate resource attributes", zap.Error(err))
			}
			s.resources[*resource.ID] = &azureResource{
				attributes: attributes,
				tags:       resource.Tags,
			}
			delete(existingResources, *resource.ID)
			return nil
		},
	); err != nil {
		s.settings.Logger.Error("failed to get Azure Resources data", zap.Error(err))
		return
	}
	if len(existingResources) > 0 {
		for idToDelete := range existingResources {
			delete(s.resources, idToDelete)
		}
	}

	s.resourcesUpdated = time.Now()
}

func (s *azureScraper) updateMetricsDefinitions(ctx context.Context, resourceID string) {
	var defs []*azquery.MetricDefinition
	pager := s.clientMetricsDefinitions.NewListPager(resourceID, nil)
	err := azuresdk.ProcessPager(
		ctx,
		pager,
		func(page azquery.MetricsClientListDefinitionsResponse) []*azquery.MetricDefinition {
			return page.Value
		},
		func(ctx context.Context, v *azquery.MetricDefinition) error {
			defs = append(defs, v)
			// s.metricDefinitions[*v.ID] = v
			return nil
		},
	)
	s.metricDefinitions[resourceID] = defs
}
*/
func (s *azureScraper) processResourceMetricsValues(ctx context.Context, resourceID string) {
	for _, resources := range s.resources {
		for _, resource := range resources {
			s.getResourceMetricsValues(ctx, resourceID)
		}
	}
	res := *s.resources[resourceID]

	for compositeKey, metricsByGrain := range res.metricsByCompositeKey {
		if time.Since(metricsByGrain.metricsValuesUpdated).Seconds() < float64(timeGrains[compositeKey.timeGrain]) {
			continue
		}
		// metricsByGrain.metricsValuesUpdated = time.Now()
		start := 0

		for start < len(metricsByGrain.metrics) {

			end := start + s.cfg.MaximumNumberOfMetricsInACall
			if end > len(metricsByGrain.metrics) {
				end = len(metricsByGrain.metrics)
			}

			opts := getResourceMetricsValuesRequestOptions(
				metricsByGrain.metrics,
				compositeKey.dimensions,
				compositeKey.timeGrain,
				start,
				end,
			)
			start = end

			result, err := s.clientMetricsValues.List(
				ctx,
				resourceID,
				&opts,
			)
			if err != nil {
				s.settings.Logger.Error("failed to get Azure Metrics values data", zap.Error(err))
				return
			}

			err := s.

			for _, metric := range result.Value {

				for _, timeseriesElement := range metric.Timeseries {

					if timeseriesElement.Data != nil {
						attributes := map[string]*string{}
						for name, value := range res.attributes {
							attributes[name] = value
						}
						for _, value := range timeseriesElement.Metadatavalues {
							name := metadataPrefix + *value.Name.Value
							attributes[name] = value.Value
						}
						if s.cfg.AppendTagsAsAttributes {
							for tagName, value := range res.tags {
								name := tagPrefix + tagName
								attributes[name] = value
							}
						}
						for _, metricValue := range timeseriesElement.Data {
							s.processTimeseriesData(resourceID, metric, metricValue, attributes)
						}
					}
				}
			}
		}
	}
}

// NOTE: we're dropping this idea and could re-implement with something simpler if we
// enable user config to specify dimensions. We no longer need this to group by dimensions
// for batching fetches since azquery methods let us fetch all metrics for a resource in
// a single call.
func getResourceMetricsValuesRequestOptions(
	metrics []string,
	dimensionsStr string,
	timeGrain string,
	start int,
	end int,
) armmonitor.MetricsClientListOptions {
	resType := strings.Join(metrics[start:end], ",")
	filter := armmonitor.MetricsClientListOptions{
		Metricnames: &resType,
		Interval:    to.Ptr(timeGrain),
		// TODO: check this because it looks odd since .Timespan is a timeGrain here
		// but SDK doc says it's expressed with the format 'startDateTimeISO/endDateTimeISO'.
		Timespan:    to.Ptr(timeGrain),
		Aggregation: to.Ptr(strings.Join(aggregations, ",")),
	}

	// NOTE: we're dropping this idea and could re-implement with something simpler if we
	// enable user config to specify dimensions. We no longer need this to group by dimensions
	// for batching fetches since azquery methods let us fetch all metrics for a resource in
	// a single call.
	if len(dimensionsStr) > 0 {
		var dimensionsFilter bytes.Buffer
		dimensions := strings.Split(dimensionsStr, ",")
		for i, dimension := range dimensions {
			dimensionsFilter.WriteString(dimension)
			dimensionsFilter.WriteString(" eq '*' ")
			if i < len(dimensions)-1 {
				dimensionsFilter.WriteString(" and ")
			}
		}
		dimensionFilterString := dimensionsFilter.String()
		filter.Filter = &dimensionFilterString
	}

	return filter
}

func (s *azureScraper) processTimeseriesData(
	resourceID string,
	metric *armmonitor.Metric,
	metricValue *armmonitor.MetricValue,
	attributes map[string]*string,
) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// NOTE: we need to check whether the AzureSDK is providing a timestamp
	// and prefer that one whenever possible. If anything we should be checking
	// if the AzureSDK timestamp is valid - maybe a nil check combined with
	// a time range check
	ts := pcommon.NewTimestampFromTime(time.Now())

	aggregations := azuresdk.ExtractAggregations(metricValue)

	for _, aggregation := range aggregationsData {
		if aggregation.value != nil {
			// TODO: workout a more OTel idiom for AddDataPoint
			s.mb.AddDataPoint(
				resourceID,
				*metric.Name.Value,
				aggregation.name,
				string(*metric.Unit),
				attributes,
				// metric.Timestamp,
				ts,
				*aggregation.value,
			)
		}
	}
}

/*
type azuerProvider struct {
	clientFactory 			*armresources.ClientFactory
	resourceGroupsClients	*armresources.ResourceGroupsClient
	resourcesClient 		*armresources.Client
	metricsClient			*azquery.MetricsClient
}

func cloud(cfg *Config) cloud.Configuration {
	cloudToUse := cloud.AzurePublic
	if cfg.Cloud == azureGovernmentCloud {
		cloudToUse = cloud.AzureGovernment
	}
	return cloudToUse
}
func newResourcesProvider(ctx context.Context, cfg *Config) (*resourcesProvider, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	opts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud(cfg)
		},
	}

	clientFactory := armresources.NewClientFactory(cfg.SubscriptionID, cred, opts)
	resourceGroupsClient := clientFactory.NewResourceGroupsClient(cfg.SubscriptionID)
	resourcesClient := clientFactory.NewClient(cfg.SubscriptionID)
	metricsClient := azquery.NewMetricsClient(cfg.SubscriptionID)

	return &resourcesProvider{
		clientFactory: clientFactory,
		resourceGroupsClients: resourceGroupsClient,
		resourcesClient: resourcesClient,
		metricsClient: metricsClient,
	}, nil
}

func (rp *resourcesProvider) getResourceGroups(ctx context.Context) ([]*armresources.ResourceGroup, error) {
	pager := rp.resourceGroupsClient.NewListPager(nil)
	var resourceGroups []*armresources.ResourceGroup
	if err := azuresdk.ProcessPager(
		ctx,
		pager,
		func(page armresources.ResourceGroupsClientListResponse) []*armresources.ResourceGroup {
			return page.Value
		},
		func(ctx context.Context, rg *armresources.ResourceGroup) error {
			resourceGroups = append(resourceGroups, rg)
			return nil
		},
	) err != nil {
		return nil, err
	}
	return resourceGroups, nil
}

func (rp *resourcesProvider) getResourceGroupsAsync(ctx context.Context, opts *armresources.ResourceGroupsClientListOptions) ([]*armresources.ResourceGroup, error) {
	pager := rp.resourceGroupsClient.NewListPager(nil)
	c := make(chan *armresources.ResourceGroup)
	if err := azuresdk.ProcessPager(
		ctx,
		pager,
		func(page armresources.ResourceGroupsClientListResponse) []*armresources.ResourceGroup {
			return page.Value
		},
		func(ctx context.Context, rg *armresources.ResourceGroup) error {
			c <- rg
			return nil
		},
	) err != nil {
		return nil, err
	}
	return c, nil
}



func (rp *resourcesProvider) getResources(ctx context.Context, opts *armresources.ClientListOptions) ([]*armresources.GenericResourceExpanded, error) {
	pager := rp.resourcesClient.NewListPager(opts)
	var resources []*armresources.GenericResourceExpanded
	if err := azuresdk.ProcessPager(
		ctx,
		pager,
		func(page armresources.ClientListResponse) []*armresources.GenericResourceExpanded {
			return page.Value
		},
		func(ctx context.Context, resource *armresources.GenericResourceExpanded) error {
			resources = append(resources, resource)
			return nil
		},
	) err != nil {
		return nil, err
	}
	return resources, nil
}

func (rp *resourcesProvider) getResourcesAsync(ctx context.Context, opts *armresources.ClientListOptions) (chan<- *armresources.GenericResourceExpanded, error) {
	pager := rp.resourcesClient.NewListPager(opts)
	c := make(chan *armresources.GenericResourceExpanded)
	if err := azuresdk.ProcessPager(
		ctx,
		pager,
		func(page armresources.ClientListResponse) []*armresources.GenericResourceExpanded {
			return page.Value
		},
		func(ctx context.Context, resource *armresources.GenericResourceExpanded) error {
			c <- resource
			return nil
		},
	) err != nil {
		return nil, err
	}
	return c, nil
}

type MonitorClient interface {
	GetResourceGroups(ctx context.Context) ([]*armresources.ResourceGroup, error)
	GetResources(ctx context.Context) ([]*armresources.GenericResourceExpanded, error)
	GetMetricsDefinitions(ctx context.Context, resourceID string) ([]*azquery.MetricDefinition, error)
	GetMetrics(ctx context.Context, resourceID string) ([]*azquery.Metric, error)
}
*/
