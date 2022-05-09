package zenduty

import (
	"context"
	"errors"

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
						"position": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
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

func CreateEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) *client.EscalationPolicy {
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
		if v, ok := rule_map["position"]; ok {
			new_rule.Position = v.(int)
		}
		if v, ok := rule_map["unique_id"]; ok {
			new_rule.Unique_Id = v.(string)
		}
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
	return new_esp

}

func resourceCreateEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	var diags diag.Diagnostics
	new_esp := CreateEsp(Ctx, d, m)

	esp, err := apiclient.Esp.CreateEscalationPolicy(new_esp.Team, new_esp)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(esp.Unique_Id)
	return diags

}

func resourceUpdateEsp(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	new_esp := CreateEsp(Ctx, d, m)
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
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	esp, err := apiclient.Esp.GetEscalationPolicyById(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", esp.Name)
	return diags
}
