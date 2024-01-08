//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azuresdk_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/azuresdk"
)

func TestMetricsClient(t *testing.T) {
	client, err := azquery.NewMetricsClient(credential, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	c := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "https://...",
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Audience: "",
				Endpoint: "",
			},
		},
	}
	opts := azcore.ClientOptions{Cloud: c}
	cloudClient, err := azquery.NewMetricsClient(credential, &azquery.MetricsClientOptions{ClientOptions: opts})
	require.Error(t, err)
	require.Equal(t, err.Error(), "provided Cloud field is missing Azure Monitor Metrics configuration")
	require.Nil(t, cloudClient)
}

func TestQueryResource_BasicQuerySuccess(t *testing.T) {
	client := startMetricsTest(t)
	timespan := azquery.TimeInterval("PT12H")
	res, err := client.QueryResource(context.Background(), resourceURI,
		&azquery.MetricsClientQueryResourceOptions{
			Timespan:        to.Ptr(timespan),
			Interval:        to.Ptr("PT1M"),
			MetricNames:     nil,
			Aggregation:     to.SliceOfPtrs(azquery.AggregationTypeAverage, azquery.AggregationTypeCount),
			Top:             nil,
			OrderBy:         to.Ptr("Average asc"),
			Filter:          nil,
			ResultType:      nil,
			MetricNamespace: to.Ptr("Microsoft.AppConfiguration/configurationStores"),
		})
	require.NoError(t, err)
	require.NotNil(t, res.Response.Timespan)
	require.Equal(t, *res.Response.Value[0].ErrorCode, "Success")
	require.Equal(t, *res.Response.Namespace, "Microsoft.AppConfiguration/configurationStores")

	// This is how you handle a metric query batch response
	extractor := azuresdk.ExtractMetrics
	process := func(ctx context.Context, v *azquery.Metric) error {
		require.NotNil(t, v)
		for _, timeSeriesElement := range v.TimeSeries {
			require.NotNil(t, timeSeriesElement)
			require.NotNil(t, timeSeriesElement.Data)
			for _, metricValue := range timeSeriesElement.Data {
				t.Log(*metricValue.TimeStamp)
			}
		}
		return nil
	}

	testSerde(t, &res)
	testSerde(t, res.Value[0])
	testSerde(t, res.Value[0].Name)
	testSerde(t, res.Value[0].TimeSeries[0])
}

func TestQueryResource_BasicQueryFailure(t *testing.T) {
	client := startMetricsTest(t)
	invalidResourceURI := "123"
	var httpErr *azcore.ResponseError

	res, err := client.QueryResource(context.Background(), invalidResourceURI, nil)

	require.Error(t, err)
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "MissingSubscription")
	require.Equal(t, httpErr.StatusCode, 404)
	require.Nil(t, res.Timespan)
	require.Nil(t, res.Value)
	require.Nil(t, res.Cost)
	require.Nil(t, res.Interval)
	require.Nil(t, res.Namespace)
	require.Nil(t, res.ResourceRegion)

	testSerde(t, &res)
}

func TestNewListDefinitionsPager_Success(t *testing.T) {
	client := startMetricsTest(t)

	pager := client.NewListDefinitionsPager(resourceURI, nil)
	extractor := azuresdk.ExtractMetricDefinitions
	process := func(ctx context.Context, v *azquery.MetricDefinition) error {
		require.NotNil(t, v)
		return nil
	}
	err := azuresdk.ProcessPager(context.Background(), pager, extractor, process)
	require.NoError(t, err)
}

func TestNewListDefinitionsPager_Failure(t *testing.T) {
	client := startMetricsTest(t)

	pager := client.NewListDefinitionsPager(resourceURI, nil)
	extractor := azuresdk.ExtractMetricDefinitions
	process := func(ctx context.Context, v *azquery.MetricDefinition) error {
		require.NotNil(t, v)
		return nil
	}

	err := azuresdk.ProcessPager(context.Background(), pager, extractor, process)
	require.NoError(t, err)
	// test if first page is valid
	// if pager.More() {
	// 	res, err := pager.NextPage(context.Background())
	// 	require.NoError(t, err)
	// 	require.NotNil(t, res.Value)
	// 	testSerde(t, &res.MetricDefinitionCollection)
	// } else {
	// 	t.Fatal("no response")
	// }
}

func TestNewListNamespacesPager_Success(t *testing.T) {
	client := startMetricsTest(t)

	pager := client.NewListNamespacesPager(resourceURI, &azquery.MetricsClientListNamespacesOptions{})
	extract := func(page azquery.MetricsClientListNamespacesResponse) []*azquery.MetricNamespace {
		return page.Value
	}
	process := func(ctx context.Context, v *azquery.MetricNamespace) error {
		require.NotNil(t, v)
		return nil
	}
	require.NoError(
		t,
		azuresdk.ProcessPager(context.Background(), pager, extract, process),
	)
}

func TestNewListNamespacesPager_Failure(t *testing.T) {
	client := startMetricsTest(t)
	invalidResourceURI := "123"
	var httpErr *azcore.ResponseError

	pager := client.NewListNamespacesPager(invalidResourceURI, nil)
	if pager.More() {
		res, err := pager.NextPage(context.Background())
		require.Error(t, err)
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, httpErr.ErrorCode, "MissingSubscription")
		require.Equal(t, httpErr.StatusCode, 404)
		require.Nil(t, res.Value)
	} else {
		t.Fatal("no response")
	}

}
