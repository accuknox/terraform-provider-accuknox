# Terraform-Provider-AccuKnox

The AccuKnox terraform provider allows managing KubeArmor resources on kubernetes cluster or host environment.

## Requirements

* [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) v1.0.0
* [Go](https://go.dev/doc/install) 1.20 (to build the provider plugin)

## Building the provider

* Build the plugin using `go build -o terraform-provider-accuknox`
* Copy provider executable to plugins directory, for Linux-system - `~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}`.
```
cp terraform-provider-accuknox ~/.terraform.d/plugins/terraform.example.com/local/accuknox/1.0.0/linux_amd64
```

## Using the provider

* Create `main.tf` config file with resources and data sources
* Execute [`terraform init`](https://developer.hashicorp.com/terraform/cli/commands/init) command to install the provider.
* Then run [`terraform plan`](https://developer.hashicorp.com/terraform/cli/commands/plan) to create an execution plan.
* Finally run [`terraform apply`] to apply the config file with defined resources.

### Defining Provider

```
terraform {
  required_providers {
    accuknox = {
      source = "terraform.example.com/local/accuknox"
      version = "1.0.0"
    }
  }
}

provider "accuknox" {
}
```

## Example - KubeArmor Resources

### Managing KubeArmor Security Policy

There are ways to define `accuknox_kubearmor_security_policy` resource:

1. using yaml format

```
resource "accuknox_kubearmor_security_policy" "ksp-ubuntu-5-net-icmp-audit" {
  policy= <<-EOT
  apiVersion: security.kubearmor.com/v1
  kind: KubeArmorPolicy
  metadata:
    name: ksp-ubuntu-5-net-icmp-audit
  spec:
    severity: 8
    selector:
      matchLabels:
        container: ubuntu-5
    network:
      matchProtocols:
      - protocol: icmp
    action:
      Audit  
  EOT
}
```
2. using defined schema format

```
resource "accuknox_kubearmor_security_policy" "ksp-vault-protect" {
  name="ksp-vault-protect"
  namespace= "default"
  severity= 7
  action= "Allow"
  
  selector {
    match_labels= {
      "app.kubernetes.io/name": "vault",
      "component": "server"
    }
  }
  
  file {
    match_directories{
        dir= "/vault/"
        recursive= true
        action= "Block"
      }
    match_directories{
        dir= "/"
        recursive= true
      }
    match_directories{
        dir= "/vault/"
        recursive= true
        from_source{
            path= "/bin/vault"
          }
      }
  }

  process {
    match_paths{
        path= "/bin/busybox"
      }
    match_paths{
        path= "/bin/vault"
      }
  }
}
```

> In this format currently `Syscalls` are not supported due to an [issue](https://github.com/kubearmor/KubeArmor/issues/1332).

### Managing KubeArmor Host Security Policy

There are ways to define `accuknox_kubearmor_host_security_policy` resource:

1. using yaml format

```
resource "accuknox_kubearmor_host_security_policy" "hsp-kubearmor-dev-proc-path-block" {
  policy= <<-EOT
  apiVersion: security.kubearmor.com/v1
  kind: KubeArmorHostPolicy
  metadata:
    name: hsp-kubearmor-dev-proc-path-block
  spec:
    nodeSelector:
      matchLabels:
        kubernetes.io/hostname: kubearmor-dev
    severity: 5
    process:
      matchPaths:
      - path: /usr/bin/diff
    action:
      Block
  EOT
}
```

2. using defined schema format

```
resource "accuknox_kubearmor_host_security_policy" "hsp-kubearmor-dev-file-path-audit" {
  name="hsp-kubearmor-dev-file-path-audit"
  severity= 5
  action= "Audit"
  
  node_selector {
    match_labels= {
      "kubernetes.io/hostname": "kubearmor-dev",
    }
  }
  
  file {
    match_paths{
        path= "/etc/passwd"
      }
  }
}
```

### Managing KubeArmor Configuration

```
resource "accuknox_kubearmor_configuration" "conf" {
  name="kubearmor-config"
  namespace="kube-system"
  data={
    "defaultCapabilitiesPosture"="audit",
  }
}
```

### Managing Namsepace Posture

```
resource "accuknox_kubearmor_namespace_posture" "ns_pos" {
  namespace     = "wordpress-mysql"
  file          = "block"
  network       = "block"
  capabilities  = "audit"
}
```

### Managing Namsepace Visibility

```
resource "accuknox_kubearmor_namespace_visibility" "visib-ns1" {
  namespace     = "wordpress-mysql"
  file          = true
  network       = true
  capabilities  = true
  process       = true
}
```

## Example - KubeArmor Data Sources

### Read KubeArmor Configuration

```
data "accuknox_kubearmor_configuration" "data_cm" {
  name="kubearmor-config"
  namespace="kube-system"
}

output "data_cm" {
  value = data.kubearmor_configuration.data_cm.data
}
```

### Read KubeArmor Security Policy

```
data "accuknox_kubearmor_security_policy" "pkg-mgmt" {
  name="block-pkg-mgmt-tools-exec"
  namespace="default"
}

output "sp" {
  value = data.kubearmor_security_policy.pkg-mgmt.policy
}
```

### Read KubeArmor Host Security Policy

```
data "accuknox_kubearmor_host_security_policy" "host-policy" {
  name="hsp-kubearmor-dev-proc-path-block"
}

output "host-policy" {
  value = data.kubearmor_host_security_policy.host-policy.policy
}
```

### Read KubeArmor Installed Version

```
data "accuknox_kubearmor_installed_version" "installed_version" {}

output "installed_version" {
  value = data.kubearmor_installed_version.installed_version.version
}
```

### Read KubeArmor Node Information

```
data "accuknox_kubearmor_node" "k_node" {}

output "k_node" {
  value = data.kubearmor_node.k_node.node_data
}
```

### Read Namespace Posture

```
data "accuknox_kubearmor_namespace_posture" "ns_ps" {
  name="kube-system"
}

output "ns_ps" {
  value = data.kubearmor_namespace_posture.ns_ps.annotation
}
```

### Read Namespace Visibility

```
data "accuknox_kubearmor_namespace_visibility" "ns_vs" {
  name="kube-system"
}

output "ns_vs" {
  value = data.kubearmor_namespace_visibility.ns_vs.visibility
}
```

### Read KubeArmror Stable Version

```
data "accuknox_kubearmor_stable_version" "stable_version" {}

output "stable_version" {
  value = data.kubearmor_stable_version.stable_version.version
}
```