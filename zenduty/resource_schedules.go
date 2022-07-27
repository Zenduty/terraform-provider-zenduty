package zenduty

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
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
		Importer: &schema.ResourceImporter{
			State: resourceScheduleImporter,
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

func buildScheduleLayerRescrition(new_layer *client.CreateLayers, layer_map map[string]interface{}, d *schema.ResourceData) ([]client.Restrictions, diag.Diagnostics) {
	if v, ok := layer_map["restriction_type"]; ok {
		new_layer.RestrictionType = v.(int)
	}

	if v, ok := layer_map["restrictions"]; ok {
		restrictions := v.([]interface{})
		Restrictions := make([]client.Restrictions, len(restrictions))
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

			Restrictions[j] = new_restriction
		}
		return Restrictions, nil
	}
	return nil, nil
}

func buildScheduleLayer(TimeZone string, ctx context.Context, d *schema.ResourceData) ([]client.CreateLayers, diag.Diagnostics) {
	layers := d.Get("layers").([]interface{})
	Layers := make([]client.CreateLayers, len(layers))

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
			loc, _ := time.LoadLocation(TimeZone)

			parsed_time, parse_err := time.ParseInLocation("2006-01-02 15:04", v.(string), loc)
			if parse_err == nil {
				new_layer.RotationStartTime = parsed_time.In(time.UTC).Format(time.RFC3339)
			}

		}
		if v, ok := layer_map["rotation_end_time"]; ok {
			new_layer.RotationEndTime = v.(string)
			loc, _ := time.LoadLocation(TimeZone)
			parsed_time, parsed_err := time.ParseInLocation("2006-01-02 15:04", v.(string), loc)
			if parsed_err == nil {
				new_layer.RotationEndTime = parsed_time.In(time.UTC).Format(time.RFC3339)
			}

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
		new_restriction, restrictionErr := buildScheduleLayerRescrition(&new_layer, layer_map, d)
		if restrictionErr != nil {
			return nil, restrictionErr
		}
		if new_restriction != nil {
			new_layer.Restrictions = new_restriction
		}
		Layers[i] = new_layer
	}
	return Layers, nil

}

func buildScheduleOverride(new_schedule *client.CreateSchedule, d *schema.ResourceData) ([]client.Overrides, diag.Diagnostics) {
	overrides := d.Get("overrides").([]interface{})
	Overrides := make([]client.Overrides, len(overrides))

	for o, override := range overrides {
		override := override.(map[string]interface{})
		new_override := client.Overrides{}

		if v, ok := override["name"]; ok {

			new_override.Name = v.(string)
		}
		if v, ok := override["start_time"]; ok {

			new_override.StartTime = v.(string)

			loc, zone_err := time.LoadLocation(new_schedule.Time_zone)
			if zone_err != nil {
				return nil, diag.FromErr(errors.New(zone_err.Error()))
			}
			parsed_time, err := time.ParseInLocation("2006-01-02 15:04", v.(string), loc)
			if err != nil {
				return nil, diag.FromErr(errors.New(err.Error()))
			}

			new_override.StartTime = parsed_time.In(time.UTC).Format(time.RFC3339)

		}
		if v, ok := override["end_time"]; ok {
			new_override.EndTime = v.(string)

			loc, zone_err := time.LoadLocation(new_schedule.Time_zone)
			if zone_err != nil {
				return nil, diag.FromErr(errors.New(zone_err.Error()))
			}
			parsed_end_time, err := time.ParseInLocation("2006-01-02 15:04", new_override.EndTime, loc)

			if err != nil {
				return nil, diag.FromErr(errors.New(err.Error()))
			}

			new_override.EndTime = parsed_end_time.In(time.UTC).Format(time.RFC3339)

		}
		if v, ok := override["user"]; ok {

			new_override.User = v.(string)
		}
		Overrides[o] = new_override
	}
	return Overrides, nil
}

