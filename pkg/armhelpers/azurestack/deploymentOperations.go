// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package azurestack

import (
	"context"

	"github.com/Azure/aks-engine/pkg/armhelpers"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	azresources "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
)

//DeploymentOperationsListResultPageClient Deployment Operations List Result Page Client
type DeploymentOperationsListResultPageClient struct {
	dolrp resources.DeploymentOperationsListResultPage
	err   error
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *DeploymentOperationsListResultPageClient) Next() error {
	return page.dolrp.Next()
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page DeploymentOperationsListResultPageClient) NotDone() bool {
	return page.dolrp.NotDone()
}

// Response returns the raw server response from the last page request.
func (page DeploymentOperationsListResultPageClient) Response() azresources.DeploymentOperationsListResult {
	r := azresources.DeploymentOperationsListResult{}
	DeepAssignment(&r, page.dolrp.Response())
	return r
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page DeploymentOperationsListResultPageClient) Values() []azresources.DeploymentOperation {
	r := []azresources.DeploymentOperation{}
	DeepAssignment(r, page.dolrp.Values())
	return r
}

// ListDeploymentOperations gets all deployments operations for a deployment.
func (az *AzureClient) ListDeploymentOperations(ctx context.Context, resourceGroupName string, deploymentName string, top *int32) (armhelpers.DeploymentOperationsListResultPage, error) {
	list, err := az.deploymentOperationsClient.List(ctx, resourceGroupName, deploymentName, top)
	c := DeploymentOperationsListResultPageClient{
		dolrp: list,
		err:   err,
	}
	return &c, err
}
