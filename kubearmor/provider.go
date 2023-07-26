package kubearmor

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"accuknox_kubearmor_namespace_visibility": resourceNsVisibility(),
			"accuknox_kubearmor_namespace_posture":    resourceNsPosture(),
			"accuknox_kubearmor_security_policy":      resourceKubearmorSecurityPolicy(),
			"accuknox_kubearmor_host_security_policy": resourceKubearmorHostSecurityPolicy(),
			"accuknox_kubearmor_configuration":        resourceKubearmorConfiguration(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"accuknox_kubearmor_configuration":        dataSourceKubearmorConfiguration(),
			"accuknox_kubearmor_namespace_visibility": dataSourceKubearmorNsVisibility(),
			"accuknox_kubearmor_namespace_posture":    dataSourceKubearmorNsPosture(),
			"accuknox_kubearmor_security_policy":      dataSourceKubearmorSecurityPolicy(),
			"accuknox_kubearmor_host_security_policy": dataSourceKubearmorHostSecurityPolicy(),
			"accuknox_kubearmor_installed_version":    dataSourceKubearmorInstalledVersion(),
			"accuknox_kubearmor_stable_version":       dataSourceKubearmorStableVersion(),
			"accuknox_kubearmor_node":                 dataSourceKubearmorNode(),
		},
	}
}
