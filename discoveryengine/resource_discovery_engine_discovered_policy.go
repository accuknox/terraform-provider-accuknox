package discoveryengine

import (
	"context"
	"encoding/json"

	dspv1 "github.com/accuknox/auto-policy-discovery/pkg/discoveredpolicy/api/security.kubearmor.com/v1"
	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/accuknox/terraform-provider-accuknox/kubearmor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ResourceDiscoveryEngineDiscoveredPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiscoveryEngineDiscoveredPolicyCreate,
		ReadContext:   resourceDiscoveryEngineDiscoveredPolicyRead,
		UpdateContext: resourceDiscoveryEngineDiscoveredPolicyUpdate,
		DeleteContext: resourceDiscoveryEngineDiscoveredPolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDiscoveryEngineDiscoveredPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	dsp := &dspv1.DiscoveredPolicy{}
	dsp.Name = name
	dsp.Namespace = namespace
	dsp.Spec.PolicyStatus = "Inactive"

	pol := d.Get("policy")
	policy, err := json.Marshal(pol)
	if err != nil {
		return diag.FromErr(err)
	}

	polSpec := apiextv1.JSON{
		Raw: policy,
	}

	dsp.Spec.Policy = &polSpec

	_, err = client.DiscoveredPolicies(namespace).Create(context.TODO(), dsp, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDiscoveryEngineDiscoveredPolicyRead(ctx, d, meta)

}

func resourceDiscoveryEngineDiscoveredPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectDiscoveryEngineClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := kubearmor.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dsp, err := client.DiscoveredPolicies(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", dsp.Name)
	d.Set("namespace", dsp.Namespace)
	d.Set("policy", dsp.Spec.Policy)

	return nil
}

func resourceDiscoveryEngineDiscoveredPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//

	// d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceDiscoveryEngineDiscoveredPolicyRead(ctx, d, meta)
}

func resourceDiscoveryEngineDiscoveredPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectDiscoveryEngineClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := kubearmor.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DiscoveredPolicies(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
