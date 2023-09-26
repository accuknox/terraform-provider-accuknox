terraform {
  required_providers {
    accuknox = {
      source  = "terraform.example.com/local/accuknox"
      version = "1.0.0"
    }
  }
}

provider "accuknox" {
}

// resources

resource "accuknox_kubearmor_security_policy" "ksp-vault-protect" {
  name      = "ksp-vault-protect"
  namespace = "default"
  severity  = 7
  action    = "Allow"

  selector {
    match_labels = {
      "app.kubernetes.io/name" : "vault",
      "component" : "server"
    }
  }

  file {
    match_directories {
      dir       = "/vault/"
      recursive = true
      action    = "Block"
    }
    match_directories {
      dir       = "/"
      recursive = true
    }
    match_directories {
      dir       = "/vault/"
      recursive = true
      from_source {
        path = "/bin/vault"
      }
    }
  }

  process {
    match_paths {
      path = "/bin/busybox"
    }
    match_paths {
      path = "/bin/vault"
    }
  }
}

resource "accuknox_kubearmor_security_policy" "block-pkg-mgmt-tools-exec" {
  policy = <<-EOT
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


resource "accuknox_kubearmor_host_security_policy" "hsp-kubearmor-dev-proc-path-block" {
  policy = <<-EOT
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

resource "accuknox_kubearmor_host_security_policy" "hsp-kubearmor-dev-file-path-audit" {
  name     = "hsp-kubearmor-dev-file-path-audit"
  severity = 5
  action   = "Audit"

  node_selector {
    match_labels = {
      "kubernetes.io/hostname" : "kubearmor-dev",
    }
  }

  file {
    match_paths {
      path = "/etc/passwd"
    }
  }
}


resource "accuknox_kubearmor_namespace_visibility" "visib-ns1" {
  namespace    = "wordpress-mysql"
  file         = true
  network      = true
  capabilities = true
  process      = true
}

resource "accuknox_kubearmor_namespace_posture" "ns_pos" {
  namespace    = "wordpress-mysql"
  file         = "block"
  network      = "block"
  capabilities = "audit"
}

resource "accuknox_kubearmor_configuration" "conf" {
  name      = "kubearmor-config"
  namespace = "kube-system"
  data = {
    "defaultCapabilitiesPosture" = "audit",
  }
}

// data-sources

data "accuknox_kubearmor_configuration" "data_cm" {
  name      = "kubearmor-config"
  namespace = "kube-system"
}

output "data_cm" {
  value = data.accuknox_kubearmor_configuration.data_cm.data
}

data "accuknox_kubearmor_host_security_policy" "host-policy" {
  name = "hsp-kubearmor-dev-proc-path-block"
}

output "host-policy" {
  value = data.accuknox_kubearmor_host_security_policy.host-policy.policy
}

data "accuknox_kubearmor_installed_version" "installed_version" {}

output "installed_version" {
  value = data.accuknox_kubearmor_installed_version.installed_version.version
}

data "accuknox_kubearmor_node" "k_node" {}

output "k_node" {
  value = data.accuknox_kubearmor_node.k_node.node_data
}

data "accuknox_kubearmor_namespace_posture" "ns_ps" {
  name = "kube-system"
}

output "ns_ps" {
  value = data.accuknox_kubearmor_namespace_posture.ns_ps.annotation
}

data "accuknox_kubearmor_namespace_visibility" "ns_vs" {
  name = "kube-system"
}

output "ns_vs" {
  value = data.accuknox_kubearmor_namespace_visibility.ns_vs.visibility
}


data "accuknox_kubearmor_security_policy" "pkg-mgmt" {
  name      = "block-pkg-mgmt-tools-exec"
  namespace = "default"
}

output "pkg-mgmt" {
  value = data.accuknox_kubearmor_security_policy.pkg-mgmt.policy
}

data "accuknox_kubearmor_stable_version" "stable_version" {}

output "stable_version" {
  value = data.accuknox_kubearmor_stable_version.stable_version.version
}
