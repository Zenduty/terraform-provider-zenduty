package zenduty

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchedules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScheduleReads,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedules": {
				Type:        schema.TypeList,
				Description: "List of schedules",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
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
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"team": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unique_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"shift_length": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rotation_start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rotation_end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_edited": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"restriction_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_active": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"users": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeString,
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
											},
										},
									},
									"restrictions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"duration": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"start_day_of_week": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"start_time_of_day": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"unique_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},

						"overrides": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"unique_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceScheduleReads(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	schedule_id := d.Get("schedule_id").(string)

	var diags diag.Diagnostics
	if schedule_id != "" {
		schedule, err := apiclient.Schedules.GetScheduleByID(team_id, schedule_id)
		if err != nil {
			return diag.FromErr(err)
		}
		items := make([]map[string]interface{}, 1)
		item := make(map[string]interface{})
		item["unique_id"] = schedule.Unique_Id
		item["name"] = schedule.Name
		item["summary"] = schedule.Summary
		item["description"] = schedule.Description
		item["time_zone"] = schedule.Time_zone
		item["team"] = schedule.Team
		layers := make([]map[string]interface{}, len(schedule.Layers))
		for j, layer := range schedule.Layers {
			layers[j] = map[string]interface{}{
				"shift_length":        layer.ShiftLength,
				"name":                layer.Name,
				"rotation_start_time": layer.RotationStartTime,
				"rotation_end_time":   layer.RotationEndTime,
				"unique_id":           layer.UniqueId,
				"last_edited":         layer.LastEdited,
				"restriction_type":    layer.RestrictionType,
				"is_active":           layer.IsActive,
			}
			if layer.Restrictions != nil {
				restrictions := make([]map[string]interface{}, len(layer.Restrictions))
				for k, restriction := range layer.Restrictions {
					restrictions[k] = map[string]interface{}{
						"duration":          restriction.Duration,
						"start_day_of_week": restriction.StartDayOfWeek,
						"start_time_of_day": restriction.StartTimeOfDay,
						"unique_id":         restriction.Unique_Id,
					}
				}
				layers[j]["restrictions"] = restrictions
			}
			if layer.Users != nil {
				users := make([]map[string]interface{}, len(layer.Users))
				for k, user := range layer.Users {
					users[k] = map[string]interface{}{
						"user": user.User,
					}
				}
				layers[j]["users"] = users
			}

		}
		item["layers"] = layers
		if schedule.Overrides != nil {
			overrides := make([]map[string]interface{}, len(schedule.Overrides))
			for j, override := range schedule.Overrides {
				overrides[j] = map[string]interface{}{
					"name":       override.Name,
					"user":       override.User,
					"start_time": override.StartTime,
					"end_time":   override.EndTime,
					"unique_id":  override.Unique_Id,
				}
			}
			item["overrides"] = overrides
		}
		items[0] = item

		if err := d.Set("schedules", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())

		return diags
	} else {

		schedules, err := apiclient.Schedules.GetSchedules(team_id)
		if err != nil {
			return diag.FromErr(err)
		}
		items := make([]map[string]interface{}, len(schedules))
		for i, schedule := range schedules {
			item := make(map[string]interface{})
			item["unique_id"] = schedule.Unique_Id
			item["name"] = schedule.Name
			item["summary"] = schedule.Summary
			item["description"] = schedule.Description
			item["time_zone"] = schedule.Time_zone
			item["team"] = schedule.Team
			layers := make([]map[string]interface{}, len(schedule.Layers))
			for j, layer := range schedule.Layers {
				layers[j] = map[string]interface{}{
					"shift_length":        layer.ShiftLength,
					"name":                layer.Name,
					"rotation_start_time": layer.RotationStartTime,
					"rotation_end_time":   layer.RotationEndTime,
					"unique_id":           layer.UniqueId,
					"last_edited":         layer.LastEdited,
					"restriction_type":    layer.RestrictionType,
					"is_active":           layer.IsActive,
				}
				if layer.Restrictions != nil {
					restrictions := make([]map[string]interface{}, len(layer.Restrictions))
					for k, restriction := range layer.Restrictions {
						restrictions[k] = map[string]interface{}{
							"duration":          restriction.Duration,
							"start_day_of_week": restriction.StartDayOfWeek,
							"start_time_of_day": restriction.StartTimeOfDay,
							"unique_id":         restriction.Unique_Id,
						}
					}
					layers[j]["restrictions"] = restrictions
				}
				if layer.Users != nil {
					users := make([]map[string]interface{}, len(layer.Users))
					for k, user := range layer.Users {
						users[k] = map[string]interface{}{
							"user": user.User,
						}
					}
					layers[j]["users"] = users
				}

			}
			item["layers"] = layers
			if schedule.Overrides != nil {
				overrides := make([]map[string]interface{}, len(schedule.Overrides))
				for j, override := range schedule.Overrides {
					overrides[j] = map[string]interface{}{
						"name":       override.Name,
						"user":       override.User,
						"start_time": override.StartTime,
						"end_time":   override.EndTime,
						"unique_id":  override.Unique_Id,
					}
				}
				item["overrides"] = overrides
			}
			items[i] = item
		}
		if err := d.Set("schedules", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())

		return diags
	}
}
