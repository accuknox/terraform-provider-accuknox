package kubearmor

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kcV1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func expandfs(fs map[string]interface{}) kcV1.MatchSourceType {
	mst := kcV1.MatchSourceType{}
	mst.Path = kcV1.MatchPathType(fs["path"].(string))
	return mst
}

func expandFromSource(fromSource []interface{}) []kcV1.MatchSourceType {
	fromSourceType := []kcV1.MatchSourceType{}

	for _, fs := range fromSource {
		fromSourceType = append(fromSourceType, expandfs(fs.(map[string]interface{})))
	}

	return fromSourceType
}

func expandMatchFilePath(matchPath map[string]interface{}) kcV1.FilePathType {
	matchPathType := kcV1.FilePathType{}
	matchPathType.Path = kcV1.MatchPathType(matchPath["path"].(string))
	matchPathType.OwnerOnly = matchPath["owner_only"].(bool)
	matchPathType.ReadOnly = matchPath["read_omly"].(bool)
	matchPathType.Action = kcV1.ActionType(matchPath["action"].(string))
	matchPathType.Severity = kcV1.SeverityType(matchPath["severity"].(int))
	matchPathType.Tags = expandTags(matchPath["tags"].([]interface{}))
	matchPathType.Message = matchPath["message"].(string)
	matchPathType.FromSource = expandFromSource(matchPath["from_source"].([]interface{}))

	return matchPathType
}

func expandMatchFilePaths(matchPaths []interface{}) []kcV1.FilePathType {
	if len(matchPaths) < 1 {
		return nil
	}
	matchPathsType := []kcV1.FilePathType{}

	for _, matchPath := range matchPaths {
		matchPathsType = append(matchPathsType, expandMatchFilePath(matchPath.(map[string]interface{})))
	}

	return matchPathsType
}

func expandMatchFileDir(matchDir map[string]interface{}) kcV1.FileDirectoryType {
	matchDirType := kcV1.FileDirectoryType{}
	matchDirType.Directory = kcV1.MatchDirectoryType(matchDir["dir"].(string))
	matchDirType.OwnerOnly = matchDir["owner_only"].(bool)
	matchDirType.ReadOnly = matchDir["read_only"].(bool)
	matchDirType.Recursive = matchDir["recursive"].(bool)
	matchDirType.Action = kcV1.ActionType(matchDir["action"].(string))
	matchDirType.Severity = kcV1.SeverityType(matchDir["severity"].(int))
	matchDirType.Tags = expandTags(matchDir["tags"].([]interface{}))
	matchDirType.Message = matchDir["message"].(string)

	matchDirType.FromSource = expandFromSource(matchDir["from_source"].([]interface{}))

	return matchDirType
}

func expandMatchFileDirectories(matchDirectories []interface{}) []kcV1.FileDirectoryType {
	if len(matchDirectories) < 1 {
		return []kcV1.FileDirectoryType{}
	}
	matchdirectories := []kcV1.FileDirectoryType{}

	for _, matchDirectory := range matchDirectories {
		matchdirectories = append(matchdirectories, expandMatchFileDir(matchDirectory.(map[string]interface{})))
	}

	return matchdirectories
}

func expandMatchFilePattern(matchPattern map[string]interface{}) kcV1.FilePatternType {
	matchPatternType := kcV1.FilePatternType{}
	matchPatternType.Pattern = matchPattern["pattern"].(string)
	matchPatternType.Action = kcV1.ActionType(matchPattern["action"].(string))
	matchPatternType.Severity = kcV1.SeverityType(matchPattern["severity"].(int))
	matchPatternType.Tags = expandTags(matchPattern["tags"].([]interface{}))
	matchPatternType.Message = matchPattern["message"].(string)

	matchPatternType.OwnerOnly = matchPattern["owner_only"].(bool)

	return matchPatternType
}

