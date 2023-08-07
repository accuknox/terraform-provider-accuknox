package discoveryengine

import (
	"context"
	"encoding/json"

	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceDiscoveryEngineDiscoveredPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDiscoveryEngineDiscoveredPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDiscoveryEngineDiscoveredPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectDiscoveryEngineClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	if namespace == "" {
		namespace = "default"
	}
	d.SetId(namespace + "/" + name)

	dsp, err := client.DiscoveredPolicies(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	var policy string

	err = json.Unmarshal(dsp.Spec.Policy.Raw, &policy)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy", policy); err != nil {
		return diag.FromErr(err)
	}
	d.Set("policy_status", dsp.Spec.PolicyStatus)

	return nil
}