func createSchedule(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.CreateSchedule, diag.Diagnostics) {
	new_schedule := &client.CreateSchedule{}

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
		if emptyString(v.(string)) {
			return nil, diag.FromErr(errors.New("time_zone must not be empty"))
		}
		_, zone_err := time.LoadLocation(v.(string))
		if zone_err != nil {
			return nil, diag.FromErr(errors.New(zone_err.Error()))
		}
		new_schedule.Time_zone = v.(string)

	}
	if v, ok := d.GetOk("team_id"); ok {
		if emptyString(v.(string)) {
			return nil, diag.FromErr(errors.New("team_id must not be empty"))
		}
		new_schedule.Team = v.(string)

	}

	Layers, layerErr := buildScheduleLayer(new_schedule.Time_zone, Ctx, d)
	if layerErr != nil {
		return nil, layerErr
	}
	Overrides, overrideErr := buildScheduleOverride(new_schedule, d)
	if overrideErr != nil {
		return nil, overrideErr
	}

	new_schedule.Overrides = Overrides
	new_schedule.Layers = Layers

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
	d.Set("time_zone", service.Time_zone)
	d.Set("team_id", service.Team)
	if err := d.Set("layers", flattenLayer(service.Time_zone, service.Layers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("overrides", flattenScheduleOverrides(service.Time_zone, service.Overrides)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenLayer(TimeZone string, layers []client.Layers) []map[string]interface{} {

	var layer_list []map[string]interface{}
	for i, layer := range layers {
		if emptyString(layer.Name) {
			layer.Name = fmt.Sprintf("Layer-%d", i+1)
		}
		layer_list = append(layer_list, map[string]interface{}{
			"name":                layer.Name,
			"shift_length":        layer.ShiftLength,
			"rotation_start_time": createScheduleLayerTimeFormat(layer.RotationStartTime, TimeZone),
			"rotation_end_time":   createScheduleLayerTimeFormat(layer.RotationEndTime, TimeZone),
			"users":               flattenLayerUsers(layer.Users),
			"restriction_type":    layer.RestrictionType,
			"restrictions":        flattenLayerRestrictions(layer.Restrictions),
		})
	}
	return layer_list

}

func createScheduleLayerTimeFormat(timestamp, Zone string) string {
	RFC3339local := "2006-01-02T15:04:05Z"
	loc, zone_err := time.LoadLocation(Zone)
	if zone_err != nil {
		return timestamp
	}
	t, parse_err := time.ParseInLocation(RFC3339local, timestamp, time.UTC)
	if parse_err != nil {
		return timestamp
	}

	return t.In(loc).Format("2006-01-02 15:04")

}

func flattenLayerUsers(users []client.Users) []string {

	var user_list []string
	for _, user := range users {
		user_list = append(user_list, user.User)
	}
	return user_list

}

func flattenLayerRestrictions(restrictions []client.Restrictions) []map[string]interface{} {

	var restriction_list []map[string]interface{}
	for _, restriction := range restrictions {
		if restriction.Duration == 0 {
			restriction.Duration = 1
		}
		restriction_list = append(restriction_list, map[string]interface{}{
			"duration":          restriction.Duration,
			"start_day_of_week": restriction.StartDayOfWeek,
			"start_time_of_day": restriction.StartTimeOfDay,
		})
	}
	return restriction_list

}

func flattenScheduleOverrides(TimeZone string, overrides []client.Overrides) []map[string]interface{} {

	var override_list []map[string]interface{}
	for i, override := range overrides {
		if emptyString(override.Name) {
			override.Name = fmt.Sprintf("Override-%d", i+1)
		}
		override_list = append(override_list, map[string]interface{}{
			"name":       override.Name,
			"start_time": createScheduleLayerTimeFormat(override.StartTime, TimeZone),
			"end_time":   createScheduleLayerTimeFormat(override.EndTime, TimeZone),
			"user":       override.User,
		})
	}

	return override_list
}

func resourceScheduleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format. Expecting team_id/schedule_id, got: %s", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id: %s", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid schedule_id: %s", parts[1])
	}

	d.Set("team_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
