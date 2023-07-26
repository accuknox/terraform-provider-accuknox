package kubearmor

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kcV1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
)

func expandHostSelector(selector []interface{}) kcV1.NodeSelectorType {
	selectorType := kcV1.NodeSelectorType{}
	in := selector[0].(map[string]interface{})
	if v, ok := in["match_labels"].(map[string]interface{}); ok && len(v) > 0 {
		selectorType.MatchLabels = expandStringMap(v)
	}
	return selectorType
}

func expandHostCapability(matchCapability map[string]interface{}) kcV1.MatchHostCapabilitiesType {
	matchcapability := kcV1.MatchHostCapabilitiesType{}
	matchcapability.Capability = kcV1.MatchCapabilitiesStringType(matchCapability["capabilities"].(string))
	matchcapability.Action = kcV1.ActionType(matchCapability["action"].(string))
	matchcapability.FromSource = expandFromSource(matchCapability["from_source"].([]interface{}))

	return matchcapability
}

func expandHostMatchCapabilities(matchCapabilities []interface{}) []kcV1.MatchHostCapabilitiesType {
	matchcapabilities := []kcV1.MatchHostCapabilitiesType{}

	for _, matchCapability := range matchCapabilities {
		matchcapabilities = append(matchcapabilities, expandHostCapability(matchCapability.(map[string]interface{})))
	}

	return matchcapabilities
}

func expandHostCapabilities(capabilities []interface{}) kcV1.HostCapabilitiesType {
	capabilitiesType := kcV1.HostCapabilitiesType{}
	if len(capabilities) > 0 {
		in := capabilities[0].(map[string]interface{})
		capabilitiesType.MatchCapabilities = expandHostMatchCapabilities(in["match_paths"].([]interface{}))
	}

	return capabilitiesType
}

func expandHostProtocol(matchProtocol map[string]interface{}) kcV1.MatchHostNetworkProtocolType {
	matchprotocol := kcV1.MatchHostNetworkProtocolType{}
	matchprotocol.Protocol = kcV1.MatchNetworkProtocolStringType(matchProtocol["protocol"].(string))
	matchprotocol.Action = kcV1.ActionType(matchProtocol["action"].(string))
	matchprotocol.FromSource = expandFromSource(matchProtocol["from_source"].([]interface{}))

	return matchprotocol
}

func expandHostMatchProtocols(protocol []interface{}) []kcV1.MatchHostNetworkProtocolType {
	matchnetwork := []kcV1.MatchHostNetworkProtocolType{}

	for _, matchProtocol := range protocol {
		matchnetwork = append(matchnetwork, expandHostProtocol(matchProtocol.(map[string]interface{})))
	}

	return matchnetwork
}

func expandHostNetwork(network []interface{}) kcV1.HostNetworkType {
	networksType := kcV1.HostNetworkType{}
	if len(network) > 0 {
		in := network[0].(map[string]interface{})
		networksType.MatchProtocols = expandHostMatchProtocols(in["protocol"].([]interface{}))
	}

	return networksType
}

func expandHostSpec(d *schema.ResourceData) kcV1.KubeArmorHostPolicySpec {
	spec := kcV1.KubeArmorHostPolicySpec{}

	spec.Severity = kcV1.SeverityType(d.Get("severity").(int))
	spec.Action = kcV1.ActionType(d.Get("action").(string))
	spec.NodeSelector = expandHostSelector(d.Get("node_selector").([]interface{}))
	spec.File = expandFile(d.Get("file").([]interface{}))
	spec.Process = expandProcess(d.Get("process").([]interface{}))

	spec.Capabilities = expandHostCapabilities(d.Get("capabilities").([]interface{}))
	spec.Capabilities = kcV1.HostCapabilitiesType{
		MatchCapabilities: append([]kcV1.MatchHostCapabilitiesType{}, spec.Capabilities.MatchCapabilities...),
	}

	spec.Network = expandHostNetwork(d.Get("network").([]interface{}))
	spec.Network = kcV1.HostNetworkType{
		MatchProtocols: append([]kcV1.MatchHostNetworkProtocolType{}, spec.Network.MatchProtocols...),
	}
	// spec.Syscalls = expandSyscalls(d.Get("syscalls").([]interface{}))

	return spec
}

func flattenHostPolicy(policy *kcV1.KubeArmorHostPolicy) []interface{} {
	pol := make(map[string]interface{})
	pol["name"] = policy.ObjectMeta.Name
	pol["namespace"] = policy.ObjectMeta.Namespace
	pol["action"] = policy.Spec.Action
	pol["severity"] = policy.Spec.Severity
	pol["file"] = flattenFile(policy.Spec.File)
	pol["process"] = flattenProcess(policy.Spec.Process)
	pol["capabilities"] = flattenHostCapabilities(policy.Spec.Capabilities)
	pol["network"] = flattenHostNetworks(policy.Spec.Network)
	// pol["syscalls"] = flattenSyscalls(policy.Spec.Syscalls)
	pol["node_selector"] = flattenHostSelector(policy.Spec.NodeSelector)

	return []interface{}{pol}
}

func flattenHostCapabilities(capabilities kcV1.HostCapabilitiesType) []interface{} {
	capabilitiesType := make(map[string]interface{})
	capabilitiesType["match_capabilities"] = flattenMatchHostCapabilities(capabilities.MatchCapabilities)

	return []interface{}{capabilitiesType}
}

func flattenMatchHostCapabilities(in []kcV1.MatchHostCapabilitiesType) []interface{} {
	matchCapabilities := make([]interface{}, len(in))
	for i, match_capabilities := range in {
		m := make(map[string]interface{})
		m["capabilities"] = match_capabilities.Capability
		m["action"] = match_capabilities.Action
		m["from_source"] = flattenFromSource(match_capabilities.FromSource)
		matchCapabilities[i] = m
	}
	return matchCapabilities
}

func flattenHostNetworks(network kcV1.HostNetworkType) []interface{} {
	networkType := make(map[string]interface{})
	networkType["match_protocols"] = flattenMatchHostProtocol(network.MatchProtocols)

	return []interface{}{networkType}
}

func flattenMatchHostProtocol(in []kcV1.MatchHostNetworkProtocolType) []interface{} {
	matchProtocol := make([]interface{}, len(in))
	for i, match_protocols := range in {
		m := make(map[string]interface{})
		m["protocol"] = match_protocols.Protocol
		m["from_source"] = flattenFromSource(match_protocols.FromSource)
		matchProtocol[i] = m
	}
	return matchProtocol
}

func flattenHostSelector(in kcV1.NodeSelectorType) []interface{} {
	node_selectorType := make(map[string]interface{})
	node_selectorType["match_labels"] = flattenMatchLabels(in.MatchLabels)

	return []interface{}{node_selectorType}
}
