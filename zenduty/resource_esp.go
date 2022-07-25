package zenduty

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEsp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateEsp,
		UpdateContext: resourceUpdateEsp,
		DeleteContext: resourceDeleteEsp,
		ReadContext:   resourceReadEsp,
		Importer: &schema.ResourceImporter{
			State: resourceEscalationPolicyImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"summary": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"targets": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_type": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"repeat_policy": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"move_to_next": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func CreateEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.EscalationPolicy, diag.Diagnostics) {
	new_esp := &client.EscalationPolicy{}
	rules := d.Get("rules").([]interface{})
	if v, ok := d.GetOk("name"); ok {
		new_esp.Name = v.(string)
	}
	if v, ok := d.GetOk("summary"); ok {
		new_esp.Summary = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		new_esp.Description = v.(string)
		if emptyString(new_esp.Description) {
			return nil, diag.FromErr(errors.New("description is empty"))
		}
	}
	if v, ok := d.GetOk("team_id"); ok {
		new_esp.Team = v.(string)
	}
	if v, ok := d.GetOk("repeat_policy"); ok {
		new_esp.Repeat_Policy = v.(int)
	}
	if v, ok := d.GetOk("move_to_next"); ok {
		new_esp.Move_To_Next = v.(bool)
	}
	new_esp.Rules = make([]client.Rules, len(rules))
	for i, rule := range rules {
		rule_map := rule.(map[string]interface{})
		new_rule := client.Rules{}
		if v, ok := rule_map["delay"]; ok {
			new_rule.Delay = v.(int)
		}
		// if v, ok := rule_map["position"]; ok {
		// 	new_rule.Position = v.(int)
		// }
		// if v, ok := rule_map["unique_id"]; ok {
		// 	new_rule.Unique_Id = v.(string)
		// }
		if v, ok := rule_map["targets"]; ok {
			targets := v.([]interface{})
			new_rule.Targets = make([]client.Targets, len(targets))
			for j, target := range targets {
				target_map := target.(map[string]interface{})
				new_target := client.Targets{}
				if v, ok := target_map["target_type"]; ok {
					new_target.Target_type = v.(int)
				}
				if v, ok := target_map["target_id"]; ok {
					new_target.Target_id = v.(string)
				}
				new_rule.Targets[j] = new_target
			}
		}
		new_esp.Rules[i] = new_rule
	}
	return new_esp, nil

}

func flattenRules(rules []client.Rules) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {

		result[i] = map[string]interface{}{
			"delay": rule.Delay,

			"targets": flattenTargets(rule.Targets),
		}
	}
	return result
}

func flattenTargets(targets []client.Targets) []map[string]interface{} {
	result := make([]map[string]interface{}, len(targets))
	for i, target := range targets {
		result[i] = map[string]interface{}{
			"target_type": target.Target_type,
			"target_id":   target.Target_id,
		}
	}
	return result
}

func resourceCreateEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	var diags diag.Diagnostics
	new_esp, createErr := CreateEsp(Ctx, d, m)
	if createErr != nil {
		return createErr
	}

	esp, err := apiclient.Esp.CreateEscalationPolicy(new_esp.Team, new_esp)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(esp.Unique_Id)

	readErr := resourceReadEsp(Ctx, d, m)
	if readErr != nil {
		return readErr
	}

	return diags

}

func resourceUpdateEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	new_esp, createErr := CreateEsp(Ctx, d, m)
	if createErr != nil {
		return createErr
	}
	id := d.Id()
	var diags diag.Diagnostics

	_, err := apiclient.Esp.UpdateEscalationPolicy(new_esp.Team, id, new_esp)

	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceDeleteEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	log.Printf("team_idss: %s", team_id)
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	err := apiclient.Esp.DeleteEscalationPolicy(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func resourceReadEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	if v, ok := d.GetOk("team_id"); ok {
		team_id = v.(string)
	}

	log.Printf("team_id: %s", team_id)

	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required "))
	}
	var diags diag.Diagnostics
	esp, err := apiclient.Esp.GetEscalationPolicyById(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", esp.Name)
	d.Set("team_id", esp.Team)
	d.Set("summary", esp.Summary)
	d.Set("description", esp.Description)
	d.Set("repeat_policy", esp.Repeat_Policy)
	d.Set("move_to_next", esp.Move_To_Next)
	if err := d.Set("rules", flattenRules(esp.Rules)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceEscalationPolicyImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<esp_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid escalationpolicyid (%q)", parts[1])
	}
	d.Set("team_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
