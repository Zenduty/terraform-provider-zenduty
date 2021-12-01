package zenduty

import (
	"context"
	"strconv"

	"github.com/Kdheeru12/zenduty-test/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIncidents() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentsCreate,
		UpdateContext: resourceIncidentUpdate,
		DeleteContext: resourceIncidentDelete,
		ReadContext:   resourceIncidentRead,
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"escalation_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"summary": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceIncidentsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	new_incident := &client.Incident{}
	var diags diag.Diagnostics
	summary := d.Get("summary").(string)
	if summary != "" {
		new_incident.Summary = summary
	}
	if v, ok := d.GetOk("title"); ok {
		new_incident.Title = v.(string)
	}
	if v, ok := d.GetOk("user"); ok {
		new_incident.User = v.(string)
	}
	if v, ok := d.GetOk("escalation_policy"); ok {
		new_incident.EscalationPolicy = v.(string)
	}
	if v, ok := d.GetOk("service"); ok {
		new_incident.Service = v.(string)
	}

	incident, err := apiclient.Incidents.CreateIncident(new_incident)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(incident.Incident_Number))
	return diags
}

func resourceIncidentUpdate(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	status := d.Get("status").(int)
	if status != 0 {
		apiclient, _ := m.(*Config).Client()

		id := d.Id()
		new_status := &client.IncidentStatus{}
		_, err := apiclient.Incidents.UpdateIncident(id, new_status)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	var diags diag.Diagnostics
	return diags
}

func resourceIncidentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceIncidentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
