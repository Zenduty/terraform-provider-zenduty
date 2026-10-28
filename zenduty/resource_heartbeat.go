package zenduty

import (
	"context"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceHeartbeat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateHeartbeat,
		ReadContext:   resourceReadHeartbeat,
		UpdateContext: resourceUpdateHeartbeat,
		DeleteContext: resourceDeleteHeartbeat,
		Importer: &schema.ResourceImporter{
			State: resourceHeartbeatImporter,
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
				Description:      "The service paged when the heartbeat expires; it must have an escalation policy.",
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
				Description:      "Name of the heartbeat, unique within the service.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"interval_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(client.HeartbeatMinIntervalSeconds, client.HeartbeatMaxIntervalSeconds),
				Description:  "Seconds without a ping before the heartbeat expires and pages the service, 60 to 7776000 (90 days).",
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"alert_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Message of the incident raised when the heartbeat expires.",
			},
			"alert_urgency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Description:  "Urgency of the incident raised: 0 low, 1 high. Unset leaves it to the API default.",
			},
			"ping_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      client.HeartbeatPingWebhook,
				ValidateFunc: validation.IntBetween(client.HeartbeatPingWebhook, client.HeartbeatPingEmail),
				Description:  "How the heartbeat is pinged: 0 webhook (ping_url), 1 email (ping_email).",
			},
			"email_subject_contains": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "For email pings, only emails whose subject contains this text count as a ping.",
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ping_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ping_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "disabled, pending (no ping yet), healthy or expired.",
			},
		},
	}
}

// configuredInt is the attribute's value when set in config, else nil.
func configuredInt(d *schema.ResourceData, key string) *int {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	attr := raw.GetAttr(key)
	if attr.IsNull() || !attr.IsKnown() {
		return nil
	}
	v := d.Get(key).(int)
	return &v
}

func expandHeartbeat(d *schema.ResourceData) *client.Heartbeat {
	return &client.Heartbeat{
		Service:              d.Get("service_id").(string),
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		IntervalSeconds:      d.Get("interval_seconds").(int),
		IsEnabled:            d.Get("is_enabled").(bool),
		AlertMessage:         d.Get("alert_message").(string),
		AlertUrgency:         configuredInt(d, "alert_urgency"),
		PingType:             d.Get("ping_type").(int),
		EmailSubjectContains: d.Get("email_subject_contains").(string),
	}
}

func resourceCreateHeartbeat(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.Heartbeats.CreateHeartbeat(expandHeartbeat(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadHeartbeat(ctx, d, m)
}

func resourceUpdateHeartbeat(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	heartbeat := expandHeartbeat(d)
	// the service is fixed after creation and ignored by the API on update
	heartbeat.Service = ""
	if _, err := apiclient.Heartbeats.UpdateHeartbeat(d.Id(), heartbeat); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadHeartbeat(ctx, d, m)
}

func resourceReadHeartbeat(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	heartbeat, err := apiclient.Heartbeats.GetHeartbeat(d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	d.Set("service_id", heartbeat.Service)
	d.Set("name", heartbeat.Name)
	d.Set("description", heartbeat.Description)
	d.Set("interval_seconds", heartbeat.IntervalSeconds)
	d.Set("is_enabled", heartbeat.IsEnabled)
	d.Set("alert_message", heartbeat.AlertMessage)
	if heartbeat.AlertUrgency != nil {
		d.Set("alert_urgency", *heartbeat.AlertUrgency)
	}
	d.Set("ping_type", heartbeat.PingType)
	d.Set("email_subject_contains", heartbeat.EmailSubjectContains)
	d.Set("team_id", heartbeat.Team)
	d.Set("ping_url", heartbeat.PingURL)
	d.Set("ping_email", heartbeat.PingEmail)
	d.Set("status", heartbeat.Status)
	return nil
}

func resourceDeleteHeartbeat(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.Heartbeats.DeleteHeartbeat(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceHeartbeatImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if !IsValidUUID(d.Id()) {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <heartbeat_id>", d.Id())
	}
	return []*schema.ResourceData{d}, nil
}
