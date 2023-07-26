package kubearmor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceKubearmorNsPosture() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubearmorNsPostureRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceKubearmorNsPostureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	annotation := make(map[string]string)
	client, err := connectK8sClient()
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

	if value, ok := namespace.Annotations["kubearmor-file-posture"]; ok {
		annotation["kubearmor-file-posture"] = value
	}
	if value, ok := namespace.Annotations["kubearmor-capabilities-posture"]; ok {
		annotation["kubearmor-capabilities-posture"] = value
	}
	if value, ok := namespace.Annotations["kubearmor-network-posture"]; ok {
		annotation["kubearmor-network-posture"] = value
	}

	d.Set("annotation", annotation)

	return nil
}
