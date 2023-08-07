package discoveryengine

import (
	"context"

	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/accuknox/terraform-provider-accuknox/kubearmor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceDiscoveryEngineConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDiscoveryEngineConfigurationRead,
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

func dataSourceDiscoveryEngineConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	objectMeta := metav1.ObjectMeta{
		Namespace: namespace,
		Name:      name,
	}

	d.SetId(kubearmor.BuildId(objectMeta))

	cm, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("data", cm.Data)

	return nil
}
