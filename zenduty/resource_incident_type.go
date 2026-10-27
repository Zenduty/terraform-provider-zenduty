package zenduty

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIncidentType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateIncidentType,
		ReadContext:   resourceReadIncidentType,
		UpdateContext: resourceUpdateIncidentType,
		DeleteContext: resourceDeleteIncidentType,
		Importer: &schema.ResourceImporter{
			State: resourceIncidentTypeImporter,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
				Description:      "Name of the incident type, unique within the account. Deleting a type only deactivates it, so its name stays reserved.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this is the account's default incident type. Setting it clears the flag on the previous default.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the type can be picked for new incidents. A type is always created active; false deactivates it right after creation.",
			},
			"form_toggles": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
				Description: "Built-in incident form sections shown for this type, e.g. { impacted_services = true, responders = false }.",
			},
			"custom_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom fields on this type's incident form, in display order.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_field": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ValidateUUID(),
							Description:      "unique_id of the zenduty_incident_custom_field.",
						},
						"position": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Display position. Defaults to the block's position in the list.",
						},
						"required": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"default_value": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ValidateJSONString(),
							DiffSuppressFunc: suppressEquivalentJSONDiffs,
							Description:      "Default value as JSON, e.g. jsonencode(\"eu-west\") or jsonencode([\"a\", \"b\"]) for a multiselect.",
						},
					},
				},
			},
		},
	}
}

func expandIncidentTypeInput(d *schema.ResourceData) *client.IncidentTypeInput {
	input := &client.IncidentTypeInput{
		Name:         d.Get("name").(string),
		IsDefault:    d.Get("is_default").(bool),
		FormToggles:  []map[string]interface{}{},
		CustomFields: []client.IncidentTypeFieldMapping{},
	}
	toggles := d.Get("form_toggles").(map[string]interface{})
	keys := make([]string, 0, len(toggles))
	for key := range toggles {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		enabled, _ := toggles[key].(bool)
		input.FormToggles = append(input.FormToggles, map[string]interface{}{key: enabled})
	}
	for index, raw := range d.Get("custom_fields").([]interface{}) {
		block, _ := raw.(map[string]interface{})
		mapping := client.IncidentTypeFieldMapping{
			UniqueID: block["custom_field"].(string),
			Required: block["required"].(bool),
			Position: index,
		}
		if position, _ := block["position"].(int); position != 0 {
			mapping.Position = position
		}
		if value, _ := block["default_value"].(string); value != "" {
			mapping.DefaultValue = json.RawMessage(value)
		}
		input.CustomFields = append(input.CustomFields, mapping)
	}
	return input
}

// incidentCustomFieldKeys maps every custom field's backend field_key to its
// unique_id; the type's read shape is keyed by field_key, the write shape by
// unique_id.
func incidentCustomFieldKeys(apiclient *client.Client) (map[string]string, error) {
	fields, err := apiclient.IncidentCustomField.GetIncidentCustomFields()
	if err != nil {
		return nil, fmt.Errorf("listing incident custom fields: %w", err)
	}
	keys := make(map[string]string, len(fields))
	for _, field := range fields {
		keys[field.FieldKey] = field.UniqueID
	}
	return keys, nil
}

// incidentTypeFieldEntry is one mapped custom field of an incident type in
// a shape both the resource and the data source can flatten.
type incidentTypeFieldEntry struct {
	FieldKey     string
	UniqueID     string
	Position     int
	Required     bool
	DefaultValue string
}

// incidentTypeFieldEntries rebuilds the mapped fields from the JSON Schema the
// API returns, in position order. Fields whose definition is gone are left out.
func incidentTypeFieldEntries(incidentType *client.IncidentType, keyToID map[string]string) []incidentTypeFieldEntry {
	required := make(map[string]bool, len(incidentType.CustomFields.Required))
	for _, key := range incidentType.CustomFields.Required {
		required[key] = true
	}
	entries := make([]incidentTypeFieldEntry, 0, len(incidentType.CustomFields.Properties))
	for key, field := range incidentType.CustomFields.Properties {
		uniqueID, known := keyToID[key]
		if !known {
			continue
		}
		entry := incidentTypeFieldEntry{FieldKey: key, UniqueID: uniqueID, Position: field.Position, Required: required[key]}
		if len(field.Default) > 0 && string(field.Default) != "null" {
			if normalized, err := normalizeJSON(string(field.Default)); err == nil {
				entry.DefaultValue = normalized
			} else {
				entry.DefaultValue = string(field.Default)
			}
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Position != entries[j].Position {
			return entries[i].Position < entries[j].Position
		}
		return entries[i].FieldKey < entries[j].FieldKey
	})
	return entries
}

// flattenIncidentTypeCustomFields builds the resource's custom_fields blocks.
func flattenIncidentTypeCustomFields(incidentType *client.IncidentType, keyToID map[string]string) []interface{} {
	entries := incidentTypeFieldEntries(incidentType, keyToID)
	blocks := make([]interface{}, 0, len(entries))
	for _, e := range entries {
		blocks = append(blocks, map[string]interface{}{
			"custom_field":  e.UniqueID,
			"position":      e.Position,
			"required":      e.Required,
			"default_value": e.DefaultValue,
		})
	}
	return blocks
}

// flattenFormToggles turns the API's single-key objects back into a map of
// booleans. A toggle decorated with widget config ({"enabled": true, ...})
// counts as enabled.
func flattenFormToggles(toggles []map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(toggles))
	for _, item := range toggles {
		for key, value := range item {
			switch v := value.(type) {
			case bool:
				out[key] = v
			case map[string]interface{}:
				enabled, _ := v["enabled"].(bool)
				out[key] = enabled
			default:
				out[key] = false
			}
		}
	}
	return out
}

func resourceCreateIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.IncidentType.CreateIncidentType(expandIncidentTypeInput(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	// new types are always active on the API; deactivate afterwards if asked
	if !d.Get("is_active").(bool) {
		inactive := false
		input := expandIncidentTypeInput(d)
		input.IsActive = &inactive
		if _, err := apiclient.IncidentType.UpdateIncidentType(created.UniqueID, input); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceReadIncidentType(ctx, d, m)
}

func resourceUpdateIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	input := expandIncidentTypeInput(d)
	active := d.Get("is_active").(bool)
	input.IsActive = &active
	if _, err := apiclient.IncidentType.UpdateIncidentType(d.Id(), input); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadIncidentType(ctx, d, m)
}

func resourceReadIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	incidentType, err := apiclient.IncidentType.GetIncidentType(d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	keyToID, err := incidentCustomFieldKeys(apiclient)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", incidentType.Name)
	d.Set("is_default", incidentType.IsDefault)
	d.Set("is_active", incidentType.IsActive)
	if err := d.Set("form_toggles", flattenFormToggles(incidentType.FormToggles)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_fields", flattenIncidentTypeCustomFields(incidentType, keyToID)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// Deleting an incident type soft-deletes it on the API (is_active = false).
func resourceDeleteIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.IncidentType.DeleteIncidentType(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceIncidentTypeImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if !IsValidUUID(d.Id()) {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <incident_type_id>", d.Id())
	}
	return []*schema.ResourceData{d}, nil
}
