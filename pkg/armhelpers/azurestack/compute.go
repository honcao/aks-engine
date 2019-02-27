// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package azurestack

import (
	"context"

	"github.com/Azure/aks-engine/pkg/armhelpers"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
	azcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
)

//VirtualMachineListResultPageClient Virtual Machine List Result Page Client
type VirtualMachineListResultPageClient struct {
	vmlrp compute.VirtualMachineListResultPage
	err   error
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *VirtualMachineListResultPageClient) NextWithContext(ctx context.Context) (err error) {
	return page.vmlrp.NextWithContext(ctx)
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *VirtualMachineListResultPageClient) Next() error {
	return page.vmlrp.Next()
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page VirtualMachineListResultPageClient) NotDone() bool {
	return page.vmlrp.NotDone()
}

// Response returns the raw server response from the last page request.
func (page VirtualMachineListResultPageClient) Response() azcompute.VirtualMachineListResult {
	r := azcompute.VirtualMachineListResult{}
	DeepAssignment(&r, page.vmlrp.Response())
	return r
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page VirtualMachineListResultPageClient) Values() []azcompute.VirtualMachine {
	r := []azcompute.VirtualMachine{}
	vms := page.vmlrp.Values()
	for _, vm := range vms {
		dst := azcompute.VirtualMachine{}
		DeepAssignment(&dst, vm)
		r = append(r, dst)
	}
	return r
}

// ListVirtualMachines returns (the first page of) the machines in the specified resource group.
func (az *AzureClient) ListVirtualMachines(ctx context.Context, resourceGroup string) (armhelpers.VirtualMachineListResultPage, error) {
	page, err := az.virtualMachinesClient.List(ctx, resourceGroup)
	c := VirtualMachineListResultPageClient{
		vmlrp: page,
		err:   err,
	}
	return &c, err
}

// GetVirtualMachine returns the specified machine in the specified resource group.
func (az *AzureClient) GetVirtualMachine(ctx context.Context, resourceGroup, name string) (azcompute.VirtualMachine, error) {
	vm, err := az.virtualMachinesClient.Get(ctx, resourceGroup, name, "")
	azVM := azcompute.VirtualMachine{}
	DeepAssignment(&azVM, vm)
	return azVM, err
}

// DeleteVirtualMachine handles deletion of a CRP/VMAS VM (aka, not a VMSS VM).
func (az *AzureClient) DeleteVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	future, err := az.virtualMachinesClient.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachinesClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachinesClient)
	return err
}

// VirtualMachineScaleSetListResultPageClient Virtual Machine Scale Set List Result Page Client
type VirtualMachineScaleSetListResultPageClient struct {
	vmsslrp compute.VirtualMachineScaleSetListResultPage
	err     error
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *VirtualMachineScaleSetListResultPageClient) NextWithContext(ctx context.Context) (err error) {
	return page.vmsslrp.NextWithContext(ctx)
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *VirtualMachineScaleSetListResultPageClient) Next() error {
	return page.vmsslrp.Next()
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page VirtualMachineScaleSetListResultPageClient) NotDone() bool {
	return page.vmsslrp.NotDone()
}

// Response returns the raw server response from the last page request.
func (page VirtualMachineScaleSetListResultPageClient) Response() azcompute.VirtualMachineScaleSetListResult {
	r := azcompute.VirtualMachineScaleSetListResult{}
	DeepAssignment(&r, page.vmsslrp.Response())
	return r
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page VirtualMachineScaleSetListResultPageClient) Values() []azcompute.VirtualMachineScaleSet {
	r := []azcompute.VirtualMachineScaleSet{}
	vms := page.vmsslrp.Values()

	for _, vm := range vms {
		dst := azcompute.VirtualMachineScaleSet{}
		DeepAssignment(&dst, vm)
		r = append(r, dst)
	}
	return r
}

