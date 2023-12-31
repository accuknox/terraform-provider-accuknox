---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "accuknox_kubearmor_configuration Resource - terraform-provider-accuknox"
subcategory: "kubearmor"
description: |-
  This resource will give us the ability to manage the kubearmor configuration.

---

# accuknox_kubearmor_configuration (Resource)

`accuknox_kubearmor_configuration` will give us the ability to manage the kubearmor configuration. KubeArmor configuration handles global default posture, default visibility, cluster and gRPC port. Read more about [global default posture](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/default_posture.md#configuring-default-posture) and [global default visibility](https://github.com/kubearmor/KubeArmor/blob/main/getting-started/kubearmor_visibility.md)

**Default Configuration**

```
data:
  cluster: default
  defaultCapabilitiesPosture: audit
  defaultFilePosture: audit
  defaultNetworkPosture: audit
  gRPC: "32767"
  visibility: process,file,network
```

## Example Usage

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

### Argument Reference

- `name` (Required) Provide the name for the config map which is by default `kubearmor-config`.
- `namespace` (Required) Provide the namespace for config map in which configuration is set.
- `data` (Required) Provide the configured data to patch in config map. It's a map type field.
