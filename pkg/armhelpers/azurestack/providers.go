// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package azurestack

import (
	"context"

	"github.com/Azure/aks-engine/pkg/armhelpers"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	azresources "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"

	"github.com/Azure/go-autorest/autorest/to"
)

// ProviderListResultPageClient contains a page of Provider values.
type ProviderListResultPageClient struct {
	plrp resources.ProviderListResultPage
	err  error
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *ProviderListResultPageClient) NextWithContext(ctx context.Context) (err error) {
	return page.plrp.NextWithContext(ctx)
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *ProviderListResultPageClient) Next() error {
	return page.plrp.Next()
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page ProviderListResultPageClient) NotDone() bool {
	return page.plrp.NotDone()
}

// Response returns the raw server response from the last page request.
func (page ProviderListResultPageClient) Response() azresources.ProviderListResult {
	r := azresources.ProviderListResult{}
	DeepAssignment(&r, page.plrp.Response())
	return r
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page ProviderListResultPageClient) Values() []azresources.Provider {
	r := []azresources.Provider{}
	DeepAssignment(&r, page.plrp.Values())
	return r
}

// ListProviders returns all the providers for a given AzureClient
func (az *AzureClient) ListProviders(ctx context.Context) (armhelpers.ProviderListResultPage, error) {
	page, err := az.providersClient.List(ctx, to.Int32Ptr(100), "")
	return &ProviderListResultPageClient{
		plrp: page,
		err:  err,
	}, err
}
