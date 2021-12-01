package zenduty

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentReads,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"results": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_reference": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"integration_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_incidents_for": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"integration_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_urgency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIncidentReads(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)
	integrationID := d.Get("integration_id").(string)
	var diags diag.Diagnostics
	if integrationID != "" {
		integration, err := apiclient.Integrations.GetIntegrationByID(teamID, serviceID, integrationID)
		if err != nil {
			return diag.Errorf("Error reading integrations: %s", err)
		}
		items := make([]map[string]interface{}, 1)
		items[0] = map[string]interface{}{

			"name":          integration.Name,
			"creation_date": integration.Creation_Date,
			"summary":       integration.Summary,
			"unique_id":     integration.Unique_Id,
			"service":       integration.Service,
			"application":   integration.Application,
			"application_reference": map[string]interface{}{
				"name":               integration.Application_Reference.Name,
				"icon_url":           integration.Application_Reference.Icon_Url,
				"summary":            integration.Application_Reference.Summary,
				"description":        integration.Application_Reference.Description,
				"unique_id":          integration.Application_Reference.Unique_Id,
				"setup_instructions": integration.Application_Reference.Setup_Instructions,
				"extension":          integration.Application_Reference.Extension,
				"categories":         integration.Application_Reference.Categories,
				"documentation_link": integration.Application_Reference.Documentation_Link,
			},
			"integration_key":      integration.Integration_key,
			"created_by":           integration.Created_By,
			"create_incidents_for": integration.Create_Incident_For,
			"integration_type":     integration.Integration_Type,
			"default_urgency":      integration.Default_Urgency,
		}

		if err := d.Set("results", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())
		return diags

	} else {

		integrations, err := apiclient.Integrations.GetIntegrations(teamID, serviceID)
		if err != nil {
			return diag.Errorf("Error reading integrations: %s", err)
		}
		items := make([]map[string]interface{}, len(integrations))
		for i, integration := range integrations {
			items[i] = map[string]interface{}{

				"name":          integration.Name,
				"creation_date": integration.Creation_Date,
				"summary":       integration.Summary,
				"unique_id":     integration.Unique_Id,
				"service":       integration.Service,
				"application":   integration.Application,
				"application_reference": map[string]interface{}{
					"name":               integration.Application_Reference.Name,
					"icon_url":           integration.Application_Reference.Icon_Url,
					"summary":            integration.Application_Reference.Summary,
					"description":        integration.Application_Reference.Description,
					"unique_id":          integration.Application_Reference.Unique_Id,
					"setup_instructions": integration.Application_Reference.Setup_Instructions,
					"extension":          integration.Application_Reference.Extension,
					"categories":         integration.Application_Reference.Categories,
					"documentation_link": integration.Application_Reference.Documentation_Link,
				},
				"integration_key":      integration.Integration_key,
				"created_by":           integration.Created_By,
				"create_incidents_for": integration.Create_Incident_For,
				"integration_type":     integration.Integration_Type,
				"default_urgency":      integration.Default_Urgency,
			}
		}

		if err := d.Set("results", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())
		return diags

	}

}
