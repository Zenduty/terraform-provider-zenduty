package zenduty

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateManintenances,
		UpdateContext: resourceUpdateManintenances,
		DeleteContext: resourceDeleteManintenances,
		ReadContext:   resourceReadManintenances,
		Importer: &schema.ResourceImporter{
			State: resourceMaintenanceImporter,
		},
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repeat_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"repeat_until": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ValidateMaintenanceWindow(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.MaintenanceWindow, diag.Diagnostics) {
	new_manintence := &client.MaintenanceWindow{}
	services := d.Get("services").([]interface{})

	if v, ok := d.GetOk("name"); ok {
		if v.(string) == "" {
			return nil, diag.FromErr(errors.New("name must not be empty"))
		}
		new_manintence.Name = v.(string)
	}
	if v, ok := d.GetOk("timezone"); ok {
		if v.(string) == "" {
			return nil, diag.FromErr(errors.New("timezone must not be empty"))
		}
		new_manintence.TimeZone = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		if !validateDate(v.(string)) {
			return nil, diag.FromErr(errors.New("start_time is invalid"))
		}

		new_manintence.StartTime = v.(string)

		loc, zone_err := time.LoadLocation(new_manintence.TimeZone)
		if zone_err != nil {
			return nil, diag.FromErr(errors.New(zone_err.Error()))
		}
		parsed_time, err := time.ParseInLocation("2006-01-02 15:04", v.(string), loc)
		if err != nil {
			return nil, diag.FromErr(errors.New(err.Error()))
		}
		new_manintence.StartTime = parsed_time.In(time.UTC).Format(time.RFC3339)

	}
	if v, ok := d.GetOk("end_time"); ok {
		if !validateDate(v.(string)) {
			return nil, diag.FromErr(errors.New("end_time is invalid"))
		}
		new_manintence.EndTime = v.(string)

		loc, zone_err := time.LoadLocation(new_manintence.TimeZone)
		if zone_err != nil {
			return nil, diag.FromErr(errors.New(zone_err.Error()))
		}
		parsed_time, err := time.ParseInLocation("2006-01-02 15:04", v.(string), loc)
		if err != nil {
			return nil, diag.FromErr(errors.New(err.Error()))
		}

		new_manintence.EndTime = parsed_time.In(time.UTC).Format(time.RFC3339)
	}

	if v, ok := d.GetOk("repeat_interval"); ok {
		if v.(int) <= 0 {
			return nil, diag.FromErr(errors.New("repeat_interval must be greater than 0"))
		}
		new_manintence.RepeatInterval = v.(int)
	}
	if v, ok := d.GetOk("repeat_until"); ok {
		if v.(string) == "" {
			return nil, diag.FromErr(errors.New("repeat_until must not be empty"))
		}
		if !validateDate(v.(string)) {
			return nil, diag.FromErr(errors.New("repeat_until is invalid"))
		}

		loc, zone_err := time.LoadLocation(new_manintence.TimeZone)
		if zone_err != nil {
			return nil, diag.FromErr(errors.New(zone_err.Error()))
		}

		parsed_time, err := time.ParseInLocation("2006-01-02 15:04", v.(string), loc)
		if err != nil {
			return nil, diag.FromErr(errors.New(err.Error()))
		}
		new_manintence.RepeatUntil = parsed_time.In(time.UTC).Format(time.RFC3339)
	}
	for _, service := range services {
		if service.(string) == "" {
			return nil, diag.FromErr(errors.New("services must not be empty"))
		}
		if !IsValidUUID(service.(string)) {
			return nil, diag.FromErr(errors.New("services must be a valid UUID"))
		}
		new_manintence.Services = append(new_manintence.Services, client.ServiceMaintenance{Service: service.(string)})
	}
	return new_manintence, nil
}

// if v, ok := d.GetOk("services"); ok {

