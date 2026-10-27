package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncidentCustomFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentCustomFieldsRead,
		Schema: map[string]*schema.Schema{
			"custom_fields": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id":  {Type: schema.TypeString, Computed: true},
						"field_key":  {Type: schema.TypeString, Computed: true},
						"field_type": {Type: schema.TypeInt, Computed: true},
						"title":      {Type: schema.TypeString, Computed: true},
						"properties": {Type: schema.TypeString, Computed: true, Description: "The field's JSON Schema property descriptor as a JSON string."},
					},
				},
			},
		},
	}
}

func dataSourceIncidentCustomFieldsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	fields, err := apiclient.IncidentCustomField.GetIncidentCustomFields()
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(fields))
	for i, f := range fields {
		properties := string(f.Properties)
		if normalized, normErr := normalizeJSON(properties); normErr == nil {
			properties = normalized
		}
		items[i] = map[string]interface{}{
			"unique_id":  f.UniqueID,
			"field_key":  f.FieldKey,
			"field_type": f.FieldType,
			"title":      jsonStringField(f.Properties, "title"),
			"properties": properties,
		}
	}
	if err := d.Set("custom_fields", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("incident_custom_fields")
	return nil
}
