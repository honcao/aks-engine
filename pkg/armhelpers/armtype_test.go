package armhelpers

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/aks-engine/pkg/api"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               3,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.10.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "azure",
				},
			},
			FeatureFlags: &api.FeatureFlags{
				BlockOutboundInternet: false,
			},
		},
	}
	armObject := CreateCustomScriptExtension(cs)

	jsonObj, _ := json.MarshalIndent(armObject, "", "   ")
	fmt.Println(string(jsonObj))
}
