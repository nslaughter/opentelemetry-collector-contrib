// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuresdk // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/azuresdk"

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

type pipeline struct {
	ctx context.Context
	resourcesClient *armresources.Client
	resourceGroupsClient *armresources.ResourceGroupsClient
	metricsClient *azquery.MetricsClient
	inC <-chan string
	errC chan error
}

func newPipeline(
	ctx context.Context,
	rc *armresources.Client,
	rgc *armresources.ResourceGroupsClient,
	mc *azquery.MetricsClient,
	inC <-chan string,
	errC chan error) *pipeline {
	return &pipeline{
		ctx: ctx,
		resourcesClient: rc,
		resourceGroupsClient: rgc,
		metricsClient: mc,
		inC: inC,
		errC: errC,
	}
}

func (p *pipeline) handleErrors() {
	go func() {
		for err := range p.errC {
			log.Println(err)
		}
	}()
}

func (p *pipeline) run() {
	// The pipeline proceeds as follows:
	//
	// send subscriptionID: 	single value producer

	// resourceGroupsStage 		->

	// resourcesStage 			->

	// metricDefinitionsStage	->

	// metricValuesStage		->

	// processMetricsStage		-> sink

}


// see example in azure
// https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/fb2b81cab55cdd6bb2c0f10ab47ae8e0a15a5dfc/sdk/resourcemanager/resource/resourcegroups/main.go#L115-L128
// func (s *azureScraper) listResourceGroups(ctx context.Context) ([]*armresources.ResourceGroup, error) {


// Generic stage takes type parameters of:
// - ctx
// - input of any (recv chan)
// - opts: Options for pager factory
// - pagerFactory: func(item IN, options OPTS) *runtime.Pager[RESP]
// - pageHandler: func(page RESP, output chan<- OUT)

// pagerFactory PagerFactory[IN, OPTS, RESP],
// pageHandler func(page RESP, outC chan<- OUT),

// genericStage's type is parameterized by IN, OUT, RESP, OPTS
// - IN is the type of the input channel
// - OUT is the type of the output channel
// - RESP is the response type of the Azure SDK pager
// - OPTS is the type of the Azure SDK options type for this pager
func genericStage[IN any, OUT any, RESP any, OPTS any](
	ctx context.Context,
	inC <-chan IN,
	opts OPTS,
	pagerFactory PagerFactory[IN, OPTS, RESP], // func(item IN, options OPTS) *runtime.Pager[RESP],
	pageHandler func(page RESP, outC chan<- OUT), // page handler needs typing from function signature to use .Value field
	) (<-chan OUT, chan error) {
	outC := make(chan OUT)
	errC := make(chan error, 1)
	go func() {
		defer close(outC)
		defer close(errC)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				select {
				case <-ctx.Done():
					return
				case item, ok := <-inC:
					if !ok {
						return
					}
					pager := pagerFactory.CreatePager(item, opts)
					for pager.More() {
						page, err := pager.NextPage(ctx)
						if err != nil {
							errC <- err
							return
						}
						pageHandler(page, outC)
					}
				}
			}
		}
	}()
	return outC, errC
}

