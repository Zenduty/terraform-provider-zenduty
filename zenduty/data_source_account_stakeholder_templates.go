package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccountStakeholderTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountStakeholderTemplatesRead,
		Schema: map[string]*schema.Schema{
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": {Type: schema.TypeString, Computed: true},
						"title":     {Type: schema.TypeString, Computed: true},
						"subject":   {Type: schema.TypeString, Computed: true},
						"template":  {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceAccountStakeholderTemplatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	templates, err := apiclient.AccountStakeholderTemplate.GetAccountStakeholderTemplates()
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(templates))
	for i, t := range templates {
		items[i] = map[string]interface{}{
			"unique_id": t.UniqueID,
			"title":     t.Title,
			"subject":   t.Subject,
			"template":  t.Template,
		}
	}
	if err := d.Set("templates", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("account_stakeholder_templates")
	return nil
}
