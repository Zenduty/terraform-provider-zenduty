package zenduty

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// incidentFieldSchemaTypes is the JSON Schema type each field_type requires.
var incidentFieldSchemaTypes = map[int]string{
	client.IncidentFieldTypeString:      "string",
	client.IncidentFieldTypeInteger:     "integer",
	client.IncidentFieldTypeSelect:      "string",
	client.IncidentFieldTypeMultiselect: "array",
	client.IncidentFieldTypeURL:         "string",
}

func incidentFieldTypes() []int {
	types := make([]int, 0, len(incidentFieldSchemaTypes))
	for t := range incidentFieldSchemaTypes {
		types = append(types, t)
	}
	return types
}

func resourceIncidentCustomField() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateIncidentCustomField,
		ReadContext:   resourceReadIncidentCustomField,
		UpdateContext: resourceUpdateIncidentCustomField,
		DeleteContext: resourceDeleteIncidentCustomField,
		CustomizeDiff: validateIncidentCustomFieldDiff,
		Importer: &schema.ResourceImporter{
			State: resourceIncidentCustomFieldImporter,
		},
		Schema: map[string]*schema.Schema{
			"field_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice(incidentFieldTypes()),
				Description:  "0 string, 2 integer, 3 select, 4 multiselect, 5 url.",
			},
			"properties": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateJSONObjectString(),
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
				Description:      "JSON Schema property descriptor, e.g. jsonencode({type = \"string\", title = \"Region\"}).",
			},
			"field_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backend-generated key (custom_field_<n>) under which the field appears in incident payloads.",
			},
		},
	}
}

// validateIncidentCustomFieldProperties mirrors the API's descriptor checks at plan time.
func validateIncidentCustomFieldProperties(fieldType int, properties string) error {
	var prop map[string]interface{}
	if err := json.Unmarshal([]byte(properties), &prop); err != nil || prop == nil {
		return errors.New("properties must be a JSON object")
	}
	title, _ := prop["title"].(string)
	if title == "" || len(title) > 256 {
		return errors.New("properties needs a non-empty title of at most 256 characters")
	}
	expected, ok := incidentFieldSchemaTypes[fieldType]
	if !ok {
		return fmt.Errorf("unsupported field_type %d", fieldType)
	}
	if got, _ := prop["type"].(string); got != expected {
		return fmt.Errorf("field_type %d requires properties.type %q, got %q", fieldType, expected, got)
	}
	switch fieldType {
	case client.IncidentFieldTypeSelect:
		if options, _ := prop["enum"].([]interface{}); len(options) == 0 {
			return errors.New("a select field needs a non-empty enum list in properties")
		}
	case client.IncidentFieldTypeMultiselect:
		items, _ := prop["items"].(map[string]interface{})
		if options, _ := items["enum"].([]interface{}); len(options) == 0 {
			return errors.New("a multiselect field needs items.enum, a non-empty list, in properties")
		}
	}
	return nil
}

func validateIncidentCustomFieldDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	if !diff.NewValueKnown("properties") || !diff.NewValueKnown("field_type") {
		return nil
	}
	return validateIncidentCustomFieldProperties(diff.Get("field_type").(int), diff.Get("properties").(string))
}

func expandIncidentCustomField(d *schema.ResourceData) *client.IncidentCustomField {
	return &client.IncidentCustomField{
		FieldType:  d.Get("field_type").(int),
		Properties: json.RawMessage(d.Get("properties").(string)),
	}
}

func resourceCreateIncidentCustomField(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.IncidentCustomField.CreateIncidentCustomField(expandIncidentCustomField(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadIncidentCustomField(ctx, d, m)
}

func resourceUpdateIncidentCustomField(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if _, err := apiclient.IncidentCustomField.UpdateIncidentCustomField(d.Id(), expandIncidentCustomField(d)); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadIncidentCustomField(ctx, d, m)
}

func resourceReadIncidentCustomField(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	field, err := apiclient.IncidentCustomField.GetIncidentCustomField(d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	d.Set("field_type", field.FieldType)
	d.Set("field_key", field.FieldKey)
	if err := setJSONAttribute(d, "properties", field.Properties); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDeleteIncidentCustomField(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.IncidentCustomField.DeleteIncidentCustomField(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceIncidentCustomFieldImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if !IsValidUUID(d.Id()) {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <custom_field_id>", d.Id())
	}
	return []*schema.ResourceData{d}, nil
}
