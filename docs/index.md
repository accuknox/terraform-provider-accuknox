---
page_title: "Provider: Accuknox"
subcategory: ""
description: |-
  The Accuknox provider is used to interact with the resources supported by KubeArmor and Discovery-Engine.
---

# Accuknox Provider

The Accuknox provider is used to interact with the resources supported by KubeArmor and Discovery-Engine.

Use the navigation to the left to read about the available resources.

## Example Usage

> As per the current version the provider section should be empty and no requirement for credentials to use Accuknox provider.

```terraform

terraform {
  required_providers {
    accuknox = {
      source = "hashicorp/accuknox"
      version = "1.0.0"
    }
  }
}

provider "accuknox" {

}

resource "accuknox_kubearmor_security_policy" "block-pkg-mgmt-tools-exec" {
  policy= <<-EOT
  apiVersion: security.kubearmor.com/v1
  kind: KubeArmorPolicy
  metadata:
    name: block-pkg-mgmt-tools-exec
  spec:
    selector:
      matchLabels:
        app: nginx
    process:
      matchPaths:
      - path: /usr/bin/apt
      - path: /usr/bin/apt-get
    action:
      Block
  EOT
}

```

## KubeArmor

[KubeArmor](https://docs.kubearmor.io/kubearmor/) is a cloud-native runtime security enforcement system that restricts the behavior (such as process execution, file access, and networking operations) of pods, containers, and nodes (VMs) at the system level. KubeArmor leverages Linux security modules (LSMs) such as AppArmor, SELinux, or BPF-LSM to enforce the user-specified policies. KubeArmor generates rich alerts/telemetry events with container/pod/namespace identities by leveraging eBPF.

## Discovery-Engine

Discovery Engine discovers the security posture for your workloads and auto-discovers the policy-set required to put the workload in least-permissive mode. The engine leverages the rich visibility provided by KubeArmor and Cilium to auto discover the systems and network security posture.