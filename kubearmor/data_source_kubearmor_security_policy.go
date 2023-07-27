package kubearmor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceKubearmorSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubearmorSecurityPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// "tags": {
						// 	Type:     schema.TypeList,
						// 	Computed: true,
						// },
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"selector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_labels": {
										Type:     schema.TypeMap,
										Computed: true,
									},
								},
							},
						},
						"file": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_directories": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dir": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"owner_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"recursive": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"from_source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"match_paths": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"owner_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"from_source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"match_patterns": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"pattern": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"owner_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"severity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Computed: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"process": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_directories": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dir": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"owner_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"recursive": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"from_source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"match_paths": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"owner_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"from_source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"match_patterns": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"pattern": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"owner_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"severity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Computed: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"capabilities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_capabilities": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"capabilities": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"from_source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"severity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Computed: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"network": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_protocols": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "tags": {
												// 	Type:     schema.TypeList,
												// 	Computed: true,
												// },
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"from_source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"severity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									// "tags": {
									// 	Type:     schema.TypeList,
									// 	Computed: true,
									// },
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						// "syscalls": {
						// 	Type:     schema.TypeList,
						// 	Computed: true,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"match_syscalls": {
						// 				Type:     schema.TypeList,
						// 				Computed: true,
						// 				Elem: &schema.Resource{
						// 					Schema: map[string]*schema.Schema{
						// 						"syscall": {
						// 							Type:     schema.TypeList,
						// 							Computed: true,
						// 						},
						// 						"from_source": {
						// 							Type:     schema.TypeList,
						// 							Computed: true,
						// 							Elem: &schema.Resource{
						// 								Schema: map[string]*schema.Schema{
						// 									"path": {
						// 										Type:     schema.TypeString,
						// 										Computed: true,
						// 									},
						// 									"dir": {
						// 										Type:     schema.TypeString,
						// 										Computed: true,
						// 									},
						// 									"recursive": {
						// 										Type:     schema.TypeString,
						// 										Computed: true,
						// 									},
						// 								},
						// 							},
						// 						},
						// 					},
						// 				},
						// 			},
						// 			"match_paths": {
						// 				Type:     schema.TypeList,
						// 				Computed: true,
						// 				Elem: &schema.Resource{
						// 					Schema: map[string]*schema.Schema{
						// 						"path": {
						// 							Type:     schema.TypeString,
						// 							Computed: true,
						// 						},
						// 						"recursive": {
						// 							Type:     schema.TypeString,
						// 							Computed: true,
						// 						},
						// 						"syscall": {
						// 							Type:     schema.TypeList,
						// 							Computed: true,
						// 						},
						// 						"from_source": {
						// 							Type:     schema.TypeList,
						// 							Computed: true,
						// 							Elem: &schema.Resource{
						// 								Schema: map[string]*schema.Schema{
						// 									"path": {
						// 										Type:     schema.TypeString,
						// 										Computed: true,
						// 									},
						// 									"dir": {
						// 										Type:     schema.TypeString,
						// 										Computed: true,
						// 									},
						// 									"recursive": {
						// 										Type:     schema.TypeString,
						// 										Computed: true,
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
					},
				},
			},
		},
	}
}

func dataSourceKubearmorSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	KSPClient, err := connectKubearmorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)
	d.SetId(namespace + "/" + name)

	policy, err := KSPClient.SecurityV1().KubeArmorPolicies(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			d.SetId("")
			return diag.FromErr(err)

		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	if err := d.Set("policy", flattenPolicy(policy)); err != nil {
		return diag.FromErr(err)

	}

	return nil
}
