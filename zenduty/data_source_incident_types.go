package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncidentTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentTypesRead,
		Schema: map[string]*schema.Schema{
			"incident_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id":  {Type: schema.TypeString, Computed: true},
						"name":       {Type: schema.TypeString, Computed: true},
						"is_default": {Type: schema.TypeBool, Computed: true},
						"is_active":  {Type: schema.TypeBool, Computed: true},
						"form_toggles": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeBool},
						},
						"custom_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_field":  {Type: schema.TypeString, Computed: true},
									"field_key":     {Type: schema.TypeString, Computed: true},
									"position":      {Type: schema.TypeInt, Computed: true},
									"required":      {Type: schema.TypeBool, Computed: true},
									"default_value": {Type: schema.TypeString, Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIncidentTypesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	types, err := apiclient.IncidentType.GetIncidentTypes()
	if err != nil {
		return diag.FromErr(err)
	}
	keyToID, err := incidentCustomFieldKeys(apiclient)
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(types))
	for i := range types {
		t := &types[i]
		fields := make([]map[string]interface{}, 0)
		for _, e := range incidentTypeFieldEntries(t, keyToID) {
			fields = append(fields, map[string]interface{}{
				"custom_field":  e.UniqueID,
				"field_key":     e.FieldKey,
				"position":      e.Position,
				"required":      e.Required,
				"default_value": e.EffectiveDefault,
			})
		}
		items[i] = map[string]interface{}{
			"unique_id":     t.UniqueID,
			"name":          t.Name,
			"is_default":    t.IsDefault,
			"is_active":     t.IsActive,
			"form_toggles":  flattenFormToggles(t.FormToggles),
			"custom_fields": fields,
		}
	}
	if err := d.Set("incident_types", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("incident_types")
	return nil
}
