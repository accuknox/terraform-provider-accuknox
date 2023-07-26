package kubearmor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func resourceNsPosture() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNsPostureCreate,
		ReadContext:   resourceNsPostureRead,
		UpdateContext: resourceNsPostureUpdate,
		DeleteContext: resourceNsPostureDelete,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capabilities": {
				Type:     schema.TypeString,
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

func resourceNsPostureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("namespace").(string)
	d.SetId(name)

	diag := resourceNsPostureUpdate(ctx, d, meta)
	if diag.HasError() {
		d.SetId("")
	}

	return resourceNsPostureRead(ctx, d, meta)
}

func resourceNsPostureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := connectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	ns := d.Id()

	namespace, err := client.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	if v, ok := namespace.Annotations["kubearmor-file-posture"]; ok {
		d.Set("file", v)
	}

	if v, ok := namespace.Annotations["kubearmor-capabilities-posture"]; ok {
		d.Set("capabilities", v)
	}

	if v, ok := namespace.Annotations["kubearmor-network-posture"]; ok {
		d.Set("network", v)
	}

	return nil
}

func resourceNsPostureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := connectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("namespace").(string)

	patchPayload := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{},
		},
	}

	if v, ok := d.GetOk("file"); ok {
		patchPayload["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["kubearmor-file-posture"] = v.(string)
	}

	if v, ok := d.GetOk("capabilities"); ok {
		patchPayload["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["kubearmor-capabilities-posture"] = v.(string)
	}

	if v, ok := d.GetOk("network"); ok {
		patchPayload["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["kubearmor-network-posture"] = v.(string)
	}

	patch, err := json.Marshal(patchPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CoreV1().Namespaces().Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceNsPostureRead(ctx, d, meta)
}

func resourceNsPostureDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := connectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("namespace").(string)

	patchPayload := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"kubearmor-file-posture":         nil,
				"kubearmor-capabilities-posture": nil,
				"kubearmor-network-posture":      nil,
			},
		},
	}

	patch, err := json.Marshal(patchPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CoreV1().Namespaces().Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
