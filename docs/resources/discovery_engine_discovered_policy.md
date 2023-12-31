---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "accuknox_discovery_engine_discovered_policy Resource - terraform-provider-accuknox"
subcategory: "Discovery_Engine"
description: |-
  Add discovered policy.
---

# accuknox_discovery_engine_discovered_policy (Resource)

Once you discovered the policies using `karmor discover`, using `accuknox_discovery_engine_discovered_policy` you can add the discovered policy in the workload system.

> This would only add the discovered policy, and will not be applied on the workloads, since this will have status set to `Inactive`.

## Example Usage

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


### Argument Reference

- `name` (Required) Name of the policy.
- `namespace` (Optional) Namespace for the policy. If not mentioned then it will be assumed as `default`.
- `policy` (Required) Add the discovered policy.

