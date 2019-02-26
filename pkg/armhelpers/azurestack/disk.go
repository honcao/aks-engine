// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package azurestack

import (
	"context"

	"github.com/Azure/aks-engine/pkg/armhelpers"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
	azcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
)

// DeleteManagedDisk deletes a managed disk.
func (az *AzureClient) DeleteManagedDisk(ctx context.Context, resourceGroupName string, diskName string) error {
	future, err := az.disksClient.Delete(ctx, resourceGroupName, diskName)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.disksClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.disksClient)
	return err
}

// DiskListPageClient contains a page of Disk values.
type DiskListPageClient struct {
	dlp compute.DiskListPage
	err error
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *DiskListPageClient) NextWithContext(ctx context.Context) (err error) {
	return page.dlp.NextWithContext(ctx)
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *DiskListPageClient) Next() error {
	return page.dlp.Next()
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page DiskListPageClient) NotDone() bool {
	return !page.dlp.NotDone()
}

// Response returns the raw server response from the last page request.
func (page DiskListPageClient) Response() azcompute.DiskList {
	l := azcompute.DiskList{}
	DeepAssignment(&l, page.dlp.Response())
	return l
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page DiskListPageClient) Values() []azcompute.Disk {
	l := []azcompute.Disk{}
	DeepAssignment(l, page.dlp.Values())
	return l
}

// ListManagedDisksByResourceGroup lists managed disks in a resource group.
func (az *AzureClient) ListManagedDisksByResourceGroup(ctx context.Context, resourceGroupName string) (result armhelpers.DiskListPage, err error) {
	page, err := az.disksClient.ListByResourceGroup(ctx, resourceGroupName)
	return &DiskListPageClient{
		dlp: page,
		err: err,
	}, err
}