func expandMatchFilePatterns(matchPatterns []interface{}) []kcV1.FilePatternType {
	matchpatterns := []kcV1.FilePatternType{}

	for _, matchPath := range matchPatterns {
		matchpatterns = append(matchpatterns, expandMatchFilePattern(matchPath.(map[string]interface{})))
	}

	return matchpatterns
}

func expandFile(file []interface{}) kcV1.FileType {
	fileType := kcV1.FileType{}
	if len(file) > 0 {
		in := file[0].(map[string]interface{})
		fileType.MatchPaths = expandMatchFilePaths(in["match_paths"].([]interface{}))
		fileType.MatchDirectories = expandMatchFileDirectories(in["match_directories"].([]interface{}))
		fileType.MatchPatterns = expandMatchFilePatterns(in["match_patterns"].([]interface{}))

		fileType.Severity = kcV1.SeverityType(in["severity"].(int))
		fileType.Action = kcV1.ActionType(in["action"].(string))
		fileType.Tags = expandTags(in["tags"].([]interface{}))
		fileType.Message = in["message"].(string)
	}

	return fileType
}

func expandMatchProcessPath(matchPath map[string]interface{}) kcV1.ProcessPathType {
	matchPathType := kcV1.ProcessPathType{}
	matchPathType.Path = kcV1.MatchPathType(matchPath["path"].(string))
	matchPathType.OwnerOnly = matchPath["owner_only"].(bool)
	matchPathType.Severity = kcV1.SeverityType(matchPath["severity"].(int))
	matchPathType.Action = kcV1.ActionType(matchPath["action"].(string))
	matchPathType.Tags = expandTags(matchPath["tags"].([]interface{}))
	matchPathType.Message = matchPath["message"].(string)

	matchPathType.FromSource = expandFromSource(matchPath["from_source"].([]interface{}))

	return matchPathType
}

func expandMatchProcessDir(matchDir map[string]interface{}) kcV1.ProcessDirectoryType {
	matchDirType := kcV1.ProcessDirectoryType{}
	matchDirType.Directory = kcV1.MatchDirectoryType(matchDir["dir"].(string))
	matchDirType.OwnerOnly = matchDir["owner_only"].(bool)
	matchDirType.Recursive = matchDir["recursive"].(bool)
	matchDirType.Action = kcV1.ActionType(matchDir["action"].(string))
	matchDirType.Severity = kcV1.SeverityType(matchDir["severity"].(int))
	matchDirType.Tags = expandTags(matchDir["tags"].([]interface{}))
	matchDirType.Message = matchDir["message"].(string)

	matchDirType.FromSource = expandFromSource(matchDir["from_source"].([]interface{}))

	return matchDirType
}

func expandMatchProcessPattern(matchPattern map[string]interface{}) kcV1.ProcessPatternType {
	matchPatternType := kcV1.ProcessPatternType{}
	matchPatternType.Pattern = matchPattern["pattern"].(string)
	matchPatternType.Action = kcV1.ActionType(matchPattern["action"].(string))
	matchPatternType.Severity = kcV1.SeverityType(matchPattern["severity"].(int))
	matchPatternType.Tags = expandTags(matchPattern["tags"].([]interface{}))
	matchPatternType.Message = matchPattern["message"].(string)

	matchPatternType.OwnerOnly = matchPattern["owner_only"].(bool)

	return matchPatternType
}

func expandMatchProcessPatterns(matchPatterns []interface{}) []kcV1.ProcessPatternType {
	matchpatterns := []kcV1.ProcessPatternType{}

	for _, matchPath := range matchPatterns {
		matchpatterns = append(matchpatterns, expandMatchProcessPattern(matchPath.(map[string]interface{})))
	}

	return matchpatterns
}

func expandMatchProcessDirectories(matchPaths []interface{}) []kcV1.ProcessDirectoryType {
	matchdirectories := []kcV1.ProcessDirectoryType{}

	for _, matchPath := range matchPaths {
		matchdirectories = append(matchdirectories, expandMatchProcessDir(matchPath.(map[string]interface{})))
	}

	return matchdirectories
}

