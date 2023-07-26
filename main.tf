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

resource "accuknox_kubearmor_security_policy" "ksp-vault-protect" {
  # policy= <<-EOT
  # apiVersion: security.kubearmor.com/v1
  # kind: KubeArmorPolicy
  # metadata:
  #   name: ksp-ubuntu-5-net-icmp-audit
  # spec:
  #   severity: 8
  #   selector:
  #     matchLabels:
  #       container: ubuntu-5
  #   network:
  #     matchProtocols:
  #     - protocol: icmp
  #   action:
  #     Audit  
  # EOT
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
        dir= "/vaultn/"
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


  # name="ksp-ubuntu-5-net-icmp-audit"
  # severity= 7
  # action= "Audit"
  
  # selector {
  #   match_labels= {
  #     "app.kubernetes.io/name": "vault",
  #     "component": "server"
  #   }
  # }
  

  # network {
  #   match_protocols{
  #       protocol= "icmp"
  #     }
  # }


}


# data "accuknox_kubearmor_security_policy" "sp" {
#   name="block-pkg-mgmt-tools-exec"
#   namespace="default"
# }

# output "sp" {
#   value = data.kubearmor_security_policy.sp.policy
# }


# resource "accuknox_kubearmor_namespace_visibility" "visib-ns1" {
#   namespace     = "wordpress-mysql"
#   file          = true
#   network       = true
#   capabilities  = true
#   process       = true
# }

# resource "accuknox_kubearmor_namespace_posture" "ns_pos" {
#   namespace     = "wordpress-mysql"
#   file          = "block"
#   network       = "block"
#   capabilities  = "audit"
# }

# resource "accuknox_kubearmor_host_security_policy" "ksp-vault-protect" {
  # policy= <<-EOT
  # apiVersion: security.kubearmor.com/v1
  # kind: KubeArmorHostPolicy
  # metadata:
  #   name: hsp-kubearmor-dev-proc-path-block
  # spec:
  #   nodeSelector:
  #     matchLabels:
  #       kubernetes.io/hostname: kubearmor-dev
  #   severity: 5
  #   process:
  #     matchPaths:
  #     - path: /usr/bin/diff
  #   action:
  #     Block
  # EOT

#   name="host-policy"
#   severity= 5
#   action= "Block"
  
#   node_selector {
#     match_labels= {
#       "kubernetes.io/hostname": "kubearmor-dev",
#     }
#   }
  
#   process {
#     match_paths{
#         path= "/usr/bin/diffd"
#       }
#   }
# }

# resource "accuknox_kubearmor_configuration" "conf" {
#   name="kubearmor-config"
#   namespace="kube-system"
#   data={
#     "defaultCapabilitiesPosture"="audit",
#   }
# }

# resource "accuknox_kubearmor_de_configuration" "conf" {
#   name="discovery-engine-config"
#   namespace="accuknox-agents"
#   data={
#     conf.yaml={
#       "new":"n"
#     }
#   }
# }

# data "accuknox_kubearmor_configuration" "data_cm" {
#   name="kubearmor-config"
#   namespace="kube-system"
# }

# output "data_cm" {
#   value = data.kubearmor_configuration.data_cm.data
# }

# data "accuknox_kubearmor_namespace_visibility" "ns_vs" {
#   name="wordpress-mysql"
# }

# output "ns_vs" {
#   value = data.kubearmor_namespace_visibility.ns_vs.visibility
# }

# data "accuknox_kubearmor_namespace_posture" "ns_ps" {
#   name="wordpress-mysql"
# }

# output "ns_ps" {
#   value = data.kubearmor_namespace_posture.ns_ps.annotation
# }

# data "accuknox_kubearmor_host_security_policy" "hp" {
#   name="hsp-kubearmor-dev-proc-path-block"
# }

# output "hp" {
#   value = data.kubearmor_host_security_policy.hp.policy
# }

# resource "accuknox_kubearmor_sec_policy" "nginx-protects" {
#   policy= <<-EOT
#   apiVersion: security.kubearmor.com/v1
#   kind: KubeArmorPolicy
#   metadata:
#     name: block-pkg-mgmt-tools
#     namespace: default
#   spec:
#     selector:
#       matchLabels:
#         app: nginx
#     process:
#       matchPaths:
#       - path: /usr/bin/apt
#       - path: /usr/bin/apt-get
#     action:
#       Block
#   EOT
# }


# data "accuknox_kubearmor_installed_version" "installed_version" {}

# output "installed_version" {
#   value = data.kubearmor_installed_version.installed_version.version
# }

# data "accuknox_kubearmor_stable_version" "stable_version" {}

# output "stable_version" {
#   value = data.kubearmor_stable_version.stable_version.version
# }

# data "accuknox_kubearmor_node" "k_node" {}

# output "k_node" {
#   value = data.kubearmor_node.k_node.node_data
# }

# data "accuknox_kubearmor_discovery_engine_configuration" "data_cm" {
#   name="discovery-engine-config"
#   namespace="accuknox-agents"
# }

# output "data_cm" {
#   value = data.kubearmor_discovery_engine_configuration.data_cm.data
# }












































# terraform {
#   required_providers {
#     kubernetes = {
#       source  = "hashicorp/kubernetes"
#       version = ">= 2.0.0"
#     }
#   }
# }
# provider "kubernetes" {
#   config_path = "/etc/rancher/k3s/k3s.yaml"
# }
# resource "kubernetes_ingress_v1" "example_ingress" {
#   metadata {
#     name = "example-ingress"
#   }

#   spec {
#     default_backend {
#       service {
#         name = "myapp-1"
#         port {
#           number = 8080
#         }
#       }
#     }

#     rule {
#       http {
#         path {
#           backend {
#             service {
#               name = "myapp-1"
#               port {
#                 number = 8080
#               }
#             }
#           }

#           path = "/app1/*"
#         }

#         path {
#           backend {
#             service {
#               name = "myapp-2"
#               port {
#                 number = 8080
#               }
#             }
#           }

#           path = "/app2/*"
#         }
#       }
#     }

#     tls {
#       secret_name = "tls-secret"
#     }
#   }
# }



# data "kubernetes_all_namespaces" "allns" {}

# output "all-ns" {
#   value = data.kubernetes_all_namespaces.allns.namespaces
# }

# output "ns-present" {
#   value = contains(data.kubernetes_all_namespaces.allns.namespaces, "kube-system")
# }