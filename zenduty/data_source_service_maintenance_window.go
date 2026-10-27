package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServiceMaintenanceWindows() *schema.Resource {
	str := func() *schema.Schema { return &schema.Schema{Type: schema.TypeString, Computed: true} }
	return &schema.Resource{
		ReadContext: dataSourceServiceMaintenanceWindowsRead,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"service_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"maintenance_windows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id":       str(),
						"name":            str(),
						"start_time":      str(),
						"end_time":        str(),
						"timezone":        str(),
						"repeat_interval": {Type: schema.TypeInt, Computed: true},
						"repeat_until":    str(),
						"creation_date":   str(),
					},
				},
			},
		},
	}
}

// localOrRaw shows an API timestamp as YYYY-MM-DD HH:MM in zone, falling
// back to the raw value when it cannot be converted.
func localOrRaw(value, zone string) string {
	if local, err := apiTimeToLocal(value, zone); err == nil {
		return local
	}
	return value
}

func dataSourceServiceMaintenanceWindowsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)

	windows, err := apiclient.MaintenanceWindow.GetServiceMaintenanceWindows(teamID, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(windows))
	for i, w := range windows {
		zone := w.TimeZone
		if zone == "" {
			zone = "UTC"
		}
		items[i] = map[string]interface{}{
			"unique_id":       w.UniqueID,
			"name":            w.Name,
			"start_time":      localOrRaw(w.StartTime, zone),
			"end_time":        localOrRaw(w.EndTime, zone),
			"timezone":        zone,
			"repeat_interval": w.RepeatInterval,
			"repeat_until":    localOrRaw(w.RepeatUntil, zone),
			"creation_date":   w.CreationDate,
		}
	}
	if err := d.Set("maintenance_windows", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(serviceID)
	return nil
}
