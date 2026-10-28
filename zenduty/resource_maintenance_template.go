package zenduty

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMaintenanceTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateMaintenanceTemplate,
		ReadContext:   resourceReadMaintenanceTemplate,
		UpdateContext: resourceUpdateMaintenanceTemplate,
		DeleteContext: resourceDeleteMaintenanceTemplate,
		Importer: &schema.ResourceImporter{
			State: resourceMaintenanceTemplateImporter,
		},
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"rule_json": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateJSONObjectString(),
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
				Description:      "Alert-matching rule in the zenduty_alertrules rule_json format.",
			},
		},
	}
}

func expandMaintenanceTemplate(d *schema.ResourceData) *client.MaintenanceTemplate {
	return &client.MaintenanceTemplate{
		Name:      d.Get("name").(string),
		IsEnabled: d.Get("is_enabled").(bool),
		RuleJSON:  d.Get("rule_json").(string),
	}
}

func resourceCreateMaintenanceTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.MaintenanceTemplate.CreateMaintenanceTemplate(d.Get("team_id").(string), expandMaintenanceTemplate(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadMaintenanceTemplate(ctx, d, m)
}

func resourceUpdateMaintenanceTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if _, err := apiclient.MaintenanceTemplate.UpdateMaintenanceTemplate(d.Get("team_id").(string), d.Id(), expandMaintenanceTemplate(d)); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadMaintenanceTemplate(ctx, d, m)
}

func resourceReadMaintenanceTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	template, err := apiclient.MaintenanceTemplate.GetMaintenanceTemplate(d.Get("team_id").(string), d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	d.Set("name", template.Name)
	d.Set("is_enabled", template.IsEnabled)
	if normalized, normErr := normalizeJSON(template.RuleJSON); normErr == nil {
		d.Set("rule_json", normalized)
	} else {
		d.Set("rule_json", template.RuleJSON)
	}
	return nil
}

func resourceDeleteMaintenanceTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.MaintenanceTemplate.DeleteMaintenanceTemplate(d.Get("team_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceMaintenanceTemplateImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<maintenance_template_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid maintenance_template_id (%q)", parts[1])
	}
	d.Set("team_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
