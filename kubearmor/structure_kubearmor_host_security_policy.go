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
	matchCapabilityType := kcV1.MatchHostCapabilitiesType{}
	matchCapabilityType.Capability = kcV1.MatchCapabilitiesStringType(matchCapability["capabilities"].(string))
	matchCapabilityType.Severity = kcV1.SeverityType(matchCapability["severity"].(int))
	matchCapabilityType.Action = kcV1.ActionType(matchCapability["action"].(string))
	matchCapabilityType.Tags = expandTags(matchCapability["tags"].([]interface{}))
	matchCapabilityType.Message = matchCapability["message"].(string)
	matchCapabilityType.FromSource = expandFromSource(matchCapability["from_source"].([]interface{}))

	return matchCapabilityType
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
		capabilitiesType.Severity = kcV1.SeverityType(in["severity"].(int))
		capabilitiesType.Action = kcV1.ActionType(in["action"].(string))
		capabilitiesType.Tags = expandTags(in["tags"].([]interface{}))
		capabilitiesType.Message = in["message"].(string)
	}

	return capabilitiesType
}

func expandHostProtocol(matchProtocol map[string]interface{}) kcV1.MatchHostNetworkProtocolType {
	matchProtocolType := kcV1.MatchHostNetworkProtocolType{}
	matchProtocolType.Protocol = kcV1.MatchNetworkProtocolStringType(matchProtocol["protocol"].(string))
	matchProtocolType.Action = kcV1.ActionType(matchProtocol["action"].(string))
	matchProtocolType.Severity = kcV1.SeverityType(matchProtocol["severity"].(int))
	matchProtocolType.Tags = expandTags(matchProtocol["tags"].([]interface{}))
	matchProtocolType.Message = matchProtocol["message"].(string)
	matchProtocolType.FromSource = expandFromSource(matchProtocol["from_source"].([]interface{}))

	return matchProtocolType
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
		networksType.Severity = kcV1.SeverityType(in["severity"].(int))
		networksType.Action = kcV1.ActionType(in["action"].(string))
		networksType.Tags = expandTags(in["tags"].([]interface{}))
		networksType.Message = in["message"].(string)
	}

	return networksType
}

func expandHostSpec(d *schema.ResourceData) kcV1.KubeArmorHostPolicySpec {
	spec := kcV1.KubeArmorHostPolicySpec{}

	spec.Severity = kcV1.SeverityType(d.Get("severity").(int))
	spec.Action = kcV1.ActionType(d.Get("action").(string))
	spec.Tags = expandTags(d.Get("tags").([]interface{}))
	spec.Message = d.Get("message").(string)
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
	spec.Syscalls = expandSyscalls(d.Get("syscalls").([]interface{}))

	return spec
}

func flattenHostPolicy(policy *kcV1.KubeArmorHostPolicy) []interface{} {
	pol := make(map[string]interface{})
	pol["name"] = policy.ObjectMeta.Name
	pol["namespace"] = policy.ObjectMeta.Namespace
	pol["action"] = policy.Spec.Action
	pol["severity"] = policy.Spec.Severity
	pol["tags"] = policy.Spec.Tags
	pol["message"] = policy.Spec.Message
	pol["file"] = flattenFile(policy.Spec.File)
	pol["process"] = flattenProcess(policy.Spec.Process)
	pol["capabilities"] = flattenHostCapabilities(policy.Spec.Capabilities)
	pol["network"] = flattenHostNetworks(policy.Spec.Network)
	pol["syscalls"] = flattenSyscalls(policy.Spec.Syscalls)
	pol["node_selector"] = flattenHostSelector(policy.Spec.NodeSelector)

	return []interface{}{pol}
}

func flattenHostCapabilities(capabilities kcV1.HostCapabilitiesType) []interface{} {
	capabilitiesType := make(map[string]interface{})
	capabilitiesType["match_capabilities"] = flattenMatchHostCapabilities(capabilities.MatchCapabilities)

	capabilitiesType["action"] = capabilities.Action
	capabilitiesType["severity"] = capabilities.Severity
	capabilitiesType["tags"] = capabilities.Tags
	capabilitiesType["message"] = capabilities.Message

	return []interface{}{capabilitiesType}
}

func flattenMatchHostCapabilities(in []kcV1.MatchHostCapabilitiesType) []interface{} {
	matchCapabilities := make([]interface{}, len(in))
	for i, match_capabilities := range in {
		m := make(map[string]interface{})
		m["capabilities"] = match_capabilities.Capability
		m["action"] = match_capabilities.Action
		m["severity"] = match_capabilities.Severity
		m["tags"] = match_capabilities.Tags
		m["message"] = match_capabilities.Message
		m["from_source"] = flattenFromSource(match_capabilities.FromSource)
		matchCapabilities[i] = m
	}
	return matchCapabilities
}

func flattenHostNetworks(network kcV1.HostNetworkType) []interface{} {
	networkType := make(map[string]interface{})
	networkType["match_protocols"] = flattenMatchHostProtocol(network.MatchProtocols)

	networkType["action"] = network.Action
	networkType["severity"] = network.Severity
	networkType["tags"] = network.Tags
	networkType["message"] = network.Message

	return []interface{}{networkType}
}

func flattenMatchHostProtocol(in []kcV1.MatchHostNetworkProtocolType) []interface{} {
	matchProtocol := make([]interface{}, len(in))
	for i, match_protocols := range in {
		m := make(map[string]interface{})
		m["protocol"] = match_protocols.Protocol
		m["from_source"] = flattenFromSource(match_protocols.FromSource)

		m["action"] = match_protocols.Action
		m["severity"] = match_protocols.Severity
		m["tags"] = match_protocols.Tags
		m["message"] = match_protocols.Message

		matchProtocol[i] = m
	}
	return matchProtocol
}

func flattenHostSelector(in kcV1.NodeSelectorType) []interface{} {
	node_selectorType := make(map[string]interface{})
	node_selectorType["match_labels"] = flattenMatchLabels(in.MatchLabels)

	return []interface{}{node_selectorType}
}
