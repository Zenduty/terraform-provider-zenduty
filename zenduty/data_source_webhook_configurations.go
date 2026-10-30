package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWebhookConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebhookConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"webhook_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": {Type: schema.TypeString, Computed: true},
						"url":       {Type: schema.TypeString, Computed: true},
						"event": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceWebhookConfigurationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	configs, err := apiclient.WebhookConfiguration.GetWebhookConfigurations()
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(configs))
	for i, c := range configs {
		items[i] = map[string]interface{}{
			"unique_id": c.UniqueID,
			"url":       c.URL,
			"event":     c.Event,
			"filter":    c.Filter,
		}
	}
	if err := d.Set("webhook_configurations", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("webhook_configurations")
	return nil
}
