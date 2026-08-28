package zenduty

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// validStaticMatchFields mirrors the built-in options the dashboard offers
// for content-based collation: incident title, incident summary, and the
// alert's originating integration.
var validStaticMatchFields = []string{"incident_title", "incident_summary", "integration"}

func resourceAlertGroupingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAlertGroupingPolicy,
		UpdateContext: resourceUpdateAlertGroupingPolicy,
		DeleteContext: resourceDeleteAlertGroupingPolicy,
		ReadContext:   resourceReadAlertGroupingPolicy,
		Importer: &schema.ResourceImporter{
			State: resourceAlertGroupingPolicyImporter,
		},

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"service_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
				Description:      "The service this content-based collation policy applies to. Set the service's collation attribute to 3 to enable content-based collation.",
			},
			"match_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 2),
				Description:  "1 requires only one of the listed fields (static_fields and custom_fields) to match, 2 requires every listed field to match.",
			},
			"static_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validStaticMatchFields, false),
				},
				Description: "Built-in alert attributes to compare: incident_title, incident_summary, integration.",
			},
			"custom_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "JSON keys from the integration's alert payload to compare (e.g. message).",
			},
			"time_window": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 20),
				Description:  "Minutes to wait for matching alerts before closing correlation eligibility.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether this grouping policy is enforced.",
			},
		},
	}
}

func validateAndCreateAlertGroupingPolicy(d *schema.ResourceData) (*client.AlertGroupingPolicy, error) {
	staticFieldsRaw := d.Get("static_fields").([]interface{})
	customFieldsRaw := d.Get("custom_fields").([]interface{})
	if len(staticFieldsRaw) == 0 && len(customFieldsRaw) == 0 {
		return nil, errors.New("at least one of static_fields or custom_fields is required")
	}

	staticFields := make([]string, len(staticFieldsRaw))
	for i, v := range staticFieldsRaw {
		staticFields[i] = v.(string)
	}
	customFields := make([]string, len(customFieldsRaw))
	for i, v := range customFieldsRaw {
		customFields[i] = v.(string)
	}

	return &client.AlertGroupingPolicy{
		Service:   d.Get("service_id").(string),
		MatchMode: d.Get("match_mode").(int),
		MatchFields: client.AlertGroupingMatchFields{
			StaticFields: staticFields,
			CustomFields: customFields,
		},
		TimeWindow: d.Get("time_window").(int),
		IsActive:   d.Get("is_active").(bool),
	}, nil
}

func resourceCreateAlertGroupingPolicy(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)

	policy, err := validateAndCreateAlertGroupingPolicy(d)
	if err != nil {
		return diag.FromErr(err)
	}

	created, err := apiclient.AlertGroupingPolicy.CreateAlertGroupingPolicy(teamID, serviceID, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.UniqueID)
	return resourceReadAlertGroupingPolicy(Ctx, d, m)
}

func resourceUpdateAlertGroupingPolicy(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)

	policy, err := validateAndCreateAlertGroupingPolicy(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = apiclient.AlertGroupingPolicy.UpdateAlertGroupingPolicy(teamID, serviceID, d.Id(), policy)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceReadAlertGroupingPolicy(Ctx, d, m)
}

func resourceDeleteAlertGroupingPolicy(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)

	var diags diag.Diagnostics
	err := apiclient.AlertGroupingPolicy.DeleteAlertGroupingPolicy(teamID, serviceID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceReadAlertGroupingPolicy(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)
	if emptyString(teamID) || emptyString(serviceID) {
		return diag.FromErr(errors.New("team_id and service_id are required"))
	}

	var diags diag.Diagnostics
	policy, err := apiclient.AlertGroupingPolicy.GetAlertGroupingPolicy(teamID, serviceID, d.Id())
	if err != nil {
		return handleReadError(d, err)
	}

	d.Set("team_id", teamID)
	d.Set("service_id", policy.Service)
	d.Set("match_mode", policy.MatchMode)
	d.Set("static_fields", policy.MatchFields.StaticFields)
	d.Set("custom_fields", policy.MatchFields.CustomFields)
	d.Set("time_window", policy.TimeWindow)
	d.Set("is_active", policy.IsActive)

	return diags
}

func resourceAlertGroupingPolicyImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<service_id>/<alert_grouping_policy_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid service_id (%q)", parts[1])
	} else if !IsValidUUID(parts[2]) {
		return nil, fmt.Errorf("invalid alert_grouping_policy_id (%q)", parts[2])
	}
	d.Set("team_id", parts[0])
	d.Set("service_id", parts[1])
	d.SetId(parts[2])
	return []*schema.ResourceData{d}, nil
}
