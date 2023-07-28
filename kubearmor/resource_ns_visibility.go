package kubearmor

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func resourceNsVisibility() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceNsVisibiltyCreate,
		ReadContext:   resourceNsVisibiltyRead,
		UpdateContext: resourceNsVisibiltyUpdate,
		DeleteContext: resourceNsVisibiltyDelete,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"file": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"network": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"capabilities": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"process": {
				Type:     schema.TypeBool,
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

func resourceNsVisibiltyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("namespace").(string)
	d.SetId(name)

	diag := resourceNsVisibiltyUpdate(ctx, d, meta)
	if diag.HasError() {
		d.SetId("")
	}

	return resourceNsVisibiltyRead(ctx, d, meta)
}

func resourceNsVisibiltyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var annotation string
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

	annotation = namespace.Annotations["kubearmor-visibility"]

	if strings.Contains(annotation, "file") {
		d.Set("file", true)
	} else {
		d.Set("file", false)
	}

	if strings.Contains(annotation, "file") {
		d.Set("process", true)
	} else {
		d.Set("process", false)
	}

	if strings.Contains(annotation, "file") {
		d.Set("capabilities", true)
	} else {
		d.Set("capabilities", false)
	}

	if strings.Contains(annotation, "file") {
		d.Set("network", true)
	} else {
		d.Set("network", false)
	}

	return nil
}

func resourceNsVisibiltyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var visibility []string
	client, err := connectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("namespace").(string)

	_, err = client.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return diag.Errorf("The namespace %q does not exist", d.Id())
	}

	if d.Get("file").(bool) {
		visibility = append(visibility, "file")
	}

	if d.Get("process").(bool) {
		visibility = append(visibility, "process")
	}

	if d.Get("capabilities").(bool) {
		visibility = append(visibility, "capabilities")
	}

	if d.Get("network").(bool) {
		visibility = append(visibility, "network")
	}

	patchPayload := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"kubearmor-visibility": strings.Join(visibility, ","),
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

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceNsVisibiltyRead(ctx, d, meta)
}

func resourceNsVisibiltyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := connectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("namespace").(string)

	patchPayload := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"kubearmor-visibility": nil,
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
