---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "accuknox_kubearmor_namespace_visibility Data Source - terraform-provider-accuknox"
subcategory: "kubearmor"
description: |-
  Provides the namespace posture based on name.
---

# accuknox_kubearmor_namespace_visibility (Data Source)

This data source provides a mechanism to view the configuration of the namespace visibility based on namespace name.

## Example Usage

```
data "accuknox_kubearmor_namespace_visibility" "ns_vs" {
  name="kube-system"
}

output "ns_vs" {
  value = data.accuknox_kubearmor_namespace_visibility.ns_vs.visibility
}
```

### Argument Reference

- `name` (Required) Namespace name.
