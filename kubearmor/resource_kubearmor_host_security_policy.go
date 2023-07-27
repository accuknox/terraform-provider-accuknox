package kubearmor

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kcV1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func resourceKubearmorHostSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubearmorHostSecurityPolicyCreate,
		ReadContext:   resourceKubearmorHostSecurityPolicyRead,
		UpdateContext: resourceKubearmorHostSecurityPolicyUpdate,
		DeleteContext: resourceKubearmorHostSecurityPolicyDelete,
		Schema: map[string]*schema.Schema{
			"policy": {
				Type:     schema.TypeString,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// "tags": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// },
			"message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_selector": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_labels": {
							Type:     schema.TypeMap,
							Required: true,
						},
					},
				},
			},
			"file": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_directories": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dir": {
										Type:     schema.TypeString,
										Required: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"owner_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"recursive": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"from_source": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"match_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"owner_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"from_source": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"match_patterns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"owner_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"severity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// "tags": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// },
						"message": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"process": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_directories": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dir": {
										Type:     schema.TypeString,
										Required: true,
									},
									"owner_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"recursive": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"from_source": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"match_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"owner_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"from_source": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"match_patterns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
									},
									"owner_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"severity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// "tags": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// },
						"message": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_capabilities": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capabilities": {
										Type:     schema.TypeString,
										Required: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"from_source": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"severity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// "tags": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// },
						"message": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_protocols": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
									},
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Optional: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"from_source": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"severity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// "tags": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// },
						"message": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// "syscalls": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"match_syscalls": {
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"syscall": {
			// 							Type:     schema.TypeList,
			// 							Optional: true,
			// 						},
			// 						"from_source": {
			// 							Type:     schema.TypeList,
			// 							Optional: true,
			// 							Elem: &schema.Resource{
			// 								Schema: map[string]*schema.Schema{
			// 									"path": {
			// 										Type:     schema.TypeString,
			// 										Optional: true,
			// 									},
			// 									"dir": {
			// 										Type:     schema.TypeString,
			// 										Optional: true,
			// 									},
			// 									"recursive": {
			// 										Type:     schema.TypeString,
			// 										Optional: true,
			// 									},
			// 								},
			// 							},
			// 						},
			// 					},
			// 				},
			// 			},
			// 			"match_paths": {
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"path": {
			// 							Type:     schema.TypeString,
			// 							Required: true,
			// 						},
			// 						"recursive": {
			// 							Type:     schema.TypeString,
			// 							Optional: true,
			// 						},
			// 						"syscall": {
			// 							Type:     schema.TypeList,
			// 							Optional: true,
			// 						},
			// 						"from_source": {
			// 							Type:     schema.TypeList,
			// 							Optional: true,
			// 							Elem: &schema.Resource{
			// 								Schema: map[string]*schema.Schema{
			// 									"path": {
			// 										Type:     schema.TypeString,
			// 										Optional: true,
			// 									},
			// 									"dir": {
			// 										Type:     schema.TypeString,
			// 										Optional: true,
			// 									},
			// 									"recursive": {
			// 										Type:     schema.TypeString,
			// 										Optional: true,
			// 									},
			// 								},
			// 							},
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceKubearmorHostSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	KSPClient, err := connectKubearmorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	ksp := &kcV1.KubeArmorHostPolicy{}

	if policyYAML, ok := d.GetOk("policy"); ok {
		ksp.Spec.Capabilities = kcV1.HostCapabilitiesType{
			MatchCapabilities: append([]kcV1.MatchHostCapabilitiesType{}, ksp.Spec.Capabilities.MatchCapabilities...),
		}
		ksp.Spec.Network = kcV1.HostNetworkType{
			MatchProtocols: append([]kcV1.MatchHostNetworkProtocolType{}, ksp.Spec.Network.MatchProtocols...),
		}

		err = yaml.Unmarshal([]byte(policyYAML.(string)), ksp)
		if err != nil {
			return diag.FromErr(err)
		}

	} else if _, ok := d.GetOk("name"); ok {
		ksp.ObjectMeta = metav1.ObjectMeta{
			Name: d.Get("name").(string),
		}
		ksp.Spec = expandHostSpec(d)
	}

	d.SetId(ksp.Name)

	_, err = KSPClient.SecurityV1().KubeArmorHostPolicies().Create(context.Background(), ksp, metav1.CreateOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("Policy %s already exists ...", ksp.Name)
		}
		return diag.FromErr(err)
	}

	return resourceKubearmorHostSecurityPolicyRead(ctx, d, meta)
}

func resourceKubearmorHostSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	KSPClient, err := connectKubearmorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := KSPClient.SecurityV1().KubeArmorHostPolicies().Get(context.Background(), d.Id(), metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	if _, ok := d.GetOk("policy"); ok {
		d.Set("policy", policy)
	} else if _, ok := d.GetOk("name"); ok {
		d.Set("severity", policy.Spec.Severity)
		d.Set("action", policy.Spec.Action)
		d.Set("node_selector", flattenHostSelector(policy.Spec.NodeSelector))
		d.Set("file", flattenFile(policy.Spec.File))
		d.Set("process", flattenProcess(policy.Spec.Process))
		d.Set("capabilities", flattenHostCapabilities(policy.Spec.Capabilities))
		d.Set("network", flattenHostNetworks(policy.Spec.Network))
		// d.Set("syscalls", flattenSyscalls(policy.Spec.Syscalls))
	}
	return nil
}

func resourceKubearmorHostSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	KSPClient, err := connectKubearmorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	ksp, err := KSPClient.SecurityV1().KubeArmorHostPolicies().Get(context.Background(), d.Id(), metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	ksp.ObjectMeta.Name = d.Get("name").(string)
	ksp.Spec = expandHostSpec(d)

	_, err = KSPClient.SecurityV1().KubeArmorHostPolicies().Update(context.Background(), ksp, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceKubearmorHostSecurityPolicyRead(ctx, d, meta)
}

func resourceKubearmorHostSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	KSPClient, err := connectKubearmorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	err = KSPClient.SecurityV1().KubeArmorHostPolicies().Delete(context.Background(), d.Id(), metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
