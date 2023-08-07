package kubearmor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func ResourceKubearmorConfiguration() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceKubeArmorConfigurationCreate,
		ReadContext:   resourceKubeArmorConfigurationRead,
		UpdateContext: resourceKubeArmorConfigurationUpdate,
		DeleteContext: resourceKubeArmorConfigurationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"data": {
				Type:     schema.TypeMap,
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

func resourceKubeArmorConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	d.SetId(namespace + "/" + name)
	diag := resourceKubeArmorConfigurationUpdate(ctx, d, meta)
	if diag.HasError() {
		d.SetId("")
	}

	return resourceKubeArmorConfigurationRead(ctx, d, meta)
}

func resourceKubeArmorConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cm, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("data", cm.Data)

	return nil
}

func resourceKubeArmorConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	patchPayload := map[string]interface{}{
		"data": d.Get("data").(map[string]interface{}),
	}

	patch, err := json.Marshal(patchPayload)
	if err != nil {
		return diag.FromErr(err)
	}
	namespace, name, err := IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.CoreV1().ConfigMaps(namespace).Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceKubeArmorConfigurationRead(ctx, d, meta)
}

func resourceKubeArmorConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	patchPayload := map[string]interface{}{
		"data": map[string]string{
			"cluster":                    "default",
			"defaultCapabilitiesPosture": "audit",
			"defaultFilePosture":         "audit",
			"defaultNetworkPosture":      "audit",
			"gRPC":                       "32767",
			"visibility":                 "process,file,network",
		},
	}

	patch, err := json.Marshal(patchPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CoreV1().ConfigMaps(namespace).Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
