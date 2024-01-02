// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuresdk // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/azuresdk"

import (
	// "context"
	// "fmt"
	// "log"
	// "net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"

	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/fake"

	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

func fakePager() azfake.PagerResponder[armresources.ClientListResponse] {
	p := azfake.PagerResponder[armresources.ClientListResponse]{}
	p.AddPage(
		200,
		armresources.ClientListResponse{
			ResourceListResult: armresources.ResourceListResult{
				Value: []*armresources.GenericResourceExpanded{
					{
						ID:   toPtr("fakeID1"),
						Name: toPtr("fakeName1"),
					},
					{
						ID:   toPtr("fakeID2"),
						Name: toPtr("fakeName2"),
					},
				},
			},
		},
		nil,
	)
	return p
}

func fakeResourceGroupClient() *armresources.ResourceGroupsClient {
	fs := FakeServer()
	c, _ := armresources.NewResourceGroupsClient(
		"subscriptionID",
		&azfake.TokenCredential{},
		&arm.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: fake.NewServerTransport(&fs),
			},
		},
	)
	return c
}

func fakeResourcesClient() *armresources.Client {
	fs := FakeServer()
	c, _ := armresources.NewClient(
		"subscriptionID",
		&azfake.TokenCredential{},
		&arm.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: fake.NewServerTransport(&fs),
			},
		},
	)
	return c
}

func fakeMetricsClient() *azquery.MetricsClient {
	fs := FakeServer()
	c, _ := azquery.NewMetricsClient(
		&azfake.TokenCredential{},
		&azquery.MetricsClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: fake.NewServerTransport(&fs),
			},
		},
	)
	return c
}

func FakeServer() fake.Server {
	return fake.Server{
		NewListPager: func(_ *armresources.ClientListOptions) azfake.PagerResponder[armresources.ClientListResponse] {
			return fakePager()
		},
	}
}

