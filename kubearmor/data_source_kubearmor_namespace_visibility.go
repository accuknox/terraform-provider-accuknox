package kubearmor

import (
	"context"

	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceKubearmorNsVisibility() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubearmorNsVisibilityRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubearmorNsVisibilityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("name").(string))

	namespace, err := client.CoreV1().Namespaces().Get(context.Background(), d.Id(), metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	d.Set("visibility", namespace.Annotations["kubearmor-visibility"])

	return nil
}
