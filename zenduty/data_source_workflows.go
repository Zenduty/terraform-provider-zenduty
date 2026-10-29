package zenduty

import (
	"context"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkflows() *schema.Resource {
	str := func() *schema.Schema { return &schema.Schema{Type: schema.TypeString, Computed: true} }
	return &schema.Resource{
		ReadContext: dataSourceWorkflowsRead,
		Description: "The account's workflows; without workflow_id the list shape omits json_schema and service_ids.",
		Schema: map[string]*schema.Schema{
			"workflow_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ValidateUUID(),
				Description:      "Return only this workflow, with its full json_schema and team scope.",
			},
			"workflows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id":   str(),
						"name":        str(),
						"description": str(),
						"is_draft":    {Type: schema.TypeBool, Computed: true},
						"is_enabled":  {Type: schema.TypeBool, Computed: true},
						"json_schema": str(),
						"teams": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"team_id": str(),
									"name":    str(),
									"service_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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

func flattenWorkflowForDataSource(w *client.Workflow) map[string]interface{} {
	teams := make([]map[string]interface{}, 0, len(w.Teams))
	for _, team := range w.Teams {
		services := make([]string, 0, len(team.Services))
		for _, service := range team.Services {
			services = append(services, service.UniqueID)
		}
		teams = append(teams, map[string]interface{}{
			"team_id":     team.UniqueID,
			"name":        team.Name,
			"service_ids": services,
		})
	}
	return map[string]interface{}{
		"unique_id":   w.UniqueID,
		"name":        w.Name,
		"description": w.Description,
		"is_draft":    w.IsDraft,
		"is_enabled":  w.IsEnabled,
		"json_schema": normalizeRawJSON([]byte(w.JSONSchema)),
		"teams":       teams,
	}
}

func dataSourceWorkflowsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	workflowID := d.Get("workflow_id").(string)

	var items []map[string]interface{}
	if workflowID != "" {
		workflow, err := apiclient.Workflow.GetWorkflow(workflowID)
		if err != nil {
			return diag.FromErr(err)
		}
		items = []map[string]interface{}{flattenWorkflowForDataSource(workflow)}
		d.SetId(workflowID)
	} else {
		workflows, err := apiclient.Workflow.GetWorkflows()
		if err != nil {
			return diag.FromErr(err)
		}
		items = make([]map[string]interface{}, len(workflows))
		for i := range workflows {
			items[i] = flattenWorkflowForDataSource(&workflows[i])
		}
		d.SetId("workflows")
	}
	if err := d.Set("workflows", items); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