// ListVirtualMachineScaleSets returns (the first page of) the vmss resources in the specified resource group.
func (az *AzureClient) ListVirtualMachineScaleSets(ctx context.Context, resourceGroup string) (armhelpers.VirtualMachineScaleSetListResultPage, error) {
	page, err := az.virtualMachineScaleSetsClient.List(ctx, resourceGroup)
	c := VirtualMachineScaleSetListResultPageClient{
		vmsslrp: page,
		err:     err,
	}
	return &c, err
}

// VirtualMachineScaleSetVMListResultPageClient Virtual Machine Scale Set VM List Result Page Client
type VirtualMachineScaleSetVMListResultPageClient struct {
	vmssvlrp compute.VirtualMachineScaleSetVMListResultPage
	err      error
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *VirtualMachineScaleSetVMListResultPageClient) NextWithContext(ctx context.Context) (err error) {
	return page.vmssvlrp.NextWithContext(ctx)
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *VirtualMachineScaleSetVMListResultPageClient) Next() error {
	return page.vmssvlrp.Next()
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page VirtualMachineScaleSetVMListResultPageClient) NotDone() bool {
	return page.vmssvlrp.NotDone()
}

// Response returns the raw server response from the last page request.
func (page VirtualMachineScaleSetVMListResultPageClient) Response() azcompute.VirtualMachineScaleSetVMListResult {
	r := azcompute.VirtualMachineScaleSetVMListResult{}
	DeepAssignment(&r, page.vmssvlrp.Response())
	return r
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page VirtualMachineScaleSetVMListResultPageClient) Values() []azcompute.VirtualMachineScaleSetVM {
	r := []azcompute.VirtualMachineScaleSetVM{}
	vms := page.vmssvlrp.Values()

	for _, vm := range vms {
		dst := azcompute.VirtualMachineScaleSetVM{}
		DeepAssignment(&dst, vm)
		r = append(r, dst)
	}
	return r
}

// ListVirtualMachineScaleSetVMs returns the list of VMs per VMSS
func (az *AzureClient) ListVirtualMachineScaleSetVMs(ctx context.Context, resourceGroup, virtualMachineScaleSet string) (armhelpers.VirtualMachineScaleSetVMListResultPage, error) {
	page, err := az.virtualMachineScaleSetVMsClient.List(ctx, resourceGroup, virtualMachineScaleSet, "", "", "")
	c := VirtualMachineScaleSetVMListResultPageClient{
		vmssvlrp: page,
		err:      err,
	}
	return &c, err
}

// DeleteVirtualMachineScaleSetVM deletes a VM in a VMSS
func (az *AzureClient) DeleteVirtualMachineScaleSetVM(ctx context.Context, resourceGroup, virtualMachineScaleSet, instanceID string) error {
	future, err := az.virtualMachineScaleSetVMsClient.Delete(ctx, resourceGroup, virtualMachineScaleSet, instanceID)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachineScaleSetVMsClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachineScaleSetVMsClient)
	return err
}

// DeleteVirtualMachineScaleSet deletes an entire VM Scale Set.
func (az *AzureClient) DeleteVirtualMachineScaleSet(ctx context.Context, resourceGroup, vmssName string) error {
	future, err := az.virtualMachineScaleSetsClient.Delete(ctx, resourceGroup, vmssName)
	if err != nil {
		return err
	}
	if err = future.WaitForCompletionRef(ctx, az.virtualMachineScaleSetsClient.Client); err != nil {
		return err
	}
	_, err = future.Result(az.virtualMachineScaleSetsClient)
	return err
}

// SetVirtualMachineScaleSetCapacity sets the VMSS capacity
func (az *AzureClient) SetVirtualMachineScaleSetCapacity(ctx context.Context, resourceGroup, virtualMachineScaleSet string, sku azcompute.Sku, location string) error {

	s := compute.Sku{}
	DeepAssignment(&s, sku)
	future, err := az.virtualMachineScaleSetsClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		virtualMachineScaleSet,
		compute.VirtualMachineScaleSet{
			Location: &location,
			Sku:      &s,
		})
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachineScaleSetsClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachineScaleSetsClient)
	return err
}
