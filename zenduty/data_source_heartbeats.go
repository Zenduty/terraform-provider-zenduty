package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHeartbeats() *schema.Resource {
	str := func() *schema.Schema { return &schema.Schema{Type: schema.TypeString, Computed: true} }
	return &schema.Resource{
		ReadContext: dataSourceHeartbeatsRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ValidateUUID(),
				Description:      "Only return heartbeats of this service.",
			},
			"heartbeats": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id":              str(),
						"name":                   str(),
						"description":            str(),
						"service_id":             str(),
						"service_name":           str(),
						"team_id":                str(),
						"team_name":              str(),
						"interval_seconds":       {Type: schema.TypeInt, Computed: true},
						"status":                 str(),
						"is_enabled":             {Type: schema.TypeBool, Computed: true},
						"alert_message":          str(),
						"alert_urgency":          {Type: schema.TypeInt, Computed: true, Description: "0 low, 1 high, -1 when the heartbeat uses the API default."},
						"ping_type":              {Type: schema.TypeInt, Computed: true},
						"email_subject_contains": str(),
						"ping_url":               str(),
						"ping_email":             str(),
						"last_ping_at":           str(),
						"expires_at":             str(),
					},
				},
			},
		},
	}
}

func dataSourceHeartbeatsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	serviceID := d.Get("service_id").(string)

	heartbeats, err := apiclient.Heartbeats.GetHeartbeats()
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, 0, len(heartbeats))
	for _, h := range heartbeats {
		if serviceID != "" && h.Service != serviceID {
			continue
		}
		urgency := -1
		if h.AlertUrgency != nil {
			urgency = *h.AlertUrgency
		}
		items = append(items, map[string]interface{}{
			"unique_id":              h.UniqueID,
			"name":                   h.Name,
			"description":            h.Description,
			"service_id":             h.Service,
			"service_name":           h.ServiceName,
			"team_id":                h.Team,
			"team_name":              h.TeamName,
			"interval_seconds":       h.IntervalSeconds,
			"status":                 h.Status,
			"is_enabled":             h.IsEnabled,
			"alert_message":          h.AlertMessage,
			"alert_urgency":          urgency,
			"ping_type":              h.PingType,
			"email_subject_contains": h.EmailSubjectContains,
			"ping_url":               h.PingURL,
			"ping_email":             h.PingEmail,
			"last_ping_at":           h.LastPingAt,
			"expires_at":             h.ExpiresAt,
		})
	}
	if err := d.Set("heartbeats", items); err != nil {
		return diag.FromErr(err)
	}
	if serviceID != "" {
		d.SetId(serviceID)
	} else {
		d.SetId("heartbeats")
	}
	return nil
}