// 	for _, service := range v.([]interface{}) {
// 		if service.(map[string]interface{})["services"] == nil {
// 			return nil, diag.FromErr(errors.New("services must not be empty"))
// 		}
// 		if !IsValidUUID(service.(map[string]interface{})["services"].(string)) {
// 			return nil, diag.FromErr(errors.New("services must be a valid UUID"))
// 		}

// 	}
// }

func resourceCreateManintenances(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var team_id string
	if v, ok := d.GetOk("team_id"); ok {
		if !IsValidUUID(v.(string)) {
			return diag.FromErr(errors.New("team_id must be a valid UUID"))

		}
		team_id = v.(string)
	}
	apiclient, _ := m.(*Config).Client()
	new_manintence, diags := ValidateMaintenanceWindow(ctx, d, m)
	if diags != nil {
		return diags
	}
	maintenance, err := apiclient.MaintenanceWindow.CreateMaintenanceWindow(team_id, new_manintence)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(maintenance.UniqueID)
	return nil
}

func resourceUpdateManintenances(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var team_id string
	if v, ok := d.GetOk("team_id"); ok {
		if !IsValidUUID(v.(string)) {
			return diag.FromErr(errors.New("team_id must be a valid UUID"))
		}
		team_id = v.(string)
	}
	apiclient, _ := m.(*Config).Client()
	new_manintence, diags := ValidateMaintenanceWindow(ctx, d, m)
	if diags != nil {
		return diags
	}
	maintenance, err := apiclient.MaintenanceWindow.UpdateMaintenanceWindow(team_id, d.Id(), new_manintence)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(maintenance.UniqueID)
	return nil
}

func resourceReadManintenances(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var team_id string
	if v, ok := d.GetOk("team_id"); ok {
		if !IsValidUUID(v.(string)) {
			return diag.FromErr(errors.New("team_id must be a valid UUID"))
		}
		team_id = v.(string)
	}
	apiclient, _ := m.(*Config).Client()
	maintenance, err := apiclient.MaintenanceWindow.GetMaintenanceWindowById(team_id, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", maintenance.Name)
	d.Set("repeat_interval", maintenance.RepeatInterval)
	start_time, time_err := createMaintenanceTimeFormat(maintenance.StartTime, maintenance.TimeZone)
	end_time, end_time_err := createMaintenanceTimeFormat(maintenance.EndTime, maintenance.TimeZone)
	repeat_until, repeat_until_err := createMaintenanceTimeFormat(maintenance.RepeatUntil, maintenance.TimeZone)
	if repeat_until_err == nil {
		d.Set("repeat_until", repeat_until)
	}

	if time_err == nil && end_time_err == nil {
		d.Set("start_time", start_time)
		d.Set("end_time", end_time)
	}
	d.Set("services", flattenServices(maintenance.Services))
	d.Set("timezone", maintenance.TimeZone)

	return nil
}

func flattenServices(services []client.ServiceMaintenance) []interface{} {
	var services_list []interface{}
	for _, service := range services {
		services_list = append(services_list, service.Service)
	}
	return services_list
}

func createMaintenanceTimeFormat(timestamp, Zone string) (string, error) {
	RFC3339local := "2006-01-02T15:04:05Z"
	loc, zone_err := time.LoadLocation(Zone)
	if zone_err != nil {
		return "", zone_err
	}
	t, parse_err := time.ParseInLocation(RFC3339local, timestamp, time.UTC)
	if parse_err != nil {
		return "", parse_err
	}

	return t.In(loc).Format("2006-01-02 15:04"), nil

}

func resourceDeleteManintenances(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var team_id string
	if v, ok := d.GetOk("team_id"); ok {
		if !IsValidUUID(v.(string)) {
			return diag.FromErr(errors.New("team_id must be a valid UUID"))
		}
		team_id = v.(string)
	}
	apiclient, _ := m.(*Config).Client()
	err := apiclient.MaintenanceWindow.DeleteMaintenanceWindow(team_id, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceMaintenanceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<maintenance_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid maintenance (%q)", parts[1])
	}
	d.Set("team_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
