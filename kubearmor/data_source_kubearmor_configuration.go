package kubearmor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceKubearmorConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubearmorConfigurationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceKubearmorConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	om := metav1.ObjectMeta{
		Namespace: d.Get("namespace").(string),
		Name:      d.Get("name").(string),
	}
	d.SetId(BuildId(om))

	return resourceKubeArmorConfigurationRead(ctx, d, meta)
}
