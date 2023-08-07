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
> This is required only until we do not publish our provider to official registry.

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

`accuknox_kubearmor_security_policy` will give us the ability to manage the container security policies. Read more about container security policies [here](https://docs.kubearmor.io/kubearmor/specification/security_policy_specification).

There are two ways to define `accuknox_kubearmor_security_policy` resource:

1. using yaml format(recommended way):

Directly use the policy to define it in the yaml format. Since we are applying a policy as an object, updating it would trigger `destroy and then create` operation, therefore new policy will be created.

Here, we are creating/managing a policy which blocks package management tools.

```
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

2. using defined schema format:

Use objects such as `process`, `network`, `file`, etc to define rules of the policy. We can also update the policy, except for the `name` and `namespace`, updating them would trigger `destroy and then create` operation.

Here, we are creating/managing a policy for vault which blocks access to volume mounts.

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

### Managing KubeArmor Host Security Policy

`accuknox_kubearmor_host_security_policy` will give us the ability to manage the host security policies. Read more about host security policies [here](https://docs.kubearmor.io/kubearmor/specification/host_security_policy_specification).

There are two ways to define `accuknox_kubearmor_host_security_policy` resource:

1. using yaml format:

Directly use the policy to define it in the yaml format. Since we are applying a policy as an object, updating it would trigger `destroy and then create` operation, therefore new policy will be created.

Here, we are creating/managing a policy for host which blocks `/usr/bin/diff` path.

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

2. using defined schema format :

Use objects such as `process`, `network`, `file`, etc to define rules of the policy. We can also update the policy, except for the `name` and `namespace`, updating them would trigger `destroy and then create` operation.

Here, we are creating/managing a policy for host which alerts when there is an access to `/etc/passwd` path.

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

We can manage the global [default posture](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/default_posture.md), global [default visibility](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/kubearmor_visibility.md), cluster and grpc port with resource `accuknox_kubearmor_configuration`.

This is the default configuration:

```
data:
  cluster: default
  defaultCapabilitiesPosture: audit
  defaultFilePosture: audit
  defaultNetworkPosture: audit
  gRPC: "32767"
  visibility: process,file,network
```
Here, we are changing `defaultFilePosture` to block.

```
resource "accuknox_kubearmor_configuration" "conf" {
  name="kubearmor-config"
  namespace="kube-system"
  data={
    "defaultFilePosture"="block",
  }
}
```

### Managing Namsepace Posture

We can manage [namespace posture](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/default_posture.md#namespace-default-posture) with resource `accuknox_kubearmor_namespace_posture`.

Here, we are updating posture for namespace `wordpress-mysql`.

```
resource "accuknox_kubearmor_namespace_posture" "ns_pos" {
  namespace     = "wordpress-mysql"
  file          = "block"
  network       = "block"
  capabilities  = "audit"
}
```

### Managing Namsepace Visibility

We can manage [namespace visibility](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/kubearmor_visibility.md#updating-namespace-visibility) with resource `accuknox_kubearmor_namespace_visibility`.

Here, we are updating visibility for namespace `wordpress-mysql`.

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

data source `accuknox_kubearmor_configuration` provides the current kubearmor configuration.

```
data "accuknox_kubearmor_configuration" "data_cm" {
  name="kubearmor-config"
  namespace="kube-system"
}

output "data_cm" {
  value = data.accuknox_kubearmor_configuration.data_cm.data
}
```

### Read KubeArmor Security Policy

data source `accuknox_kubearmor_security_policy` provides the security policy based on `name` and `namespace`.

```
data "accuknox_kubearmor_security_policy" "pkg-mgmt" {
  name="block-pkg-mgmt-tools-exec"
  namespace="default"
}

output "sp" {
  value = data.accuknox_kubearmor_security_policy.pkg-mgmt.policy
}
```

### Read KubeArmor Host Security Policy

data source `accuknox_kubearmor_host_security_policy` provides the security policy based on `name`.

```
data "accuknox_kubearmor_host_security_policy" "host-policy" {
  name="hsp-kubearmor-dev-proc-path-block"
}

output "host-policy" {
  value = data.accuknox_kubearmor_host_security_policy.host-policy.policy
}
```

### Read KubeArmor Installed Version

data source `accuknox_kubearmor_installed_version` provides the current installed verion of kubearmor.

```
data "accuknox_kubearmor_installed_version" "installed_version" {}

output "installed_version" {
  value = data.accuknox_kubearmor_installed_version.installed_version.version
}
```

### Read KubeArmor Node Information

data source `accuknox_kubearmor_node` provides the node information.

```
data "accuknox_kubearmor_node" "k_node" {}

output "k_node" {
  value = data.accuknox_kubearmor_node.k_node.node_data
}
```

### Read Namespace Posture

data source `accuknox_kubearmor_namespace_posture` provides the current namespace based posture information.

```
data "accuknox_kubearmor_namespace_posture" "ns_ps" {
  name="kube-system"
}

output "ns_ps" {
  value = data.accuknox_kubearmor_namespace_posture.ns_ps.annotation
}
```

### Read Namespace Visibility

data source `accuknox_kubearmor_namespace_visibility` provides the current namespace based visibility information.

```
data "accuknox_kubearmor_namespace_visibility" "ns_vs" {
  name="kube-system"
}

output "ns_vs" {
  value = data.accuknox_kubearmor_namespace_visibility.ns_vs.visibility
}
```

### Read KubeArmror Stable Version

data source `accuknox_kubearmor_stable_version` provides the stable verion of kubearmor.

```
data "accuknox_kubearmor_stable_version" "stable_version" {}

output "stable_version" {
  value = data.accuknox_kubearmor_stable_version.stable_version.version
}
```

## Example - Discovery Engine Resources

### Managing Discovered Policy

Once you discovered the policies using `karmor discover`, using `accuknox_discovery_engine_discovered_policy` you can add the discovered secuity policy in the workload system.

> This would only add the discovered policy, and will not be applied on the workloads, since this will have status set to `Inactive`.

```
resource "accuknox_discovery_engine_discovered_policy" "dsp"{
  name="autopol-system-3960684242"
  namespace="wordpress-mysql"
  policy= <<-EOT
apiVersion: security.kubearmor.com/v1
kind: KubeArmorPolicy
metadata:
  name: autopol-system-3960684242
  namespace: wordpress-mysql
spec:
  action: Allow
  file:
    matchPaths:
    - fromSource:
      - path: /usr/sbin/apache2
      path: /dev/urandom
    - fromSource:
      - path: /usr/local/bin/php
      path: /etc/hosts
  network:
    matchProtocols:
    - fromSource:
      - path: /usr/local/bin/php
      protocol: tcp
    - fromSource:
      - path: /usr/local/bin/php
      protocol: udp
  process:
    matchPaths:
    - path: /usr/sbin/apache2
    - path: /usr/local/bin/php
  selector:
    matchLabels:
      app: wordpress
  severity: 1
  EOT
}
```

### Enabling Discovered Policy

To make a dicovered policy work we need to update the discovered Policy in `Active` state, therefore we use `accuknox_discovery_engine_enable_discovered_policy` to do so. This will update the policy's state

```
resource "accuknox_discovery_engine_enable_discovered_policy" "dsp_enable"{
  name="autopol-system-3960684242"
  namespace="wordpress-mysql"
}
```

### Managing Discovery Engine Configuration

We can manage the discovery engine configuration with resource `accuknox_kubearmor_configuration`.

```
resource "accuknox_discovery_engine_configuration" "example_config" {
  name      = "discovery-engine-config"
  namespace = "accuknox-agents"
  data = {
    "conf.yaml" = <<EOT
      application:
        name: discovery-engine
        network:
          operation-mode: 1
          cron-job-time-interval: "0h0m10s"
          operation-trigger: 5
          network-log-from: "kubearmor"
          network-log-file: "./flow.json"
          network-policy-to: "db"
          network-policy-dir: "./"
          namespace-filter:
          - "!kube-system"
        system:
          operation-mode: 1
          cron-job-time-interval: "0h0m10s"
          operation-trigger: 5
          system-log-from: "kubearmor"
          system-log-file: "./log.json"
          system-policy-to: "db"
          system-policy-dir: "./"
          deprecate-old-mode: true
          namespace-filter:
          - "!kube-system"
          fromsource-filter:
          - "knoxAutoPolicy"
        admission-controller:
          generic-policy-list:
          - "restrict-deprecated-registry"
          - "prevent-cr8escape"
          - "check-kernel-version"
          - "restrict-ingress-defaultbackend"
          - "restrict-nginx-ingress-annotations"
          - "restrict-ingress-paths"
          - "prevent-naked-pods"
          - "restrict-wildcard-verbs"
          - "restrict-wildcard-resources"
          - "require-requests-limits"
          - "require-pod-probes"
          - "drop-cap-net-raw"
        cluster:
          cluster-info-from: "k8sclient"
      observability:
        enable: true
        cron-job-time-interval: "0h0m10s"
        dbname: ./accuknox-obs.db
        system-observability: true
        network-observability: false
        write-logs-to-db: false
        summary-jobs:
          publisher: true
          write-summary-to-db: true
          cron-interval: "0h1m00s"
      database:
        driver: sqlite3
        host: mysql.explorer.svc.cluster.local
        port: 3306
        user: root
        password: password
        dbname: discovery-engine
        table-configuration: auto_policy_config
        table-network-log: network_log
        table-network-policy: network_policy
        table-system-log: system_log
        table-system-policy: system_policy
      feed-consumer:
        driver: "pulsar"
        servers:
          - "pulsar-proxy.accuknox-dev-pulsar.svc.cluster.local:6650"
        topic:
          cilium: "persistent://accuknox/datapipeline/ciliumalertsflowv1"
          kubearmor: "persistent://accuknox/datapipeline/kubearmoralertsflowv1"
        encryption:
          enable: false
          ca-cert: /kafka-ssl/ca.pem
        auth:
          enable: false
          cert: /kafka-ssl/user.cert.pem
          key: /kafka-ssl/user.key.pem
      logging:
        level: "INFO"
      cilium-hubble:
        url: hubble-relay.kube-system.svc.cluster.local
        port: 80
      kubearmor:
        url: kubearmor.kube-system.svc.cluster.local
        port: 32767
      recommend:
        operation-mode: 1
        cron-job-time-interval: "1h0m00s"
        recommend-host-policy: true
        template-version: ""
        admission-controller-policy: false
      license:
        enabled: false
        validate: "user-id"
      dsp:
        auto-deploy-dsp: false
    EOT
  }
}
```

## Example - Discovery Engine Data Sources

### Read Applied Discovered Policy

data source `accuknox_discovery_engine_discovered_policy` provides the discovered security policy based on `name` and `namespace`.

```
data "accuknox_discovery_engine_discovered_policy" "dsp"{
  name="autopol-system-3960684242"
  namespace="wordpress-mysql"
}

output "dsp" {
  value = data.accuknox_discovery_engine_discovered_policy.dsp.policy
}
```

### Read Discovery Engine Configuration

data source `accuknox_discovery_engine_configuration` provides the current discovery engine's configuration.

```
data "accuknox_discovery_engine_configuration" "data_cm" {
  name="discovery-engine-config"
  namespace="accuknox-agents"
}

output "data_cm" {
  value = data.accuknox_discovery_engine_configuration.data_cm.data
}
```
