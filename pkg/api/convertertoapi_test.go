// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/apimachinery/pkg/api/equality"

	"github.com/Azure/aks-engine/pkg/api/common"
	v20170701 "github.com/Azure/aks-engine/pkg/api/v20170701"
	"github.com/Azure/aks-engine/pkg/api/vlabs"
	"github.com/Azure/go-autorest/autorest/azure"
)

func TestAddDCOSPublicAgentPool(t *testing.T) {
	expectedNumPools := 2
	for _, masterCount := range [2]int{1, 3} {
		profiles := []*AgentPoolProfile{}
		profile := makeAgentPoolProfile(1, "agentprivate", "test-dcos-pool", "Standard_D2_v2", "Linux")
		profiles = append(profiles, profile)
		master := makeMasterProfile(masterCount, "test-dcos", "Standard_D2_v2")
		props := getProperties(profiles, master)
		expectedPublicPoolName := props.AgentPoolProfiles[0].Name + publicAgentPoolSuffix
		expectedPublicDNSPrefix := props.AgentPoolProfiles[0].DNSPrefix
		expectedPrivateDNSPrefix := ""
		expectedPublicOSType := props.AgentPoolProfiles[0].OSType
		expectedPublicVMSize := props.AgentPoolProfiles[0].VMSize
		addDCOSPublicAgentPool(props)
		if len(props.AgentPoolProfiles) != expectedNumPools {
			t.Fatalf("incorrect agent pools count. expected=%d actual=%d", expectedNumPools, len(props.AgentPoolProfiles))
		}
		if props.AgentPoolProfiles[1].Name != expectedPublicPoolName {
			t.Fatalf("incorrect public pool name. expected=%s actual=%s", expectedPublicPoolName, props.AgentPoolProfiles[1].Name)
		}
		if props.AgentPoolProfiles[1].DNSPrefix != expectedPublicDNSPrefix {
			t.Fatalf("incorrect public pool DNS prefix. expected=%s actual=%s", expectedPublicDNSPrefix, props.AgentPoolProfiles[1].DNSPrefix)
		}
		if props.AgentPoolProfiles[0].DNSPrefix != expectedPrivateDNSPrefix {
			t.Fatalf("incorrect private pool DNS prefix. expected=%s actual=%s", expectedPrivateDNSPrefix, props.AgentPoolProfiles[0].DNSPrefix)
		}
		if props.AgentPoolProfiles[1].OSType != expectedPublicOSType {
			t.Fatalf("incorrect public pool OS type. expected=%s actual=%s", expectedPublicOSType, props.AgentPoolProfiles[1].OSType)
		}
		if props.AgentPoolProfiles[1].VMSize != expectedPublicVMSize {
			t.Fatalf("incorrect public pool VM size. expected=%s actual=%s", expectedPublicVMSize, props.AgentPoolProfiles[1].VMSize)
		}
		for i, port := range [3]int{80, 443, 8080} {
			if props.AgentPoolProfiles[1].Ports[i] != port {
				t.Fatalf("incorrect public pool port assignment. expected=%d actual=%d", port, props.AgentPoolProfiles[1].Ports[i])
			}
		}
		if props.AgentPoolProfiles[1].Count != masterCount {
			t.Fatalf("incorrect public pool VM size. expected=%d actual=%d", masterCount, props.AgentPoolProfiles[1].Count)
		}
	}
}

func makeAgentPoolProfile(count int, name, dNSPrefix, vMSize string, oSType OSType) *AgentPoolProfile {
	return &AgentPoolProfile{
		Name:      name,
		Count:     count,
		DNSPrefix: dNSPrefix,
		OSType:    oSType,
		VMSize:    vMSize,
	}
}

func makeMasterProfile(count int, dNSPrefix, vMSize string) *MasterProfile {
	return &MasterProfile{
		Count:     count,
		DNSPrefix: "test-dcos",
		VMSize:    "Standard_D2_v2",
	}
}

func getProperties(profiles []*AgentPoolProfile, master *MasterProfile) *Properties {
	return &Properties{
		AgentPoolProfiles: profiles,
		MasterProfile:     master,
	}
}

