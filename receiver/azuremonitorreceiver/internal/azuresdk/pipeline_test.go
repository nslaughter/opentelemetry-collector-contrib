// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuresdk // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/azuresdk"

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func FakePipeline() *pipeline {
	subC := make(chan string)
	errC := make(chan error)
	mc := fakeMetricsClient()
	rc := fakeResourcesClient()
	rgc := fakeResourceGroupClient()
	return newPipeline(
		context.Background(),
		rc,
		rgc,
		mc,
		subC,
		errC,
	)
}

func ExampleServer() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	inC := make(chan string)
	errC := make(chan error)
	p := pipeline{
		ctx: context.Background(),
		resourcesClient: fakeResourcesClient(),
		resourceGroupsClient: fakeResourceGroupClient(),
		metricsClient: fakeMetricsClient(),
		inC: inC,
	}

	h := func(page armresources.ClientListResponse, outC chan<- *armresources.GenericResourceExpanded) {
		for _, resource := range page.ResourceListResult.Value {
			outC <- resource
		}
	}

	// now create the corresponding client, connecting the fake server via the client options
	// call resourcesStage() APOI
	// [armresources.GenericResourceExpanded, *armresources.GenericResourceExpanded, armresources.ClientListResponse, *armresources.ClientListOptions]
	resC, errC := p.resourcesStage(
		context.Background(),
		armresources.ClientListOptions{},
		p,
		h,
	)

	for resource := range resC {
		fmt.Println(*resource.Name)
	}

	// Output:
	// true
	// fake for method Get not implemented
}

func toPtr(s string) *string {
	return &s
}


