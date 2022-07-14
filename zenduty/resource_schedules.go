package zenduty

import (
	"context"
	"errors"
	"regexp"
	"time"

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
						"restriction_type": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 2),
							Default:      0,
						},
						"restrictions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"duration": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 7*24*3600),
									},
									"start_day_of_week": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 7),
									},
									"start_time_of_day": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9]):([0-9]|[0-5][0-9])$`), "must be in the format HH:MM:SS"),
									},
								},
							},
						},
					},
				},
			},
			"overrides": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{1}$`), "must be a valid user id"),
						},
					},
				},
			},
		},
	}
}

func createSchedule(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.CreateSchedule, diag.Diagnostics) {
	new_schedule := &client.CreateSchedule{}

	layers := d.Get("layers").([]interface{})
	overrides := d.Get("overrides").([]interface{})
	if v, ok := d.GetOk("name"); ok {
		if v.(string) == "" {
			return nil, diag.FromErr(errors.New("name must not be empty"))
		}
		new_schedule.Name = v.(string)
	}
	if v, ok := d.GetOk("summary"); ok {

		new_schedule.Summary = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {

		new_schedule.Description = v.(string)
	}
	if v, ok := d.GetOk("time_zone"); ok {
		if v.(string) == "" {
			return nil, diag.FromErr(errors.New("time_zone must not be empty"))
		}
		new_schedule.Time_zone = v.(string)
	}
	if v, ok := d.GetOk("team_id"); ok {
		if v.(string) == "" {
			return nil, diag.FromErr(errors.New("team_id must not be empty"))
		}

		new_schedule.Team = v.(string)
	}
	new_schedule.Layers = make([]client.CreateLayers, len(layers))
	new_schedule.Overrides = make([]client.Overrides, len(overrides))

	for o, override := range overrides {
		override := override.(map[string]interface{})
		new_override := client.Overrides{}

		if v, ok := override["name"]; ok {

			new_override.Name = v.(string)
		}
		if v, ok := override["start_time"]; ok {

			new_override.StartTime = v.(string)
			parsed_time, err := time.Parse("2006-01-02 15:04", v.(string))
			if err != nil {
				return nil, diag.FromErr(errors.New(err.Error()))
			}
			loc, zone_err := time.LoadLocation(new_schedule.Time_zone)
			if zone_err != nil {
				return nil, diag.FromErr(errors.New(zone_err.Error()))
			}

			new_override.StartTime = parsed_time.In(loc).Format(time.RFC3339)

		}
		if v, ok := override["end_time"]; ok {
			new_override.EndTime = v.(string)
			parsed_end_time, err := time.Parse("2006-01-02 15:04", new_override.EndTime)

			if err != nil {
				return nil, diag.FromErr(errors.New(err.Error()))
			}

			loc, zone_err := time.LoadLocation(new_schedule.Time_zone)
			if zone_err != nil {
				return nil, diag.FromErr(errors.New(zone_err.Error()))
			}

			new_override.EndTime = parsed_end_time.In(loc).Format(time.RFC3339)

		}
		if v, ok := override["user"]; ok {

			new_override.User = v.(string)
		}
		new_schedule.Overrides[o] = new_override
	}

	for i, layer := range layers {
		layer_map := layer.(map[string]interface{})
		new_layer := client.CreateLayers{}

		if v, ok := layer_map["name"]; ok {
			if v.(string) == "" {
				return nil, diag.FromErr(errors.New("name must not be empty"))
			}

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

		if v, ok := layer_map["restriction_type"]; ok {
			new_layer.RestrictionType = v.(int)
		}
		if v, ok := layer_map["restrictions"]; ok {
			restrictions := v.([]interface{})
			new_layer.Restrictions = make([]client.Restrictions, len(restrictions))
			for j, restriction := range restrictions {
				if new_layer.RestrictionType == 0 {
					return nil, diag.FromErr(errors.New("restrictions must be set to add restrictions.. ie daily(1) or weekly(2)"))
				}
				restriction_map := restriction.(map[string]interface{})
				new_restriction := client.Restrictions{}
				if v, ok := restriction_map["duration"]; ok {
					new_restriction.Duration = v.(int)
					if new_layer.RestrictionType == 1 && new_restriction.Duration >= 86400 {
						return nil, diag.FromErr(errors.New("duration must be less than 86400 for daily restriction ie 24 hours"))
					} else if new_layer.RestrictionType == 2 && new_restriction.Duration >= 604800 {
						return nil, diag.FromErr(errors.New("duration must be less than 604800 for weekly restriction ie 7 days"))
					}
				}
				if v, ok := restriction_map["start_day_of_week"]; ok {

					if new_layer.RestrictionType == 1 {
						new_restriction.StartDayOfWeek = 7
					} else {
						new_restriction.StartDayOfWeek = v.(int)
					}

				}
				if v, ok := restriction_map["start_time_of_day"]; ok {
					new_restriction.StartTimeOfDay = v.(string)
				}

				new_layer.Restrictions[j] = new_restriction
			}
		}

		new_schedule.Layers[i] = new_layer
	}
	return new_schedule, nil
}

func resourceCreateSchedule(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	var diags diag.Diagnostics
	new_schedule, create_error := createSchedule(Ctx, d, m)
	if create_error != nil {
		return create_error
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
	team_id := d.Get("team_id").(string)
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	new_schedule, create_error := createSchedule(Ctx, d, m)
	if create_error != nil {
		return create_error
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
