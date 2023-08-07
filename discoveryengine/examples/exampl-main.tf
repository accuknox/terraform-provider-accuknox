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

// resources

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

resource "accuknox_discovery_engine_enable_discovered_policy" "dsp_enable"{
  name="autopol-system-3960684242"
  namespace="wordpress-mysql"
}

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

//data sources

data "accuknox_discovery_engine_discovered_policy" "dsp"{
  name="autopol-system-3960684242"
  namespace="wordpress-mysql"
}

output "dsp" {
  value = data.accuknox_discovery_engine_discovered_policy.dsp.policy
}

data "accuknox_discovery_engine_configuration" "data_cm" {
  name="discovery-engine-config"
  namespace="accuknox-agents"
}

output "data_cm" {
  value = data.accuknox_discovery_engine_configuration.data_cm.data
}