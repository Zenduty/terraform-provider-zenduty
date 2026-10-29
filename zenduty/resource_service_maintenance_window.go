package zenduty

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceServiceMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateServiceMaintenanceWindow,
		ReadContext:   resourceReadServiceMaintenanceWindow,
		UpdateContext: resourceUpdateServiceMaintenanceWindow,
		DeleteContext: resourceDeleteServiceMaintenanceWindow,
		Importer: &schema.ResourceImporter{
			State: resourceServiceMaintenanceWindowImporter,
		},
		Description: "A maintenance window on a single service.",
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"service_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
			},
			"start_time": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateLocalTime(),
				Description:      "Start of the window as YYYY-MM-DD HH:MM in timezone.",
			},
			"end_time": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateLocalTime(),
				Description:      "End of the window as YYYY-MM-DD HH:MM in timezone.",
			},
			"timezone": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateTimeZone(),
				Description:      "IANA time zone the times are expressed in, e.g. UTC or Europe/Berlin.",
			},
			"repeat_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "Repeat interval of the window. 0 (the default) means the window does not repeat.",
			},
			"repeat_until": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ValidateLocalTime(),
				Description:      "Last repetition as YYYY-MM-DD HH:MM in timezone.",
			},
		},
	}
}

func expandServiceMaintenanceWindow(d *schema.ResourceData) (*client.ServiceMaintenanceWindow, error) {
	zone := d.Get("timezone").(string)
	start, end, err := localPeriodToAPI("start_time", d.Get("start_time").(string), "end_time", d.Get("end_time").(string), zone)
	if err != nil {
		return nil, err
	}
	window := &client.ServiceMaintenanceWindow{
		Name:           d.Get("name").(string),
		StartTime:      start,
		EndTime:        end,
		TimeZone:       zone,
		RepeatInterval: d.Get("repeat_interval").(int),
	}
	if v := d.Get("repeat_until").(string); v != "" {
		until, err := localTimeToAPI(v, zone)
		if err != nil {
			return nil, fmt.Errorf("repeat_until: %w", err)
		}
		window.RepeatUntil = &until
	}
	return window, nil
}

func resourceCreateServiceMaintenanceWindow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	window, err := expandServiceMaintenanceWindow(d)
	if err != nil {
		return diag.FromErr(err)
	}
	created, err := apiclient.MaintenanceWindow.CreateServiceMaintenanceWindow(d.Get("team_id").(string), d.Get("service_id").(string), window)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadServiceMaintenanceWindow(ctx, d, m)
}

func resourceUpdateServiceMaintenanceWindow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	window, err := expandServiceMaintenanceWindow(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := apiclient.MaintenanceWindow.UpdateServiceMaintenanceWindow(d.Get("team_id").(string), d.Get("service_id").(string), d.Id(), window); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadServiceMaintenanceWindow(ctx, d, m)
}

func resourceReadServiceMaintenanceWindow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)
	window, err := apiclient.MaintenanceWindow.GetServiceMaintenanceWindowByID(teamID, serviceID, d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	zone := window.TimeZone
	if zone == "" {
		zone = d.Get("timezone").(string)
	}
	d.Set("name", window.Name)
	d.Set("timezone", zone)
	d.Set("repeat_interval", window.RepeatInterval)
	if window.Service != "" {
		d.Set("service_id", window.Service)
	}
	repeatUntil := ""
	if window.RepeatUntil != nil {
		repeatUntil = *window.RepeatUntil
	}
	for key, value := range map[string]string{"start_time": window.StartTime, "end_time": window.EndTime, "repeat_until": repeatUntil} {
		local, convErr := apiTimeToLocal(value, zone)
		if convErr != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", key, convErr))
		}
		d.Set(key, local)
	}
	return nil
}

func resourceDeleteServiceMaintenanceWindow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.MaintenanceWindow.DeleteServiceMaintenanceWindow(d.Get("team_id").(string), d.Get("service_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceServiceMaintenanceWindowImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<service_id>/<maintenance_window_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid service_id (%q)", parts[1])
	} else if !IsValidUUID(parts[2]) {
		return nil, fmt.Errorf("invalid maintenance_window_id (%q)", parts[2])
	}
	d.Set("team_id", parts[0])
	d.Set("service_id", parts[1])
	d.SetId(parts[2])
	return []*schema.ResourceData{d}, nil
}
