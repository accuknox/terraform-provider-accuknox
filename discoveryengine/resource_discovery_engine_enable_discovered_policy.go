package discoveryengine

import (
	"context"
	"fmt"
	"time"

	dspv1 "github.com/accuknox/auto-policy-discovery/pkg/discoveredpolicy/api/security.kubearmor.com/v1"
	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/accuknox/terraform-provider-accuknox/kubearmor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func ResourceDiscoveryEngineEnableDiscoveredPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiscoveryEngineEnableDiscoveredPolicyCreate,
		ReadContext:   resourceDiscoveryEngineEnableDiscoveredPolicyRead,
		UpdateContext: resourceDiscoveryEngineEnableDiscoveredPolicyUpdate,
		DeleteContext: resourceDiscoveryEngineEnableDiscoveredPolicyDelete,
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
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDiscoveryEngineEnableDiscoveredPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	if namespace == "" {
		namespace = "default"
	}
	d.SetId(namespace + "/" + name)
	diag := resourceDiscoveryEngineEnableDiscoveredPolicyUpdate(ctx, d, meta)
	if diag.HasError() {
		d.SetId("")
	}

	return resourceDiscoveryEngineEnableDiscoveredPolicyRead(ctx, d, meta)

}

func resourceDiscoveryEngineEnableDiscoveredPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return nil
}

func resourceDiscoveryEngineEnableDiscoveredPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectDiscoveryEngineClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := kubearmor.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dsp := &dspv1.DiscoveredPolicy{}
	dsp.Name = name
	dsp.Namespace = namespace
	// dsp.Spec.PolicyStatus = "Inactive"

	patchData := fmt.Sprintf(`{"spec":{"status":"%s"}}`, "Active")
	patchByte := []byte(patchData)

	_, err = client.DiscoveredPolicies(namespace).Patch(context.TODO(), name, types.MergePatchType, patchByte, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceDiscoveryEngineEnableDiscoveredPolicyRead(ctx, d, meta)
}

func resourceDiscoveryEngineEnableDiscoveredPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectDiscoveryEngineClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := kubearmor.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	patchData := fmt.Sprintf(`{"spec":{"status":"%s"}}`, "Inactive")
	patchByte := []byte(patchData)

	_, err = client.DiscoveredPolicies(namespace).Patch(context.TODO(), name, types.MergePatchType, patchByte, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
