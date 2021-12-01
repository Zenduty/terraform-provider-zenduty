package zenduty

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServicesRead,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_resolve_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acknowledgement_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"under_maintenance": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"escalation_policy": {
							Type:     schema.TypeString,
							Required: true,
						},
						"team": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"summary": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"collation": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"collation_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"sla": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"task_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"team_priority": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	id := d.Get("id").(string)
	if team_id == "" {
		return diag.Errorf("team_id is required")
	}
	var diags diag.Diagnostics
	if id != "" {
		service, err := apiclient.Services.GetServicesById(team_id, id)
		if err != nil {
			return diag.FromErr(err)
		}
		items := make([]map[string]interface{}, 1)

		items[0] = map[string]interface{}{
			"name":                    service.Name,
			"escalation_policy":       service.Escalation_Policy,
			"team":                    service.Team,
			"description":             service.Description,
			"summary":                 service.Summary,
			"collation":               service.Collation,
			"collation_time":          service.Collation_Time,
			"sla":                     service.Sla,
			"task_template":           service.Task_Template,
			"team_priority":           service.Team_Priority,
			"creation_date":           service.Creation_Date,
			"unique_id":               service.Unique_Id,
			"auto_resolve_timeout":    service.Auto_Resolve_Timeout,
			"created_by":              service.Created_By,
			"acknowledgement_timeout": service.Acknowledgment_Timeout,
			"status":                  service.Status,
			"under_maintenance":       service.Under_Maintenance,
		}
		if err := d.Set("services", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())
		return diags
	} else {

		// id := d.Get("id").(string)
		services, err := apiclient.Services.GetServices(team_id)
		if err != nil {
			return diag.FromErr(err)
		}
		items := make([]map[string]interface{}, len(services))
		for i, service := range services {
			items[i] = map[string]interface{}{
				"name":                    service.Name,
				"escalation_policy":       service.Escalation_Policy,
				"team":                    service.Team,
				"description":             service.Description,
				"summary":                 service.Summary,
				"collation":               service.Collation,
				"collation_time":          service.Collation_Time,
				"sla":                     service.Sla,
				"task_template":           service.Task_Template,
				"team_priority":           service.Team_Priority,
				"creation_date":           service.Creation_Date,
				"unique_id":               service.Unique_Id,
				"auto_resolve_timeout":    service.Auto_Resolve_Timeout,
				"created_by":              service.Created_By,
				"acknowledgement_timeout": service.Acknowledgment_Timeout,
				"status":                  service.Status,
				"under_maintenance":       service.Under_Maintenance,
			}
		}
		if err := d.Set("services", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())

		return diags
	}
}