func TestOrchestratorVersion(t *testing.T) {
	// test v20170701
	v20170701cs := &v20170701.ContainerService{
		Properties: &v20170701.Properties{
			OrchestratorProfile: &v20170701.OrchestratorProfile{
				OrchestratorType: v20170701.Kubernetes,
			},
		},
	}
	cs := ConvertV20170701ContainerService(v20170701cs, false)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != common.GetDefaultKubernetesVersion(false) {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	v20170701cs = &v20170701.ContainerService{
		Properties: &v20170701.Properties{
			OrchestratorProfile: &v20170701.OrchestratorProfile{
				OrchestratorType:    v20170701.Kubernetes,
				OrchestratorVersion: "1.7.14",
			},
		},
	}
	cs = ConvertV20170701ContainerService(v20170701cs, true)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != "1.7.14" {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}
	// test vlabs
	vlabscs := &vlabs.ContainerService{
		Properties: &vlabs.Properties{
			OrchestratorProfile: &vlabs.OrchestratorProfile{
				OrchestratorType: vlabs.Kubernetes,
			},
		},
	}
	cs = ConvertVLabsContainerService(vlabscs, false)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != common.GetDefaultKubernetesVersion(false) {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	vlabscs = &vlabs.ContainerService{
		Properties: &vlabs.Properties{
			OrchestratorProfile: &vlabs.OrchestratorProfile{
				OrchestratorType:    vlabs.Kubernetes,
				OrchestratorVersion: "1.7.15",
			},
		},
	}
	cs = ConvertVLabsContainerService(vlabscs, false)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != "1.7.15" {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}
}

func TestKubernetesVlabsDefaults(t *testing.T) {
	vp := makeKubernetesPropertiesVlabs()
	ap := makeKubernetesProperties()
	setVlabsKubernetesDefaults(vp, ap.OrchestratorProfile)
	if ap.OrchestratorProfile.KubernetesConfig == nil {
		t.Fatalf("KubernetesConfig cannot be nil after vlabs default conversion")
	}
	if ap.OrchestratorProfile.KubernetesConfig.NetworkPlugin != vlabs.DefaultNetworkPlugin {
		t.Fatalf("vlabs defaults not applied, expected NetworkPlugin: %s, instead got: %s", vlabs.DefaultNetworkPlugin, ap.OrchestratorProfile.KubernetesConfig.NetworkPlugin)
	}
	if ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy != vlabs.DefaultNetworkPolicy {
		t.Fatalf("vlabs defaults not applied, expected NetworkPolicy: %s, instead got: %s", vlabs.DefaultNetworkPolicy, ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
	}

	vp = makeKubernetesPropertiesVlabs()
	vp.WindowsProfile = &vlabs.WindowsProfile{}
	vp.AgentPoolProfiles = append(vp.AgentPoolProfiles, &vlabs.AgentPoolProfile{OSType: "Windows"})
	ap = makeKubernetesProperties()
	setVlabsKubernetesDefaults(vp, ap.OrchestratorProfile)
	if ap.OrchestratorProfile.KubernetesConfig == nil {
		t.Fatalf("KubernetesConfig cannot be nil after vlabs default conversion")
	}
	if ap.OrchestratorProfile.KubernetesConfig.NetworkPlugin != vlabs.DefaultNetworkPluginWindows {
		t.Fatalf("vlabs defaults not applied, expected NetworkPlugin: %s, instead got: %s", vlabs.DefaultNetworkPluginWindows, ap.OrchestratorProfile.KubernetesConfig.NetworkPlugin)
	}
	if ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy != vlabs.DefaultNetworkPolicy {
		t.Fatalf("vlabs defaults not applied, expected NetworkPolicy: %s, instead got: %s", vlabs.DefaultNetworkPolicy, ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
	}
}

func TestConvertVLabsKubernetesConfigProfile(t *testing.T) {
	tests := map[string]struct {
		props  *vlabs.KubernetesConfig
		expect *KubernetesConfig
	}{
		"WindowsNodeBinariesURL": {
			props: &vlabs.KubernetesConfig{
				WindowsNodeBinariesURL: "http://test/test.tar.gz",
			},
			expect: &KubernetesConfig{
				WindowsNodeBinariesURL: "http://test/test.tar.gz",
			},
		},
	}

	for name, test := range tests {
		t.Logf("running scenario %q", name)
		actual := &KubernetesConfig{}
		convertVLabsKubernetesConfig(test.props, actual)
		if !equality.Semantic.DeepEqual(test.expect, actual) {
			t.Errorf(spew.Sprintf("Expected:\n%+v\nGot:\n%+v", test.expect, actual))
		}
	}
}

func makeKubernetesProperties() *Properties {
	ap := &Properties{}
	ap.OrchestratorProfile = &OrchestratorProfile{}
	ap.OrchestratorProfile.OrchestratorType = "Kubernetes"
	return ap
}

func makeKubernetesPropertiesVlabs() *vlabs.Properties {
	vp := &vlabs.Properties{}
	vp.OrchestratorProfile = &vlabs.OrchestratorProfile{}
	vp.OrchestratorProfile.OrchestratorType = "Kubernetes"
	return vp
}

func TestConvertCustomFilesToAPI(t *testing.T) {
	expectedAPICustomFiles := []CustomFile{
		{
			Source: "/test/source",
			Dest:   "/test/dest",
		},
	}
	masterProfile := MasterProfile{}

	vp := &vlabs.MasterProfile{}
	vp.CustomFiles = &[]vlabs.CustomFile{
		{
			Source: "/test/source",
			Dest:   "/test/dest",
		},
	}
	convertCustomFilesToAPI(vp, &masterProfile)
	if !equality.Semantic.DeepEqual(&expectedAPICustomFiles, masterProfile.CustomFiles) {
		t.Fatalf("convertCustomFilesToApi conversion of vlabs.MasterProfile did not convert correctly")
	}
}

func TestCustomCloudProfile(t *testing.T) {
	const (
		name                         = "AzureStackCloud"
		managementPortalURL          = "https=//management.local.azurestack.external/"
		publishSettingsURL           = "https=//management.local.azurestack.external/publishsettings/index"
		serviceManagementEndpoint    = "https=//management.azurestackci15.onmicrosoft.com/36f71706-54df-4305-9847-5b038a4cf189"
		resourceManagerEndpoint      = "https=//management.local.azurestack.external/"
		activeDirectoryEndpoint      = "https=//login.windows.net/"
		galleryEndpoint              = "https=//portal.local.azurestack.external=30015/"
		keyVaultEndpoint             = "https=//vault.azurestack.external/"
		graphEndpoint                = "https=//graph.windows.net/"
		serviceBusEndpoint           = "https=//servicebus.azurestack.external/"
		batchManagementEndpoint      = "https=//batch.azurestack.external/"
		storageEndpointSuffix        = "core.azurestack.external"
		sqlDatabaseDNSSuffix         = "database.azurestack.external"
		trafficManagerDNSSuffix      = "trafficmanager.cn"
		keyVaultDNSSuffix            = "vault.azurestack.external"
		serviceBusEndpointSuffix     = "servicebus.azurestack.external"
		serviceManagementVMDNSSuffix = "chinacloudapp.cn"
		resourceManagerVMDNSSuffix   = "cloudapp.azurestack.external"
		containerRegistryDNSSuffix   = "azurecr.io"
		tokenAudience                = "https=//management.azurestack.external/"
	)

	vlabscs := &vlabs.ContainerService{
		Properties: &vlabs.Properties{
			CustomCloudProfile: &vlabs.CustomCloudProfile{
				IdentitySystem:       ADFS,
				AuthenticationMethod: ClientCertificate,
				Environment: &azure.Environment{
					Name:                         name,
					ManagementPortalURL:          managementPortalURL,
					PublishSettingsURL:           publishSettingsURL,
					ServiceManagementEndpoint:    serviceManagementEndpoint,
					ResourceManagerEndpoint:      resourceManagerEndpoint,
					ActiveDirectoryEndpoint:      activeDirectoryEndpoint,
					GalleryEndpoint:              galleryEndpoint,
					KeyVaultEndpoint:             keyVaultEndpoint,
					GraphEndpoint:                graphEndpoint,
					ServiceBusEndpoint:           serviceBusEndpoint,
					BatchManagementEndpoint:      batchManagementEndpoint,
					StorageEndpointSuffix:        storageEndpointSuffix,
					SQLDatabaseDNSSuffix:         sqlDatabaseDNSSuffix,
					TrafficManagerDNSSuffix:      trafficManagerDNSSuffix,
					KeyVaultDNSSuffix:            keyVaultDNSSuffix,
					ServiceBusEndpointSuffix:     serviceBusEndpointSuffix,
					ServiceManagementVMDNSSuffix: serviceManagementVMDNSSuffix,
					ResourceManagerVMDNSSuffix:   resourceManagerVMDNSSuffix,
					ContainerRegistryDNSSuffix:   containerRegistryDNSSuffix,
					TokenAudience:                tokenAudience,
				},
			},
		},
	}

	cs := ConvertVLabsContainerService(vlabscs, false)

	if cs.Properties.CustomCloudProfile.AuthenticationMethod != ClientCertificate {
		t.Errorf("incorrect AuthenticationMethod, expect: '%s', actual: '%s'", ClientCertificate, cs.Properties.CustomCloudProfile.AuthenticationMethod)
	}
	if cs.Properties.CustomCloudProfile.IdentitySystem != ADFS {
		t.Errorf("incorrect IdentitySystem, expect: '%s', actual: '%s'", ADFS, cs.Properties.CustomCloudProfile.IdentitySystem)
	}

	if cs.Properties.CustomCloudProfile.Environment.Name != name {
		t.Errorf("incorrect Name, expect: '%s', actual: '%s'", name, cs.Properties.CustomCloudProfile.Environment.Name)
	}
	if cs.Properties.CustomCloudProfile.Environment.ManagementPortalURL != managementPortalURL {
		t.Errorf("incorrect ManagementPortalURL, expect: '%s', actual: '%s'", managementPortalURL, cs.Properties.CustomCloudProfile.Environment.ManagementPortalURL)
	}
	if cs.Properties.CustomCloudProfile.Environment.PublishSettingsURL != publishSettingsURL {
		t.Errorf("incorrect PublishSettingsURL, expect: '%s', actual: '%s'", publishSettingsURL, cs.Properties.CustomCloudProfile.Environment.PublishSettingsURL)
	}
	if cs.Properties.CustomCloudProfile.Environment.ServiceManagementEndpoint != serviceManagementEndpoint {
		t.Errorf("incorrect ServiceManagementEndpoint, expect: '%s', actual: '%s'", serviceManagementEndpoint, cs.Properties.CustomCloudProfile.Environment.ServiceManagementEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.ResourceManagerEndpoint != resourceManagerEndpoint {
		t.Errorf("incorrect ResourceManagerEndpoint, expect: '%s', actual: '%s'", resourceManagerEndpoint, cs.Properties.CustomCloudProfile.Environment.ResourceManagerEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.ActiveDirectoryEndpoint != activeDirectoryEndpoint {
		t.Errorf("incorrect ActiveDirectoryEndpoint, expect: '%s', actual: '%s'", activeDirectoryEndpoint, cs.Properties.CustomCloudProfile.Environment.ActiveDirectoryEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.GalleryEndpoint != galleryEndpoint {
		t.Errorf("incorrect GalleryEndpoint, expect: '%s', actual: '%s'", galleryEndpoint, cs.Properties.CustomCloudProfile.Environment.GalleryEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.KeyVaultEndpoint != keyVaultEndpoint {
		t.Errorf("incorrect KeyVaultEndpoint, expect: '%s', actual: '%s'", keyVaultEndpoint, cs.Properties.CustomCloudProfile.Environment.KeyVaultEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.GraphEndpoint != graphEndpoint {
		t.Errorf("incorrect GraphEndpoint, expect: '%s', actual: '%s'", graphEndpoint, cs.Properties.CustomCloudProfile.Environment.GraphEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.ServiceBusEndpoint != serviceBusEndpoint {
		t.Errorf("incorrect ServiceBusEndpoint, expect: '%s', actual: '%s'", serviceBusEndpoint, cs.Properties.CustomCloudProfile.Environment.ServiceBusEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.BatchManagementEndpoint != batchManagementEndpoint {
		t.Errorf("incorrect BatchManagementEndpoint, expect: '%s', actual: '%s'", batchManagementEndpoint, cs.Properties.CustomCloudProfile.Environment.BatchManagementEndpoint)
	}
	if cs.Properties.CustomCloudProfile.Environment.StorageEndpointSuffix != storageEndpointSuffix {
		t.Errorf("incorrect StorageEndpointSuffix, expect: '%s', actual: '%s'", storageEndpointSuffix, cs.Properties.CustomCloudProfile.Environment.StorageEndpointSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.SQLDatabaseDNSSuffix != sqlDatabaseDNSSuffix {
		t.Errorf("incorrect SQLDatabaseDNSSuffix, expect: '%s', actual: '%s'", sqlDatabaseDNSSuffix, cs.Properties.CustomCloudProfile.Environment.SQLDatabaseDNSSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.TrafficManagerDNSSuffix != trafficManagerDNSSuffix {
		t.Errorf("incorrect TrafficManagerDNSSuffix, expect: '%s', actual: '%s'", trafficManagerDNSSuffix, cs.Properties.CustomCloudProfile.Environment.TrafficManagerDNSSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.KeyVaultDNSSuffix != keyVaultDNSSuffix {
		t.Errorf("incorrect KeyVaultDNSSuffix, expect: '%s', actual: '%s'", keyVaultDNSSuffix, cs.Properties.CustomCloudProfile.Environment.KeyVaultDNSSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.ServiceBusEndpointSuffix != serviceBusEndpointSuffix {
		t.Errorf("incorrect ServiceBusEndpointSuffix, expect: '%s', actual: '%s'", serviceBusEndpointSuffix, cs.Properties.CustomCloudProfile.Environment.ServiceBusEndpointSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.ServiceManagementVMDNSSuffix != serviceManagementVMDNSSuffix {
		t.Errorf("incorrect ServiceManagementVMDNSSuffix, expect: '%s', actual: '%s'", serviceManagementVMDNSSuffix, cs.Properties.CustomCloudProfile.Environment.ServiceManagementVMDNSSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.ResourceManagerVMDNSSuffix != resourceManagerVMDNSSuffix {
		t.Errorf("incorrect ResourceManagerVMDNSSuffix, expect: '%s', actual: '%s'", resourceManagerVMDNSSuffix, cs.Properties.CustomCloudProfile.Environment.ResourceManagerVMDNSSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.ContainerRegistryDNSSuffix != containerRegistryDNSSuffix {
		t.Errorf("incorrect ContainerRegistryDNSSuffix, expect: '%s', actual: '%s'", containerRegistryDNSSuffix, cs.Properties.CustomCloudProfile.Environment.ContainerRegistryDNSSuffix)
	}
	if cs.Properties.CustomCloudProfile.Environment.TokenAudience != tokenAudience {
		t.Errorf("incorrect TokenAudience, expect: '%s', actual: '%s'", tokenAudience, cs.Properties.CustomCloudProfile.Environment.TokenAudience)
	}
}

func TestConvertAzureEnvironmentSpecConfig(t *testing.T) {
	//Mock AzureEnvironmentSpecConfig
	vlabscs := &vlabs.ContainerService{
		Properties: &vlabs.Properties{
			CustomCloudProfile: &vlabs.CustomCloudProfile{
				IdentitySystem:       ADFS,
				AuthenticationMethod: ClientCertificate,
				AzureEnvironmentSpecConfig: &vlabs.AzureEnvironmentSpecConfig{
					CloudName: "AzureStackCloud",
					//DockerSpecConfig specify the docker engine download repo
					DockerSpecConfig: vlabs.DockerSpecConfig{
						DockerEngineRepo:         "DockerEngineRepo",
						DockerComposeDownloadURL: "DockerComposeDownloadURL",
					},
					//KubernetesSpecConfig - Due to Chinese firewall issue, the default containers from google is blocked, use the Chinese local mirror instead
					KubernetesSpecConfig: vlabs.KubernetesSpecConfig{
						KubernetesImageBase:              "KubernetesImageBase",
						TillerImageBase:                  "TillerImageBase",
						ACIConnectorImageBase:            "ACIConnectorImageBase",
						NVIDIAImageBase:                  "NVIDIAImageBase",
						AzureCNIImageBase:                "AzureCNIImageBase",
						EtcdDownloadURLBase:              "EtcdDownloadURLBase",
						KubeBinariesSASURLBase:           "KubeBinariesSASURLBase",
						WindowsTelemetryGUID:             "WindowsTelemetryGUID",
						CNIPluginsDownloadURL:            "CNIPluginsDownloadURL",
						VnetCNILinuxPluginsDownloadURL:   "VnetCNILinuxPluginsDownloadURL",
						VnetCNIWindowsPluginsDownloadURL: "VnetCNIWindowsPluginsDownloadURL",
						ContainerdDownloadURLBase:        "ContainerdDownloadURLBase",
					},
					DCOSSpecConfig: vlabs.DCOSSpecConfig{
						DCOS188BootstrapDownloadURL:     "DCOS188BootstrapDownloadURL",
						DCOS190BootstrapDownloadURL:     "DCOS190BootstrapDownloadURL",
						DCOS198BootstrapDownloadURL:     "DCOS198BootstrapDownloadURL",
						DCOS110BootstrapDownloadURL:     "DCOS110BootstrapDownloadURL",
						DCOS111BootstrapDownloadURL:     "DCOS111BootstrapDownloadURL",
						DCOSWindowsBootstrapDownloadURL: "DCOSWindowsBootstrapDownloadURL",
						DcosRepositoryURL:               "DcosRepositoryURL",
						DcosClusterPackageListID:        "DcosClusterPackageListID",
						DcosProviderPackageID:           "DcosProviderPackageID",
					},
					EndpointConfig: vlabs.AzureEndpointConfig{
						ResourceManagerVMDNSSuffix: "ResourceManagerVMDNSSuffix",
					},
					OSImageConfig: map[vlabs.Distro]vlabs.AzureOSImageConfig{
						vlabs.Distro("Test"): {
							ImageOffer:     "ImageOffer",
							ImageSku:       "ImageSku",
							ImagePublisher: "ImagePublisher",
							ImageVersion:   "ImageVersion",
						},
						vlabs.AKS:             vlabs.DefaultAKSOSImageConfig,
						vlabs.AKSDockerEngine: vlabs.DefaultAKSDockerEngineOSImageConfig,
					},
				},
			},
		},
	}
	cs := ConvertVLabsContainerService(vlabscs, false)

	if cs.Properties.CustomCloudProfile.AuthenticationMethod != ClientCertificate {
		t.Errorf("incorrect AuthenticationMethod, expect: '%s', actual: '%s'", ClientCertificate, cs.Properties.CustomCloudProfile.AuthenticationMethod)
	}
	if cs.Properties.CustomCloudProfile.IdentitySystem != ADFS {
		t.Errorf("incorrect IdentitySystem, expect: '%s', actual: '%s'", ADFS, cs.Properties.CustomCloudProfile.IdentitySystem)
	}

	csSpec := cs.Properties.CustomCloudProfile.AzureEnvironmentSpecConfig
	vlabscsSpec := vlabscs.Properties.CustomCloudProfile.AzureEnvironmentSpecConfig

	if csSpec.CloudName != vlabscsSpec.CloudName {
		t.Errorf("incorrect CloudName, expect: '%s', actual: '%s'", vlabscsSpec.CloudName, csSpec.CloudName)
	}

	//KubernetesSpecConfig
	if csSpec.KubernetesSpecConfig.KubernetesImageBase != vlabscsSpec.KubernetesSpecConfig.KubernetesImageBase {
		t.Errorf("incorrect KubernetesImageBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.KubernetesImageBase, csSpec.KubernetesSpecConfig.KubernetesImageBase)
	}
	if csSpec.KubernetesSpecConfig.TillerImageBase != vlabscsSpec.KubernetesSpecConfig.TillerImageBase {
		t.Errorf("incorrect TillerImageBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.TillerImageBase, csSpec.KubernetesSpecConfig.TillerImageBase)
	}
	if csSpec.KubernetesSpecConfig.ACIConnectorImageBase != vlabscsSpec.KubernetesSpecConfig.ACIConnectorImageBase {
		t.Errorf("incorrect ACIConnectorImageBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.ACIConnectorImageBase, csSpec.KubernetesSpecConfig.ACIConnectorImageBase)
	}
	if csSpec.KubernetesSpecConfig.NVIDIAImageBase != vlabscsSpec.KubernetesSpecConfig.NVIDIAImageBase {
		t.Errorf("incorrect NVIDIAImageBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.NVIDIAImageBase, csSpec.KubernetesSpecConfig.NVIDIAImageBase)
	}
	if csSpec.KubernetesSpecConfig.AzureCNIImageBase != vlabscsSpec.KubernetesSpecConfig.AzureCNIImageBase {
		t.Errorf("incorrect AzureCNIImageBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.AzureCNIImageBase, csSpec.KubernetesSpecConfig.AzureCNIImageBase)
	}
	if csSpec.KubernetesSpecConfig.EtcdDownloadURLBase != vlabscsSpec.KubernetesSpecConfig.EtcdDownloadURLBase {
		t.Errorf("incorrect EtcdDownloadURLBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.EtcdDownloadURLBase, csSpec.KubernetesSpecConfig.EtcdDownloadURLBase)
	}
	if csSpec.KubernetesSpecConfig.KubeBinariesSASURLBase != vlabscsSpec.KubernetesSpecConfig.KubeBinariesSASURLBase {
		t.Errorf("incorrect KubeBinariesSASURLBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.KubeBinariesSASURLBase, csSpec.KubernetesSpecConfig.KubeBinariesSASURLBase)
	}
	if csSpec.KubernetesSpecConfig.WindowsTelemetryGUID != vlabscsSpec.KubernetesSpecConfig.WindowsTelemetryGUID {
		t.Errorf("incorrect WindowsTelemetryGUID, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.WindowsTelemetryGUID, csSpec.KubernetesSpecConfig.WindowsTelemetryGUID)
	}
	if csSpec.KubernetesSpecConfig.CNIPluginsDownloadURL != vlabscsSpec.KubernetesSpecConfig.CNIPluginsDownloadURL {
		t.Errorf("incorrect CNIPluginsDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.CNIPluginsDownloadURL, csSpec.KubernetesSpecConfig.CNIPluginsDownloadURL)
	}
	if csSpec.KubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL != vlabscsSpec.KubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL {
		t.Errorf("incorrect VnetCNILinuxPluginsDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL, csSpec.KubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL)
	}
	if csSpec.KubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL != vlabscsSpec.KubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL {
		t.Errorf("incorrect VnetCNIWindowsPluginsDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL, csSpec.KubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL)
	}
	if csSpec.KubernetesSpecConfig.ContainerdDownloadURLBase != vlabscsSpec.KubernetesSpecConfig.ContainerdDownloadURLBase {
		t.Errorf("incorrect ContainerdDownloadURLBase, expect: '%s', actual: '%s'", vlabscsSpec.KubernetesSpecConfig.ContainerdDownloadURLBase, csSpec.KubernetesSpecConfig.ContainerdDownloadURLBase)
	}

	//DockerSpecConfig
	if csSpec.DockerSpecConfig.DockerComposeDownloadURL != vlabscsSpec.DockerSpecConfig.DockerComposeDownloadURL {
		t.Errorf("incorrect DockerComposeDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.DockerSpecConfig.DockerComposeDownloadURL, csSpec.DockerSpecConfig.DockerComposeDownloadURL)
	}
	if csSpec.DockerSpecConfig.DockerEngineRepo != vlabscsSpec.DockerSpecConfig.DockerEngineRepo {
		t.Errorf("incorrect DockerEngineRepo, expect: '%s', actual: '%s'", vlabscsSpec.DockerSpecConfig.DockerEngineRepo, csSpec.DockerSpecConfig.DockerEngineRepo)
	}

	//DCOSSpecConfig
	if csSpec.DCOSSpecConfig.DCOS188BootstrapDownloadURL != vlabscsSpec.DCOSSpecConfig.DCOS188BootstrapDownloadURL {
		t.Errorf("incorrect DCOS188BootstrapDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DCOS188BootstrapDownloadURL, csSpec.DCOSSpecConfig.DCOS188BootstrapDownloadURL)
	}
	if csSpec.DCOSSpecConfig.DCOS190BootstrapDownloadURL != vlabscsSpec.DCOSSpecConfig.DCOS190BootstrapDownloadURL {
		t.Errorf("incorrect DCOS190BootstrapDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DCOS190BootstrapDownloadURL, csSpec.DCOSSpecConfig.DCOS190BootstrapDownloadURL)
	}
	if csSpec.DCOSSpecConfig.DCOS198BootstrapDownloadURL != vlabscsSpec.DCOSSpecConfig.DCOS198BootstrapDownloadURL {
		t.Errorf("incorrect DCOS198BootstrapDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DCOS198BootstrapDownloadURL, csSpec.DCOSSpecConfig.DCOS198BootstrapDownloadURL)
	}
	if csSpec.DCOSSpecConfig.DCOS110BootstrapDownloadURL != vlabscsSpec.DCOSSpecConfig.DCOS110BootstrapDownloadURL {
		t.Errorf("incorrect DCOS110BootstrapDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DCOS110BootstrapDownloadURL, csSpec.DCOSSpecConfig.DCOS110BootstrapDownloadURL)
	}
	if csSpec.DCOSSpecConfig.DCOS111BootstrapDownloadURL != vlabscsSpec.DCOSSpecConfig.DCOS111BootstrapDownloadURL {
		t.Errorf("incorrect DCOS111BootstrapDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DCOS111BootstrapDownloadURL, csSpec.DCOSSpecConfig.DCOS111BootstrapDownloadURL)
	}
	if csSpec.DCOSSpecConfig.DCOSWindowsBootstrapDownloadURL != vlabscsSpec.DCOSSpecConfig.DCOSWindowsBootstrapDownloadURL {
		t.Errorf("incorrect DCOSWindowsBootstrapDownloadURL, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DCOSWindowsBootstrapDownloadURL, csSpec.DCOSSpecConfig.DCOSWindowsBootstrapDownloadURL)
	}
	if csSpec.DCOSSpecConfig.DcosRepositoryURL != vlabscsSpec.DCOSSpecConfig.DcosRepositoryURL {
		t.Errorf("incorrect DcosRepositoryURL, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DcosRepositoryURL, csSpec.DCOSSpecConfig.DcosRepositoryURL)
	}
	if csSpec.DCOSSpecConfig.DcosClusterPackageListID != vlabscsSpec.DCOSSpecConfig.DcosClusterPackageListID {
		t.Errorf("incorrect DcosClusterPackageListID, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DcosClusterPackageListID, csSpec.DCOSSpecConfig.DcosClusterPackageListID)
	}
	if csSpec.DCOSSpecConfig.DcosProviderPackageID != vlabscsSpec.DCOSSpecConfig.DcosProviderPackageID {
		t.Errorf("incorrect DcosProviderPackageID, expect: '%s', actual: '%s'", vlabscsSpec.DCOSSpecConfig.DcosProviderPackageID, csSpec.DCOSSpecConfig.DcosProviderPackageID)
	}

	//EndpointConfig
	if csSpec.EndpointConfig.ResourceManagerVMDNSSuffix != vlabscsSpec.EndpointConfig.ResourceManagerVMDNSSuffix {
		t.Errorf("incorrect ResourceManagerVMDNSSuffix, expect: '%s', actual: '%s'", vlabscsSpec.EndpointConfig.ResourceManagerVMDNSSuffix, csSpec.EndpointConfig.ResourceManagerVMDNSSuffix)
	}

	//OSImageConfig
	for k, v := range vlabscsSpec.OSImageConfig {
		if actualValue, ok := csSpec.OSImageConfig[Distro(string(k))]; ok {
			if v.ImageOffer != actualValue.ImageOffer {
				t.Errorf("incorrect ImageOffer for '%s', expect: '%s', actual: '%s'", string(k), v.ImageOffer, actualValue.ImageOffer)
			}
			if v.ImagePublisher != actualValue.ImagePublisher {
				t.Errorf("incorrect ImagePublisher for '%s', expect: '%s', actual: '%s'", string(k), v.ImagePublisher, actualValue.ImagePublisher)
			}
			if v.ImageSku != actualValue.ImageSku {
				t.Errorf("incorrect ImageSku for '%s', expect: '%s', actual: '%s'", string(k), v.ImageSku, actualValue.ImageSku)
			}
			if v.ImageVersion != actualValue.ImageVersion {
				t.Errorf("incorrect ImageVersion for '%s', expect: '%s', actual: '%s'", string(k), v.ImageVersion, actualValue.ImageVersion)
			}
		} else {
			t.Errorf("incorrect OSImageConfig: '%s' is missing", string(k))
		}
	}
}
