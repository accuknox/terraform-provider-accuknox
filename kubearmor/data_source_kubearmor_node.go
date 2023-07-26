package kubearmor

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kubearmor/KubeArmor/KubeArmor/types"
	"github.com/kubearmor/kubearmor-client/k8s"
	"github.com/kubearmor/kubearmor-client/probe"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceKubearmorNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubearmorNodeRead,
		Schema: map[string]*schema.Schema{
			"node_data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kernel_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kubelet_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_lsm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_security": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"container_security": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"container_default_posture": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"capabilities": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"host_default_posture": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"capabilities": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"host_visibility": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKubearmorNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := k8s.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	svc, err := client.K8sClientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{
		LabelSelector: "kubearmor-app=kubearmor-relay",
	})
	if err != nil {
		return diag.FromErr(err)
	}

	o := probe.Options{
		Namespace: svc.Items[0].Namespace,
	}
	_, nodeData, err := probe.ProbeRunningKubeArmorNodes(client, o)
	if err != nil {
		log.Printf("[ERROR] nodeData err %#v", err)
		return diag.FromErr(err)
	}

	node := make([]interface{}, 1)

	for _, v := range nodeData {
		n := map[string]interface{}{
			"os_image":                  v.OSImage,
			"kernel_version":            v.KernelVersion,
			"kubelet_version":           v.KubeletVersion,
			"container_runtime":         v.ContainerRuntime,
			"active_lsm":                v.ActiveLSM,
			"host_security":             v.HostSecurity,
			"container_security":        v.ContainerSecurity,
			"container_default_posture": flattenPosture(v.ContainerDefaultPosture),
			"host_default_posture":      flattenPosture(v.HostDefaultPosture),
			"host_visibility":           v.HostVisibility,
		}
		node[0] = n
	}
	if err := d.Set("node_data", node); err != nil {
		return diag.FromErr(err)
	}

	idsum := sha256.New()

	id := fmt.Sprintf("%x", idsum.Sum(nil))
	d.SetId(id)

	return nil
}

func flattenPosture(posture types.DefaultPosture) []interface{} {
	pos := make(map[string]interface{})
	pos["file"] = posture.FileAction
	pos["capabilities"] = posture.FileAction
	pos["network"] = posture.FileAction

	return []interface{}{pos}

}