func expandMatchProcessPaths(matchPaths []interface{}) []kcV1.ProcessPathType {
	matchpaths := []kcV1.ProcessPathType{}

	for _, matchPath := range matchPaths {
		matchpaths = append(matchpaths, expandMatchProcessPath(matchPath.(map[string]interface{})))
	}

	return matchpaths
}
func expandProcess(process []interface{}) kcV1.ProcessType {
	processType := kcV1.ProcessType{}
	if len(process) > 0 {
		in := process[0].(map[string]interface{})
		processType.MatchPaths = expandMatchProcessPaths(in["match_paths"].([]interface{}))
		processType.MatchDirectories = expandMatchProcessDirectories(in["match_directories"].([]interface{}))
		processType.MatchPatterns = expandMatchProcessPatterns(in["match_patterns"].([]interface{}))

		processType.Severity = kcV1.SeverityType(in["severity"].(int))
		processType.Action = kcV1.ActionType(in["action"].(string))
		processType.Tags = expandTags(in["tags"].([]interface{}))
		processType.Message = in["message"].(string)

	}

	return processType
}

func expandCapability(matchCapability map[string]interface{}) kcV1.MatchCapabilitiesType {
	matchCapabilityType := kcV1.MatchCapabilitiesType{}
	matchCapabilityType.Capability = kcV1.MatchCapabilitiesStringType(matchCapability["capabilities"].(string))
	matchCapabilityType.Action = kcV1.ActionType(matchCapability["action"].(string))
	matchCapabilityType.Severity = kcV1.SeverityType(matchCapability["severity"].(int))
	matchCapabilityType.Tags = expandTags(matchCapability["tags"].([]interface{}))
	matchCapabilityType.Message = matchCapability["message"].(string)

	matchCapabilityType.FromSource = expandFromSource(matchCapability["from_source"].([]interface{}))

	return matchCapabilityType
}

func expandMatchCapabilities(matchCapabilities []interface{}) []kcV1.MatchCapabilitiesType {
	matchcapabilities := []kcV1.MatchCapabilitiesType{}

	for _, matchCapability := range matchCapabilities {
		matchcapabilities = append(matchcapabilities, expandCapability(matchCapability.(map[string]interface{})))
	}

	return matchcapabilities
}

func expandCapabilities(capabilities []interface{}) kcV1.CapabilitiesType {
	capabilitiesType := kcV1.CapabilitiesType{}
	if len(capabilities) > 0 {
		in := capabilities[0].(map[string]interface{})
		capabilitiesType.MatchCapabilities = expandMatchCapabilities(in["match_paths"].([]interface{}))
		capabilitiesType.Severity = kcV1.SeverityType(in["severity"].(int))
		capabilitiesType.Action = kcV1.ActionType(in["action"].(string))
		capabilitiesType.Tags = expandTags(in["tags"].([]interface{}))
		capabilitiesType.Message = in["message"].(string)

	}

	return capabilitiesType
}

func expandProtocol(matchProtocol map[string]interface{}) kcV1.MatchNetworkProtocolType {
	matchProtocolType := kcV1.MatchNetworkProtocolType{}
	matchProtocolType.Protocol = kcV1.MatchNetworkProtocolStringType(matchProtocol["protocol"].(string))
	matchProtocolType.FromSource = expandFromSource(matchProtocol["from_source"].([]interface{}))
	matchProtocolType.Severity = kcV1.SeverityType(matchProtocol["severity"].(int))
	matchProtocolType.Action = kcV1.ActionType(matchProtocol["action"].(string))
	matchProtocolType.Tags = expandTags(matchProtocol["tags"].([]interface{}))
	matchProtocolType.Message = matchProtocol["message"].(string)

	return matchProtocolType
}

