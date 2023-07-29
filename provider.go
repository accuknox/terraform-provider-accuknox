package main

import (
	"github.com/accuknox/terraform-provider-accuknox/discoveryengine"
	"github.com/accuknox/terraform-provider-accuknox/kubearmor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"accuknox_kubearmor_namespace_visibility":            kubearmor.ResourceKubearmorNamespaceVisibility(),
			"accuknox_kubearmor_namespace_posture":               kubearmor.ResourceKubearmorNamespacePosture(),
			"accuknox_kubearmor_security_policy":                 kubearmor.ResourceKubearmorSecurityPolicy(),
			"accuknox_kubearmor_host_security_policy":            kubearmor.ResourceKubearmorHostSecurityPolicy(),
			"accuknox_kubearmor_configuration":                   kubearmor.ResourceKubearmorConfiguration(),
			"accuknox_discovery_engine_configuration":            discoveryengine.ResourceDiscoveryEngineConfiguration(),
			"accuknox_discovery_engine_discovered_policy":        discoveryengine.ResourceDiscoveryEngineDiscoveredPolicy(),
			"accuknox_discovery_engine_enable_discovered_policy": discoveryengine.ResourceDiscoveryEngineEnableDiscoveredPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"accuknox_kubearmor_configuration":            kubearmor.DataSourceKubearmorConfiguration(),
			"accuknox_kubearmor_namespace_visibility":     kubearmor.DataSourceKubearmorNsVisibility(),
			"accuknox_kubearmor_namespace_posture":        kubearmor.DataSourceKubearmorNsPosture(),
			"accuknox_kubearmor_security_policy":          kubearmor.DataSourceKubearmorSecurityPolicy(),
			"accuknox_kubearmor_host_security_policy":     kubearmor.DataSourceKubearmorHostSecurityPolicy(),
			"accuknox_kubearmor_installed_version":        kubearmor.DataSourceKubearmorInstalledVersion(),
			"accuknox_kubearmor_stable_version":           kubearmor.DataSourceKubearmorStableVersion(),
			"accuknox_kubearmor_node":                     kubearmor.DataSourceKubearmorNode(),
			"accuknox_discovery_engine_configuration":     discoveryengine.DataSourceDiscoveryEngineConfiguration(),
			"accuknox_discovery_engine_discovered_policy": discoveryengine.DataSourceDiscoveryEngineDiscoveredPolicy(),
		},
	}
}
