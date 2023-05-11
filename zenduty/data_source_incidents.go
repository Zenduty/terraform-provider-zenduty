package zenduty

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncidents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentRead,
		Schema: map[string]*schema.Schema{
			"number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"results": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_auto_resolve_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_object": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_team_priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_task_template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_acknowledgement_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_object_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_object_escalation_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_team": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_sla": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_object_collation_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_object_collation": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urgency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"merged_with": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assigned_to": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"esccalation_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"escalation_policy_object_unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"escalation_policy_object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assigned_to_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resolved_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acknowledged_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"context_window_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"context_window_end": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// "tags": {
						// 	Type:     schema.TypeList,
						// 	Computed: true,
						// },
						"sla": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sla_object": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"team_priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"team_priority_object": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIncidentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	var diags diag.Diagnostics
	incidents, err := apiclient.Incidents.GetIncidents()
	if err != nil {
		return diag.FromErr(err)
	}
	results := incidents.Results
	items := make([]map[string]interface{}, len(results))
	for i, result := range results {

		item := make(map[string]interface{})
		item["summary"] = result.Summary
		item["incident_number"] = result.IncidentNumber
		item["creation_date"] = result.CreationDate
		item["status"] = result.Status
		item["unique_id"] = result.UniqueID
		item["service_object_name"] = result.ServiceObject.Name
		item["service_object_unique_id"] = result.ServiceObject.UniqueID
		item["service_object_creation_date"] = result.ServiceObject.CreationDate
		item["service_object_status"] = result.ServiceObject.Status
		item["service_object_team"] = result.ServiceObject.Team
		item["service_object_summary"] = result.ServiceObject.Summary
		item["service_object_description"] = result.ServiceObject.Description
		item["service_object_acknowledgement_timeout"] = result.ServiceObject.AcknowledgmentTimeout
		item["service_object_auto_resolve_timeout"] = result.ServiceObject.AutoResolveTimeouts
		item["service_object_created_by"] = result.ServiceObject.CreatedBy
		item["service_object_team_priority"] = result.ServiceObject.TeamPriority
		item["service_object_task_template"] = result.ServiceObject.TaskTemplate
		item["service_object_escalation_policy"] = result.ServiceObject.EscalationPolicy
		item["service_object_team"] = result.ServiceObject.Team
		item["service_object_sla"] = result.ServiceObject.SLA
		item["service_object_collation_time"] = result.ServiceObject.CollationTime
		item["service_object_collation"] = result.ServiceObject.Collation
		item["team_priority"] = result.TeamPriority
		item["team_priority_object"] = result.TeamPriorityObject
		item["title"] = result.Title
		item["incident_key"] = result.IncidentKey
		item["service"] = result.Service
		item["urgency"] = result.Urgency
		items[i] = item
	}
	if err := d.Set("results", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(time.Now().String())
	return diags

}
