package zenduty

import (
	"context"
	"errors"
	"strconv"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlertRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAlertRules,
		UpdateContext: resourceUpdateAlertRules,
		DeleteContext: resourceDeleteAlertRules,
		ReadContext:   resourceReadAlertRules,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"actions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_type": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
func AlertRuleAction(Ctx context.Context, d *schema.ResourceData, m interface{}, newAlertRule *client.AlertRule) ([]client.AlertAction, diag.Diagnostics) {
	actions := d.Get("actions").([]interface{})
	newAlertRule.Actions = make([]client.AlertAction, len(actions))
	for i, action := range actions {
		rule_map := action.(map[string]interface{})
		new_action := client.AlertAction{}
		var value, key string

		if v, ok := rule_map["action_type"]; ok {
			new_action.ActionType = v.(int)
		}

		if (new_action.ActionType > 15) || (new_action.ActionType < 1) {
			return nil, diag.FromErr(errors.New("action_type is not valid"))
		}

		if v, ok := rule_map["value"]; ok {
			value = v.(string)
		}

		if (new_action.ActionType == 1) || (new_action.ActionType == 2) || (new_action.ActionType == 3) || (new_action.ActionType == 4) || (new_action.ActionType == 5) || (new_action.ActionType == 6) || (new_action.ActionType == 7) || (new_action.ActionType == 8) || (new_action.ActionType == 9) || (new_action.ActionType == 10) || (new_action.ActionType == 11) || (new_action.ActionType == 12) || (new_action.ActionType == 13) || (new_action.ActionType == 14) || (new_action.ActionType == 15) {
			if (new_action.ActionType != 3) && (value == "") {
				return nil, diag.FromErr(errors.New("value is required"))
			}
			if ((new_action.ActionType == 4) || (new_action.ActionType == 12) || (new_action.ActionType == 13) || (new_action.ActionType == 14) || (new_action.ActionType == 15)) && (!IsValidUUID(value)) {
				return nil, diag.FromErr(errors.New(value + " is not a valid UUID"))
			}

			if (new_action.ActionType == 7) && (!(value == "0" || value == "1")) {
				return nil, diag.FromErr(errors.New("incident urgency should be 0 or 1"))
			}
			if new_action.ActionType == 1 {
				i, err := strconv.Atoi(value)
				if i < 0 || i > 5 {
					return nil, diag.FromErr(errors.New("value should be between 0 and 5"))
				}
				if err != nil {
					return nil, diag.FromErr(errors.New("value is not valid"))
				}
				new_action.Value = value
			} else if new_action.ActionType == 3 {
				value = ""
			} else if new_action.ActionType == 4 {
				new_action.EscalationPolicy = value
				value = ""
			} else if new_action.ActionType == 6 {

				new_action.Assigned_To = value
				value = ""
			} else if new_action.ActionType == 11 {
				if v, ok := rule_map["key"]; ok {
					key = v.(string)
				}
				if key == "" {
					return nil, diag.FromErr(errors.New("key(ie..role_id) is required"))
				}
				if !IsValidUUID(key) {
					return nil, diag.FromErr(errors.New("key(ie..role_id) is not valid UUID"))
				}

				new_action.Key = key

			} else if new_action.ActionType == 14 {
				new_action.SLA = value
				value = ""
				if new_action.SLA == "" {
					return nil, diag.FromErr(errors.New("sla is required"))
				} else if !IsValidUUID(new_action.SLA) {
					return nil, diag.FromErr(errors.New("sla is not valid UUID"))
				}

			} else if new_action.ActionType == 15 {
				new_action.TeamPriority = value
				value = ""
				if new_action.TeamPriority == "" {
					return nil, diag.FromErr(errors.New("team_priority is required"))
				} else if !IsValidUUID(new_action.TeamPriority) {
					return nil, diag.FromErr(errors.New("team_priority is not valid UUID"))
				}
			}

			new_action.Value = value
		}

		newAlertRule.Actions[i] = new_action

	}
	return newAlertRule.Actions, nil

}

func ValidateAncCreateAlertRules(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.AlertRule, diag.Diagnostics) {
	newAlertRule := &client.AlertRule{}

	if v, ok := d.GetOk("description"); ok {
		newAlertRule.Description = v.(string)

	}
	if v, ok := d.GetOk("rule_json"); ok {
		newAlertRule.RuleJson = v.(string)

	}
	if newAlertRule.Description == "" {
		return nil, diag.FromErr(errors.New("description is required"))
	}
	if newAlertRule.RuleJson == "" {
		return nil, diag.FromErr(errors.New("rule_json is required"))
	}
	if !isJSONString(newAlertRule.RuleJson) {
		return nil, diag.FromErr(errors.New("rule_json is not valid JSON"))
	}
	actions, action_err := AlertRuleAction(Ctx, d, m, newAlertRule)
	if action_err != nil {
		return nil, action_err
	}

	newAlertRule.Actions = actions

	return newAlertRule, nil
}

func resourceCreateAlertRules(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	var diags diag.Diagnostics
	var team_id, service_id, integration_id string

	team_id = d.Get("team_id").(string)
	service_id = d.Get("service_id").(string)
	integration_id = d.Get("integration_id").(string)

	new_rule, rule_err := ValidateAncCreateAlertRules(Ctx, d, m)
	if rule_err != nil {
		return rule_err
	}

	alert_rule, err := apiclient.AlertRules.CreateAlertRule(team_id, service_id, integration_id, new_rule)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(alert_rule.Unique_Id)
	return diags
}

func resourceUpdateAlertRules(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	var diags diag.Diagnostics
	var team_id, service_id, integration_id string

	team_id = d.Get("team_id").(string)
	service_id = d.Get("service_id").(string)
	integration_id = d.Get("integration_id").(string)

	new_rule, rule_err := ValidateAncCreateAlertRules(Ctx, d, m)
	if rule_err != nil {
		return rule_err
	}

	_, err := apiclient.AlertRules.UpdateAlertRule(team_id, service_id, integration_id, d.Id(), new_rule)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceReadAlertRules(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	var diags diag.Diagnostics
	var team_id, service_id, integration_id string

	team_id = d.Get("team_id").(string)
	service_id = d.Get("service_id").(string)
	integration_id = d.Get("integration_id").(string)

	rule, err := apiclient.AlertRules.GetAlertRule(team_id, service_id, integration_id, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(rule.Unique_Id)
	return diags
}

func resourceDeleteAlertRules(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiclient, _ := m.(*Config).Client()
	team_id, service_id, integration_id, unique_id := d.Get("team_id").(string), d.Get("service_id").(string), d.Get("integration_id").(string), d.Id()

	err := apiclient.AlertRules.DeleteAlertRule(team_id, service_id, integration_id, unique_id)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
