package zenduty

import (
	"context"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEsp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEspsRead,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"esp_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"escalation_policies": {
				Type:        schema.TypeList,
				Description: "List of Escalation policies",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"team": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rules": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delay": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"position": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unique_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"targets": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_type": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"target_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"position": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"repeat_policy": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"move_to_next": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"global_ep": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"assignee_strategy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "1 notifies every target, 2 assigns incidents via round-robin.",
						},
						"notify_round_robin_assignee_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether only the round-robin assignee is notified or called.",
						},
					},
				},
			},
		},
	}
}

func flattenEspItem(esp *client.EscalationPolicy) map[string]interface{} {
	strategy, notifyRRAssigneeOnly := flattenEspAssignmentSettings(esp.AssignmentSettings)
	return map[string]interface{}{
		"unique_id":                        esp.UniqueID,
		"name":                             esp.Name,
		"summary":                          esp.Summary,
		"description":                      esp.Description,
		"team":                             esp.Team,
		"rules":                            flattenRules(esp.Rules),
		"repeat_policy":                    esp.RepeatPolicy,
		"move_to_next":                     esp.MoveToNext,
		"global_ep":                        esp.GlobalEp,
		"assignee_strategy":                strategy,
		"notify_round_robin_assignee_only": notifyRRAssigneeOnly,
	}
}

func dataSourceEspsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	teamID := d.Get("team_id").(string)
	espID := d.Get("esp_id").(string)

	var diags diag.Diagnostics
	if espID != "" {
		esp, err := apiclient.Esp.GetEscalationPolicyByID(teamID, espID)
		if err != nil {
			return diag.FromErr(err)
		}
		items := []map[string]interface{}{flattenEspItem(esp)}

		if err := d.Set("escalation_policies", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", teamID, espID))

		return diags
	} else {

		esps, err := apiclient.Esp.GetEscalationPolicy(teamID)
		if err != nil {
			return diag.FromErr(err)
		}
		items := make([]map[string]interface{}, len(esps))
		for i := range esps {
			items[i] = flattenEspItem(&esps[i])
		}
		if err := d.Set("escalation_policies", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", teamID, espID))

		return diags
	}
}
