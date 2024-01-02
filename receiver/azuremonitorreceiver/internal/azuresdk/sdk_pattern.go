// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuresdk // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/azuresdk"

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"

	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

// See the SDK documentation on Common Usage Patterns at
// https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-core-concepts#policy-template

type Page[T any] struct {
    Value []T
}

type ProcessorFunc[T any] func(T) error

type ResourceGroupsPage = Page[*armresources.ResourceGroup]
type ResourcesPage = Page[*armresources.GenericResourceExpanded]
type MetricsValuesPage = Page[*azquery.Metric]
type MetricDefinitionsPage = Page[*azquery.MetricDefinition]

// Generic function to process paginated results.
func ProcessPage[T any](page Page[T], process ProcessorFunc[T]) error {
    for _, item := range page.Value {
        if err := process(item); err != nil {
            return err
        }
    }
    return nil
}

// Generic function to process paginated results.
func ProcessPager[T any](ctx context.Context, pager *runtime.Pager[T], process ProcessorFunc[T]) error {
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
		 	return err
		}
		if err := process(page); err != nil {
			return err
		}
	}
	return nil
}

func processResourceGroups(ctx context.Context, pager *runtime.Pager[armresources.ResourceGroupsClientListResponse]) error {
	processPage := func(page armresources.ResourceGroupsClientListResponse) error {
		return ProcessPage(page, func(resourceGroup *armresources.ResourceGroup) error {
			// This is where we would send the resource group to the next stage.
			fmt.Println(*resourceGroup.Name)
			return nil
		})
	}
 	return ProcessPager(ctx, pager, processPage)
}

func processResources(ctx context.Context, pager *runtime.Pager[armresources.ClientListResponse]) error {
	processPage := func(page armresources.ClientListResponse) error {
		return ProcessPage(page, func(resource *armresources.GenericResourceExpanded) error {
		  	// This is where we would send the resource to the next stage.
			fmt.Println(*resource.Name)
			return nil
		})
	}
	return ProcessPager(ctx, pager, processPage)
}

func processMetricsValues(ctx context.Context, pager *runtime.Pager[azquery.MetricsClientQueryResourceResponse]) error {
	processPage := func(page azquery.MetricsClientQueryResourceResponse) error {
		return ProcessPage(page, func(metricValue *azquery.Metric) error {
			// This is where we would send the metric value to the next stage.
			fmt.Println(*metricValue.Name)
			return nil
		})
	}
	return ProcessPager(ctx, pager, processPage)
}

func processMetricDefinitions(ctx context.Context, pager *runtime.Pager[azquery.MetricsClientListDefinitionsResponse]) error {
	processPage := func(page azquery.MetricsClientListDefinitionsResponse) error {
		return ProcessPage(page, func(metricDefinition *azquery.MetricDefinition) error {
			// This is where we would send the metric definition to the next stage.
			fmt.Println(*metricDefinition.Name)
			return nil
		})
	}
	return ProcessPager(ctx, pager, processPage)
}

func processResourceGroupsPage(process ProcessorFunc[*armresources.ResourceGroup]) ProcessorFunc[armresources.ResourceGroupsClientListResponse] {
	return func(page armresources.ResourceGroupsClientListResponse) error {
		for _, resourceGroup := range page.Value {
			if err := process(resourceGroup); err != nil {
				return err
			}
		}
		return nil
	}
}

func processResourcesPage(page armresources.ClientListResponse, process ProcessorFunc[*armresources.GenericResourceExpanded]) ProcessorFunc[armresources.ClientListResponse] {
	return func(page armresources.ClientListResponse) error {
		for _, resource := range page.Value {
			if err := process(resource); err != nil {
				return err
			}
		}
		return nil
	}
}

func processMetricsValuesPage(page azquery.MetricsClientQueryResourceResponse, process ProcessorFunc[*azquery.Metric]) ProcessorFunc[azquery.MetricsClientQueryResourceResponse] {
	return func(page azquery.MetricsClientQueryResourceResponse) error {
		for _, metricValue := range page.Value {
			if err := process(metricValue); err != nil {
				return err
			}
		}
		return nil
	}
}

func processMetrics(ctx context.Context, pager *runtime.Pager[azquery.MetricsClientQueryResourceResponse]) error {
    processPage := func(page azquery.MetricsClientQueryResourceResponse) error {
        for _, metric := range page.Value {
            fmt.Println(*metric.Name)
        }
        return nil
    }

    return ProcessPager(ctx, pager, processPage)
}

func processMetricDefinitionsPage(page azquery.MetricsClientListDefinitionsResponse, process ProcessorFunc[*azquery.MetricDefinition]) ProcessorFunc[azquery.MetricsClientListDefinitionsResponse] {
	return func(page azquery.MetricsClientListDefinitionsResponse) error {
		for _, metricDefinition := range page.Value {
			if err := process(metricDefinition); err != nil {
				return err
			}
		}
		return nil
	}
}

func processMetricsPage(page azquery.MetricsClientQueryResourceResponse, process ProcessorFunc[*azquery.Metric]) ProcessorFunc[azquery.MetricsClientQueryResourceResponse] {
	return func(page azquery.MetricsClientQueryResourceResponse) error {
		for _, metricValue := range page.Value {
			if err := process(metricValue); err != nil {
				return err
			}
		}
		return nil
	}
}
