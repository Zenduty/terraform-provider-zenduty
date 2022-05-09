package zenduty

import (
	"context"
	"errors"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSchedules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateSchedule,
		UpdateContext: resourceUpdateSchedule,
		DeleteContext: resourceDeleteSchedule,
		ReadContext:   resourceReadSchedule,
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
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"layers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"shift_length": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3600, 365*24*3600),
						},
						"rotation_start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rotation_end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"users": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// "restrictions": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Elem: &schema.Schema{
						// 		Type: schema.TypeString,
						// 	},
						// },
					},
				},
			},
			"overrides": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceCreateSchedule(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	new_schedule := &client.CreateSchedule{}
	var diags diag.Diagnostics
	layers := d.Get("layers").([]interface{})
	overrides := d.Get("overrides").([]interface{})
	if v, ok := d.GetOk("name"); ok {
		new_schedule.Name = v.(string)
	}
	if v, ok := d.GetOk("summary"); ok {
		new_schedule.Summary = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		new_schedule.Description = v.(string)
	}
	if v, ok := d.GetOk("time_zone"); ok {
		new_schedule.Time_zone = v.(string)
	}
	if v, ok := d.GetOk("team_id"); ok {
		new_schedule.Team = v.(string)
	}
	new_schedule.Layers = make([]client.CreateLayers, len(layers))
	new_schedule.Overrides = make([]client.Overrides, len(overrides))

	for i, layer := range layers {
		layer_map := layer.(map[string]interface{})
		new_layer := client.CreateLayers{}

		if v, ok := layer_map["name"]; ok {
			new_layer.Name = v.(string)
		}
		// if v, ok := layer_map["time_zone"]; ok {
		// 	new_layer.Time_zone = v.(string)
		// }
		if v, ok := layer_map["shift_length"]; ok {
			new_layer.ShiftLength = v.(int)
		}
		if v, ok := layer_map["rotation_start_time"]; ok {
			new_layer.RotationStartTime = v.(string)
		}
		if v, ok := layer_map["rotation_end_time"]; ok {
			new_layer.RotationEndTime = v.(string)
		}
		if v, ok := layer_map["users"]; ok {
			users := v.([]interface{})
			new_layer.Users = make([]client.CreateUserLayer, len(users))
			for j, user := range users {
				new_user := client.CreateUserLayer{}
				new_user.User = user.(string)
				new_layer.Users[j] = new_user
			}
		}
		new_schedule.Layers[i] = new_layer
	}
	schedule, err := apiclient.Schedules.CreateSchedule(new_schedule.Team, new_schedule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(schedule.Unique_Id)
	return diags

}

func resourceUpdateSchedule(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	new_schedule := &client.CreateSchedule{}
	team_id := d.Get("team_id").(string)
	layers := d.Get("layers").([]interface{})
	overrides := d.Get("overrides").([]interface{})
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	if v, ok := d.GetOk("name"); ok {
		new_schedule.Name = v.(string)
	}
	if v, ok := d.GetOk("summary"); ok {
		new_schedule.Summary = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		new_schedule.Description = v.(string)
	}
	if v, ok := d.GetOk("time_zone"); ok {
		new_schedule.Time_zone = v.(string)
	}
	new_schedule.Layers = make([]client.CreateLayers, len(layers))
	new_schedule.Overrides = make([]client.Overrides, len(overrides))

	for i, layer := range layers {
		layer_map := layer.(map[string]interface{})
		new_layer := client.CreateLayers{}

		if v, ok := layer_map["name"]; ok {
			new_layer.Name = v.(string)
		}
		// if v, ok := layer_map["time_zone"]; ok {
		// 	new_layer.Time_zone = v.(string)
		// }
		if v, ok := layer_map["shift_length"]; ok {
			new_layer.ShiftLength = v.(int)
		}
		if v, ok := layer_map["rotation_start_time"]; ok {
			new_layer.RotationStartTime = v.(string)
		}
		if v, ok := layer_map["rotation_end_time"]; ok {
			new_layer.RotationEndTime = v.(string)
		}
		if v, ok := layer_map["users"]; ok {
			users := v.([]interface{})
			new_layer.Users = make([]client.CreateUserLayer, len(users))
			for j, user := range users {
				new_user := client.CreateUserLayer{}
				new_user.User = user.(string)
				new_layer.Users[j] = new_user
			}
		}
		new_schedule.Layers[i] = new_layer
	}

	_, err := apiclient.Schedules.UpdateScheduleByID(team_id, id, new_schedule)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceDeleteSchedule(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	err := apiclient.Schedules.DeleteScheduleByID(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func resourceReadSchedule(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	service, err := apiclient.Schedules.GetScheduleByID(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", service.Name)
	d.Set("summary", service.Summary)
	d.Set("description", service.Description)
	return diags
}