func expandMatchProtocols(protocol []interface{}) []kcV1.MatchNetworkProtocolType {
	matchnetwork := []kcV1.MatchNetworkProtocolType{}

	for _, matchProtocol := range protocol {
		matchnetwork = append(matchnetwork, expandProtocol(matchProtocol.(map[string]interface{})))
	}

	return matchnetwork
}

func expandNetwork(network []interface{}) kcV1.NetworkType {
	networksType := kcV1.NetworkType{}
	if len(network) > 0 {
		in := network[0].(map[string]interface{})
		networksType.MatchProtocols = expandMatchProtocols(in["match_protocols"].([]interface{}))
		networksType.Severity = kcV1.SeverityType(in["severity"].(int))
		networksType.Action = kcV1.ActionType(in["action"].(string))
		networksType.Tags = expandTags(in["tags"].([]interface{}))
		networksType.Message = in["message"].(string)

	}

	return networksType
}

func expandsysfs(fs map[string]interface{}) kcV1.SyscallFromSourceType {
	mst := kcV1.SyscallFromSourceType{}
	mst.Path = kcV1.MatchPathType(fs["path"].(string))
	mst.Dir = fs["dir"].(string)
	mst.Recursive = fs["recursive"].(bool)
	return mst
}

func expandSyscallFromSource(fromSource []interface{}) []kcV1.SyscallFromSourceType {
	fromSourceType := []kcV1.SyscallFromSourceType{}

	for _, fs := range fromSource {
		fromSourceType = append(fromSourceType, expandsysfs(fs.(map[string]interface{})))
	}

	return fromSourceType
}
func expandsyscall(in []interface{}) []kcV1.Syscall {
	syscall := make([]kcV1.Syscall, len(in))
	for i, v := range in {
		syscall[i] = kcV1.Syscall(v.(string))
	}
	return syscall
}

func expandMatchSyscall(matchSyscall map[string]interface{}) kcV1.SyscallMatchType {
	matchSyscallType := kcV1.SyscallMatchType{}
	matchSyscallType.Syscalls = expandsyscall(matchSyscall["syscall"].([]interface{}))
	matchSyscallType.FromSource = expandSyscallFromSource(matchSyscall["from_source"].([]interface{}))
	return matchSyscallType
}

func expandMatchSyscalls(MatchSyscalls []interface{}) []kcV1.SyscallMatchType {
	matchsyscalls := []kcV1.SyscallMatchType{}

	for _, matchSyscall := range MatchSyscalls {
		matchsyscalls = append(matchsyscalls, expandMatchSyscall(matchSyscall.(map[string]interface{})))
	}

	return matchsyscalls
}

func expandMatchPath(matchPath map[string]interface{}) kcV1.SyscallMatchPathType {
	matchPathType := kcV1.SyscallMatchPathType{}
	matchPathType.Path = kcV1.MatchSyscallPathType(matchPath["path"].(string))
	matchPathType.Recursive = matchPath["path"].(bool)
	matchPathType.Syscalls = expandsyscall(matchPath["syscall"].([]interface{}))

	matchPathType.FromSource = expandSyscallFromSource(matchPath["from_source"].([]interface{}))

	return matchPathType
}

func expandSyscallsMatchPaths(in []interface{}) []kcV1.SyscallMatchPathType {
	matchpaths := []kcV1.SyscallMatchPathType{}

	for _, matchPaths := range in {
		matchpaths = append(matchpaths, expandMatchPath(matchPaths.(map[string]interface{})))
	}

	return matchpaths
}

func expandSyscalls(syscalls []interface{}) kcV1.SyscallsType {
	syscallType := kcV1.SyscallsType{}
	if len(syscalls) > 0 {
		in := syscalls[0].(map[string]interface{})
		syscallType.MatchSyscalls = expandMatchSyscalls(in["match_syscalls"].([]interface{}))
		syscallType.MatchPaths = expandSyscallsMatchPaths(in["match_paths"].([]interface{}))

		syscallType.Message = in["message"].(string)
		syscallType.Severity = kcV1.SeverityType(in["severity"].(int))
		syscallType.Tags = expandTags(in["tags"].([]interface{}))
	}

	return syscallType
}