/*
func (p *pipeline) resourceGroupsStage(
    ctx context.Context,
    opts armresources.ResourceGroupsClientListOptions,
	inC <-chan string,
	errC chan<- error,
) (<-chan *armresources.ResourceGroup, chan<- error) {
	outC := make(chan *armresources.ResourceGroup)
	defer close(outC)

	pagerFactory := PagerFactoryWithOptions[string, armresources.ResourceGroupsClientListOptions, armresources.ResourceGroupsClientListResponse]{
    	fn: func(_ armresources.ResourceGroupsClientListOptions) *runtime.Pager[armresources.ResourceGroupsClientListResponse] {
			return p.resourceGroupsClient.NewListPager(nil)
		},
	}

	handler := func(page armresources.ResourceGroupsClientListResponse, outC chan<- *armresources.ResourceGroup) {
		for _, rg := range page.Value {
			outC <- rg
		}
	}

	return genericStage[string, *armresources.ResourceGroup, armresources.ResourceGroupsClientListResponse, armresources.ResourceGroupsClientListOptions](
		ctx,
		inC,
		opts,
		pagerFactory,
		handler,
	)
}
*/
func (s *pipeline) resourcesStage(
    ctx context.Context,
    opts *armresources.ClientListByResourceGroupOptions,
	inC <-chan armresources.ResourceGroup,
	errC <-chan error,
) (<-chan *armresources.GenericResourceExpanded, <-chan error) {
	outC := make(chan *armresources.GenericResourceExpanded)
	defer close(outC)

	pagerFactory := PagerFactoryWithItem[armresources.ResourceGroup, *armresources.ClientListByResourceGroupOptions, armresources.ClientListByResourceGroupResponse]{
		fn: func(rg armresources.ResourceGroup, opts *armresources.ClientListByResourceGroupOptions) *runtime.Pager[armresources.ClientListByResourceGroupResponse] {
			return s.resourcesClient.NewListByResourceGroupPager(*rg.Name, opts)
		},
	}

	handler := func(page armresources.ClientListByResourceGroupResponse, outC chan<- *armresources.GenericResourceExpanded) {
		for _, r := range page.Value {
			outC <- r
		}
	}

	return genericStage[armresources.ResourceGroup, *armresources.GenericResourceExpanded, armresources.ClientListByResourceGroupResponse, *armresources.ClientListByResourceGroupOptions](
		ctx,
		inC,
		opts,
		pagerFactory,
		handler,
	)
}
/*
func (s *pipeline) metricDefinitionsStage(
	ctx context.Context,
	opts *azquery.MetricsClientListDefinitionsOptions,
	inC <-chan *azquery.MetricsClientListDefinitionsResponse,
	errC <-chan error,
) (<-chan *azquery.MetricDefinition, <-chan error) {
	outC := make(chan *azquery.MetricDefinition)
	defer close(outC)

	pagerFactory := PagerFactoryWithItem[string, *azquery.MetricsClientListDefinitionsOptions, azquery.MetricsClientListDefinitionsResponse]{
		fn: func(resourceID string, opts *azquery.MetricsClientListDefinitionsOptions) *runtime.Pager[azquery.MetricsClientListDefinitionsResponse] {
			return s.metricsClient.NewListDefinitionsPager(resourceID, opts)
		},
	}

	handler := func(page azquery.MetricsClientListDefinitionsResponse, outC chan<- *azquery.MetricDefinition) {
		for _, md := range page.Value {
			outC <- md
		}
	}

	return genericStage[string, *azquery.MetricDefinition, azquery.MetricsClientListDefinitionsResponse, *azquery.MetricsClientListDefinitionsOptions](
		ctx,
		inC,
		opts,
		pagerFactory,
		handler,
	)
}

// QueryResource is the final fetch. Unlike the other fetches, this one is not paged.
func (s *pipeline) metricValuesStage(ctx context.Context, inC <-chan string, errC chan<- error) (<-chan *azquery.Metric, <-chan error) {
	outC := make(chan *azquery.MetricValue)
	defer close(outC)

	for mv, ok := range inC {
		if !ok {
			return
		}

	}

	// TODO:
	resp, err := s.metricsClient.QueryResource(ctx, res, nil)
	if err != nil {
		errC <- err
		return outC, errC
	}
	for _, m := range resp.Value {
		outC <- m
	}

	return outC, errC
}
*/
func (p *pipeline) run() {
	// The pipeline proceeds as follows:
	//
	// send subscriptionID: 	single value producer

	// resourceGroupsStage 		->

	// resourcesStage 			->

	// metricDefinitionsStage	->

	// metricValuesStage		->

	// processMetricsStage		-> sink

}

// We need to handle two signatures for the pager factory. So we're using the adapter pattern
// to create a common interface for both signatures. And we'll type assert the interface in the
// generic stage function to call the correct signature. See https://golangbyexample.com/adapter-design-pattern-go/
type PagerFactory[IN any, OPTS any, RESP any] interface {
    CreatePager(item IN, opts OPTS) *runtime.Pager[RESP]
    CreatePagerWithOptions(opts OPTS) *runtime.Pager[RESP]
}

type PagerFactoryWithItem[IN any, OPTS any, RESP any] struct {
    fn func(item IN, opts OPTS) *runtime.Pager[RESP]
}

func (p PagerFactoryWithItem[IN, OPTS, RESP]) CreatePager(item IN, opts OPTS) *runtime.Pager[RESP] {
    return p.fn(item, opts)
}

func (p PagerFactoryWithItem[IN, OPTS, RESP]) CreatePagerWithOptions(opts OPTS) *runtime.Pager[RESP] {
    // This function should not be called for this type
    panic("CreatePagerWithOptions called on PagerFactoryWithItem")
}

type PagerFactoryWithOptions[IN any, OPTS any, RESP any] struct {
    fn func(opts OPTS) *runtime.Pager[RESP]
}

func (p PagerFactoryWithOptions[IN, OPTS, RESP]) CreatePager(item IN, opts OPTS) *runtime.Pager[RESP] {
    // This function should not be called for this type
    panic("CreatePager called on PagerFactoryWithOptions")
}

func (p PagerFactoryWithOptions[IN, OPTS, RESP]) CreatePagerWithOptions(opts OPTS) *runtime.Pager[RESP] {
    return p.fn(opts)
}
