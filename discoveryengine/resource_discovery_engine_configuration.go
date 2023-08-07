package discoveryengine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/accuknox/terraform-provider-accuknox/clienthandler"
	"github.com/accuknox/terraform-provider-accuknox/kubearmor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func ResourceDiscoveryEngineConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiscoveryEngineConfigCreate,
		ReadContext:   resourceDiscoveryEngineConfigRead,
		UpdateContext: resourceDiscoveryEngineConfigUpdate,
		DeleteContext: resourceDiscoveryEngineConfigDelete,
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

func resourceDiscoveryEngineConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	d.SetId(namespace + "/" + name)
	diag := resourceDiscoveryEngineConfigUpdate(ctx, d, meta)
	if diag.HasError() {
		d.SetId("")
	}

	return resourceDiscoveryEngineConfigRead(ctx, d, meta)
}

func resourceDiscoveryEngineConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := kubearmor.IdParts(d.Id())
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

func resourceDiscoveryEngineConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	data := d.Get("data").(map[string]interface{})
	patchPayload := map[string]interface{}{
		"data": data,
	}

	patch, err := json.Marshal(patchPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := kubearmor.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CoreV1().ConfigMaps(namespace).Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceDiscoveryEngineConfigRead(ctx, d, meta)
}

func resourceDiscoveryEngineConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clienthandler.ConnectK8sClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := kubearmor.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	patchPayload := map[string]interface{}{
		"data": map[string]interface{}{
			"conf.yaml": "application:\n  name: discovery-engine\n  network:\n    operation-mode: 1                         # 1: cronjob | 2: one-time-job\n    cron-job-time-interval: \"0h0m10s\"         # format: XhYmZs\n    operation-trigger: 5\n    network-log-from: \"kubearmor\"             # db|hubble|feed-consumer|kubearmor\n    network-log-file: \"./flow.json\"           # file path\n    network-policy-to: \"db\"                   # db, file\n    network-policy-dir: \"./\"\n    namespace-filter:\n    - \"!kube-system\"\n  system:\n    operation-mode: 1                         # 1: cronjob | 2: one-time-job\n    cron-job-time-interval: \"0h0m10s\"         # format: XhYmZs\n    operation-trigger: 5\n    system-log-from: \"kubearmor\"              # db|kubearmor|feed-consumer\n    system-log-file: \"./log.json\"             # file path\n    system-policy-to: \"db\"                    # db, file\n    system-policy-dir: \"./\"\n    deprecate-old-mode: true\n    namespace-filter:\n    - \"!kube-system\"\n    fromsource-filter:\n    - \"knoxAutoPolicy\"\n  admission-controller:\n    generic-policy-list:\n    - \"restrict-deprecated-registry\"\n    - \"prevent-cr8escape\"\n    - \"check-kernel-version\"\n    - \"restrict-ingress-defaultbackend\"\n    - \"restrict-nginx-ingress-annotations\"\n    - \"restrict-ingress-paths\"\n    - \"prevent-naked-pods\"\n    - \"restrict-wildcard-verbs\"\n    - \"restrict-wildcard-resources\"\n    - \"require-requests-limits\"\n    - \"require-pod-probes\"\n    - \"drop-cap-net-raw\"\n  cluster:\n    cluster-info-from: \"k8sclient\"            # k8sclient|accuknox\n\nobservability:\n  enable: true\n  cron-job-time-interval: \"0h0m10s\"         # format: XhYmZs\n  dbname: ./accuknox-obs.db\n  system-observability: true\n  network-observability: false\n  write-logs-to-db: false\n  summary-jobs:\n    publisher: true\n    write-summary-to-db: true\n    cron-interval: \"0h1m00s\"\n\ndatabase:\n  driver: sqlite3\n  host: mysql.explorer.svc.cluster.local\n  port: 3306\n  user: root\n  password: password\n  dbname: discovery-engine\n  table-configuration: auto_policy_config\n  table-network-log: network_log\n  table-network-policy: network_policy\n  table-system-log: system_log\n  table-system-policy: system_policy\n\nfeed-consumer:\n  driver: \"pulsar\"\n  servers:\n    - \"pulsar-proxy.accuknox-dev-pulsar.svc.cluster.local:6650\"\n  topic:\n    cilium: \"persistent://accuknox/datapipeline/ciliumalertsflowv1\"\n    kubearmor: \"persistent://accuknox/datapipeline/kubearmoralertsflowv1\"\n  encryption:\n    enable: false\n    ca-cert: /kafka-ssl/ca.pem\n  auth:\n    enable: false\n    cert: /kafka-ssl/user.cert.pem\n    key: /kafka-ssl/user.key.pem\n\nlogging:\n  level: \"INFO\"\n\n# kubectl -n kube-system port-forward service/hubble-relay --address 0.0.0.0  --address :: 4245:80\n\ncilium-hubble:\n  url: hubble-relay.kube-system.svc.cluster.local\n  port: 80\n\nkubearmor:\n  url: kubearmor.kube-system.svc.cluster.local\n  port: 32767\n\n# Recommended policies configuration\n\nrecommend:\n  operation-mode: 1                       # 1: cronjob | 2: one-time-job\n  cron-job-time-interval: \"1h0m00s\"       # format: XhYmZs\n  recommend-host-policy: true\n  template-version: \"\"\n  admission-controller-policy: false\n\n# license\n\nlicense:\n  enabled: false\n  validate: \"user-id\"\n\ndsp:\n  auto-deploy-dsp: false\n",
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
