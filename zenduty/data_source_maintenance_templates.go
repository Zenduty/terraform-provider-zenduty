package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMaintenanceTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMaintenanceTemplatesRead,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"maintenance_templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id":  {Type: schema.TypeString, Computed: true},
						"name":       {Type: schema.TypeString, Computed: true},
						"is_enabled": {Type: schema.TypeBool, Computed: true},
						"rule_json":  {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceMaintenanceTemplatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	teamID := d.Get("team_id").(string)

	templates, err := apiclient.MaintenanceTemplate.GetMaintenanceTemplates(teamID)
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(templates))
	for i, t := range templates {
		items[i] = map[string]interface{}{
			"unique_id":  t.UniqueID,
			"name":       t.Name,
			"is_enabled": t.IsEnabled,
			"rule_json":  t.RuleJSON,
		}
	}
	if err := d.Set("maintenance_templates", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(teamID)
	return nil
}