func expandSelector(selector []interface{}) kcV1.SelectorType {
	selectorType := kcV1.SelectorType{}
	in := selector[0].(map[string]interface{})
	if v, ok := in["match_labels"].(map[string]interface{}); ok && len(v) > 0 {
		selectorType.MatchLabels = expandStringMap(v)
	}
	return selectorType
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func expandTags(in []interface{}) []string {
	tags := make([]string, len(in))
	for i, v := range in {
		tags[i] = v.(string)
	}
	return tags
}

func expandSpec(d *schema.ResourceData) kcV1.KubeArmorPolicySpec {
	spec := kcV1.KubeArmorPolicySpec{}

	spec.Severity = kcV1.SeverityType(d.Get("severity").(int))
	spec.Action = kcV1.ActionType(d.Get("action").(string))
	spec.Tags = expandTags(d.Get("tags").([]interface{}))
	spec.Message = d.Get("message").(string)
	spec.Selector = expandSelector(d.Get("selector").([]interface{}))
	spec.File = expandFile(d.Get("file").([]interface{}))
	spec.Process = expandProcess(d.Get("process").([]interface{}))

	spec.Capabilities = expandCapabilities(d.Get("capabilities").([]interface{}))
	spec.Capabilities = kcV1.CapabilitiesType{
		MatchCapabilities: append([]kcV1.MatchCapabilitiesType{}, spec.Capabilities.MatchCapabilities...),
	}

	spec.Network = expandNetwork(d.Get("network").([]interface{}))
	spec.Network = kcV1.NetworkType{
		MatchProtocols: append([]kcV1.MatchNetworkProtocolType{}, spec.Network.MatchProtocols...),
	}
	spec.Syscalls = expandSyscalls(d.Get("syscalls").([]interface{}))

	return spec
}

func flattenFromSource(in []kcV1.MatchSourceType) []interface{} {
	fromSource := make([]interface{}, len(in))
	for i, fs := range in {
		m := make(map[string]interface{})
		m["path"] = fs.Path
		fromSource[i] = m
	}
	return fromSource
}

func flattenMatchFilePaths(in []kcV1.FilePathType) []interface{} {
	matchPaths := make([]interface{}, len(in))
	for i, filePath := range in {
		m := make(map[string]interface{})
		m["path"] = filePath.Path
		m["read_only"] = filePath.ReadOnly
		m["owner_only"] = filePath.OwnerOnly
		m["action"] = filePath.Action
		m["severity"] = filePath.Severity
		m["tags"] = filePath.Tags
		m["message"] = filePath.Message
		m["from_source"] = flattenFromSource(filePath.FromSource)
		matchPaths[i] = m
	}
	return matchPaths
}

func flattenMatchFileDirectories(in []kcV1.FileDirectoryType) []interface{} {
	matchPaths := make([]interface{}, len(in))
	for i, filePath := range in {
		m := make(map[string]interface{})
		m["dir"] = filePath.Directory
		m["read_only"] = filePath.ReadOnly
		m["owner_only"] = filePath.OwnerOnly
		m["action"] = filePath.Action
		m["severity"] = filePath.Severity
		m["tags"] = filePath.Tags
		m["message"] = filePath.Message
		m["recursive"] = filePath.Recursive
		m["from_source"] = flattenFromSource(filePath.FromSource)
		matchPaths[i] = m
	}
	return matchPaths
}

func flattenMatchFilePatterns(in []kcV1.FilePatternType) []interface{} {
	matchPatterns := make([]interface{}, len(in))
	for i, filePattern := range in {
		m := make(map[string]interface{})
		m["pattern"] = filePattern.Pattern
		m["read_only"] = filePattern.ReadOnly
		m["owner_only"] = filePattern.OwnerOnly
		m["action"] = filePattern.Action
		m["severity"] = filePattern.Severity
		m["tags"] = filePattern.Tags
		m["message"] = filePattern.Message
		matchPatterns[i] = m
	}
	return matchPatterns
}

func flattenFile(file kcV1.FileType) []interface{} {
	fileType := make(map[string]interface{})
	fileType["match_paths"] = flattenMatchFilePaths(file.MatchPaths)
	fileType["match_directories"] = flattenMatchFileDirectories(file.MatchDirectories)
	fileType["match_patterns"] = flattenMatchFilePatterns(file.MatchPatterns)

	fileType["action"] = file.Action
	fileType["severity"] = file.Severity
	fileType["tags"] = file.Tags
	fileType["message"] = file.Message

	return []interface{}{fileType}
}

func flattenMatchProcessPaths(in []kcV1.ProcessPathType) []interface{} {
	matchPaths := make([]interface{}, len(in))
	for i, filePath := range in {
		m := make(map[string]interface{})
		m["path"] = filePath.Path
		m["owner_only"] = filePath.OwnerOnly
		m["action"] = filePath.Action
		m["severity"] = filePath.Severity
		m["tags"] = filePath.Tags
		m["message"] = filePath.Message
		m["from_source"] = flattenFromSource(filePath.FromSource)
		matchPaths[i] = m
	}
	return matchPaths
}

func flattenMatchProcessDirectories(in []kcV1.ProcessDirectoryType) []interface{} {
	matchPaths := make([]interface{}, len(in))
	for i, filePath := range in {
		m := make(map[string]interface{})
		m["path"] = filePath.Directory
		m["owner_only"] = filePath.OwnerOnly
		m["action"] = filePath.Action
		m["severity"] = filePath.Severity
		m["tags"] = filePath.Tags
		m["message"] = filePath.Message
		m["from_source"] = flattenFromSource(filePath.FromSource)
		matchPaths[i] = m
	}
	return matchPaths
}

func flattenMatchProcessPatterns(in []kcV1.ProcessPatternType) []interface{} {
	matchPatterns := make([]interface{}, len(in))
	for i, processPattern := range in {
		m := make(map[string]interface{})
		m["pattern"] = processPattern.Pattern
		m["owner_only"] = processPattern.OwnerOnly
		m["action"] = processPattern.Action
		m["severity"] = processPattern.Severity
		m["tags"] = processPattern.Tags
		m["message"] = processPattern.Message
		matchPatterns[i] = m
	}
	return matchPatterns
}

func flattenProcess(process kcV1.ProcessType) []interface{} {
	processType := make(map[string]interface{})
	processType["match_paths"] = flattenMatchProcessPaths(process.MatchPaths)
	processType["match_directories"] = flattenMatchProcessDirectories(process.MatchDirectories)
	processType["match_patterns"] = flattenMatchProcessPatterns(process.MatchPatterns)

	processType["action"] = process.Action
	processType["severity"] = process.Severity
	processType["tags"] = process.Tags
	processType["message"] = process.Message

	return []interface{}{processType}
}

func flattenMatchCapabilities(in []kcV1.MatchCapabilitiesType) []interface{} {
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

func flattenCapabilities(capabilities kcV1.CapabilitiesType) []interface{} {
	capabilitiesType := make(map[string]interface{})
	capabilitiesType["match_capabilities"] = flattenMatchCapabilities(capabilities.MatchCapabilities)

	capabilitiesType["action"] = capabilities.Action
	capabilitiesType["severity"] = capabilities.Severity
	capabilitiesType["tags"] = capabilities.Tags
	capabilitiesType["message"] = capabilities.Message

	return []interface{}{capabilitiesType}
}

func flattenMatchProtocol(in []kcV1.MatchNetworkProtocolType) []interface{} {
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

func flattenNetworks(network kcV1.NetworkType) []interface{} {
	networkType := make(map[string]interface{})
	networkType["match_protocols"] = flattenMatchProtocol(network.MatchProtocols)

	networkType["action"] = network.Action
	networkType["severity"] = network.Severity
	networkType["tags"] = network.Tags
	networkType["message"] = network.Message

	return []interface{}{networkType}
}

func flattenSyscallFromSource(in []kcV1.SyscallFromSourceType) []interface{} {
	fromSource := make([]interface{}, len(in))
	for i, fs := range in {
		m := make(map[string]interface{})
		m["path"] = fs.Path
		m["dir"] = fs.Dir
		m["recursive"] = fs.Recursive
		fromSource[i] = m
	}
	return fromSource
}

func flattenMatchSyscalls(in []kcV1.SyscallMatchType) []interface{} {
	SyscallMatchType := make([]interface{}, len(in))
	for i, match_syscalls := range in {
		m := make(map[string]interface{})
		m["syscall"] = match_syscalls.Syscalls
		m["from_source"] = flattenSyscallFromSource(match_syscalls.FromSource)
		SyscallMatchType[i] = m
	}
	return SyscallMatchType
}

func flattenMatchSyscallPaths(in []kcV1.SyscallMatchPathType) []interface{} {
	SyscallMatchPathType := make([]interface{}, len(in))
	for i, match_syscalls_paths := range in {
		m := make(map[string]interface{})
		m["path"] = match_syscalls_paths.Path
		m["recursive"] = match_syscalls_paths.Recursive
		m["syscall"] = match_syscalls_paths.Syscalls
		m["from_source"] = flattenSyscallFromSource(match_syscalls_paths.FromSource)
		SyscallMatchPathType[i] = m
	}
	return SyscallMatchPathType
}

func flattenSyscalls(syscalls kcV1.SyscallsType) []interface{} {
	SyscallsType := make(map[string]interface{})
	SyscallsType["match_syscalls"] = flattenMatchSyscalls(syscalls.MatchSyscalls)
	SyscallsType["match_paths"] = flattenMatchSyscallPaths(syscalls.MatchPaths)

	SyscallsType["message"] = syscalls.Message
	SyscallsType["severity"] = syscalls.Severity
	SyscallsType["tags"] = syscalls.Tags

	return []interface{}{SyscallsType}
}

func flattenMatchLabels(labels map[string]string) map[string]string {
	label := make(map[string]string)
	for k, v := range labels {
		label[k] = v
	}
	return label

}

func flattenSelector(in kcV1.SelectorType) []interface{} {
	node_selectorType := make(map[string]interface{})
	node_selectorType["match_labels"] = flattenMatchLabels(in.MatchLabels)

	return []interface{}{node_selectorType}
}

func flattenPolicy(policy *kcV1.KubeArmorPolicy) []interface{} {
	pol := make(map[string]interface{})
	pol["name"] = policy.ObjectMeta.Name
	pol["namespace"] = policy.ObjectMeta.Namespace
	pol["action"] = policy.Spec.Action
	pol["severity"] = policy.Spec.Severity
	pol["tags"] = policy.Spec.Tags
	pol["message"] = policy.Spec.Message
	pol["file"] = flattenFile(policy.Spec.File)
	pol["process"] = flattenProcess(policy.Spec.Process)
	pol["capabilities"] = flattenCapabilities(policy.Spec.Capabilities)
	pol["network"] = flattenNetworks(policy.Spec.Network)
	pol["syscalls"] = flattenSyscalls(policy.Spec.Syscalls)
	pol["selector"] = flattenSelector(policy.Spec.Selector)

	return []interface{}{pol}
}

func BuildId(meta metav1.ObjectMeta) string {
	return meta.Namespace + "/" + meta.Name
}

func IdParts(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		err := fmt.Errorf("unexpected ID format (%q), expected %q", id, "namespace/name")
		return "", "", err
	}

	return parts[0], parts[1], nil
}

func objectMeta(d *schema.ResourceData) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: d.Get("name").(string),
	}
}
