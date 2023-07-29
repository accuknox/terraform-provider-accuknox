package kubearmor

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceKubearmorInstalledVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubearmorInstalledVersionRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubearmorInstalledVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	svc, err := client.CoreV1().Services("").List(context.Background(), metav1.ListOptions{
		LabelSelector: "kubearmor-app=kubearmor-relay",
	})
	if err != nil {
		return diag.FromErr(err)
	}

	pods, err := client.CoreV1().Pods(svc.Items[0].Namespace).List(context.Background(), metav1.ListOptions{LabelSelector: "kubearmor-app=kubearmor"})
	if err != nil {
		return diag.FromErr(err)
	}

	image := ""
	if len(pods.Items) > 0 {
		image = pods.Items[0].Spec.Containers[0].Image
	}

	d.Set("version", image)
	idsum := sha256.New()

	id := fmt.Sprintf("%x", idsum.Sum(nil))
	d.SetId(id)

	return nil
}
