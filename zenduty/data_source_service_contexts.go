package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServiceContexts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceContextsRead,
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
			"service_contexts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id":     {Type: schema.TypeString, Computed: true},
						"integration":   {Type: schema.TypeString, Computed: true},
						"created_by":    {Type: schema.TypeString, Computed: true},
						"creation_date": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceServiceContextsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)

	contexts, err := apiclient.ServiceContext.GetServiceContexts(teamID, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(contexts))
	for i, c := range contexts {
		items[i] = map[string]interface{}{
			"unique_id":     c.UniqueID,
			"integration":   c.Integration,
			"created_by":    c.CreatedBy,
			"creation_date": c.CreationDate,
		}
	}
	if err := d.Set("service_contexts", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(serviceID)
	return nil
}
