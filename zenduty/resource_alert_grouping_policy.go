package zenduty

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validStaticMatchFields = []string{"incident_title", "incident_summary", "integration"}

func resourceAlertGroupingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAlertGroupingPolicy,
		UpdateContext: resourceUpdateAlertGroupingPolicy,
		DeleteContext: resourceDeleteAlertGroupingPolicy,
		ReadContext:   resourceReadAlertGroupingPolicy,
		CustomizeDiff: validateAlertGroupingPolicyDiff,
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
				Description:      "The service this content-based collation policy applies to. The service must have collation = 3; create and update fail otherwise.",
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
				Description: "JSONPath expressions evaluated against the integration's alert payload (e.g. message, labels.host).",
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

func expandStringList(raw []interface{}) []string {
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		s, _ := v.(string)
		out = append(out, s)
	}
	return out
}

// API rules for match fields: at least one field, no field listed twice.
func validateAlertGroupingFields(static, custom []string) error {
	if len(static)+len(custom) == 0 {
		return errors.New("at least one of static_fields or custom_fields is required")
	}
	seen := make(map[string]string, len(static)+len(custom))
	check := func(attr string, fields []string) error {
		for _, f := range fields {
			if prev, dup := seen[f]; dup {
				if prev == attr {
					return fmt.Errorf("%s lists %q more than once", attr, f)
				}
				return fmt.Errorf("%q is listed in both %s and %s; a match field may appear only once", f, prev, attr)
			}
			seen[f] = attr
		}
		return nil
	}
	if err := check("static_fields", static); err != nil {
		return err
	}
	return check("custom_fields", custom)
}

// Plan-time field check; unknown values read as "" and are skipped.
func validateAlertGroupingPolicyDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	if !diff.NewValueKnown("static_fields") || !diff.NewValueKnown("custom_fields") {
		return nil
	}
	rawStatic, _ := diff.Get("static_fields").([]interface{})
	rawCustom, _ := diff.Get("custom_fields").([]interface{})
	if len(rawStatic)+len(rawCustom) == 0 {
		return errors.New("at least one of static_fields or custom_fields is required")
	}
	known := func(fields []string) []string {
		out := fields[:0]
		for _, f := range fields {
			if f != "" {
				out = append(out, f)
			}
		}
		return out
	}
	static := known(expandStringList(rawStatic))
	custom := known(expandStringList(rawCustom))
	if len(static)+len(custom) == 0 {
		return nil
	}
	return validateAlertGroupingFields(static, custom)
}

func validateAndCreateAlertGroupingPolicy(d *schema.ResourceData) (*client.AlertGroupingPolicy, error) {
	staticFields := expandStringList(d.Get("static_fields").([]interface{}))
	customFields := expandStringList(d.Get("custom_fields").([]interface{}))
	if err := validateAlertGroupingFields(staticFields, customFields); err != nil {
		return nil, err
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

// A policy only takes effect when the service's collation is 3.
func requireContentBasedCollation(apiclient *client.Client, teamID, serviceID string) error {
	service, err := apiclient.Services.GetServicesByID(teamID, serviceID)
	if err != nil {
		return fmt.Errorf("looking up service %s: %w", serviceID, err)
	}
	if service.Collation != collationContentBased {
		return fmt.Errorf("service %s has collation = %d; set collation = %d (with a collation_time) on its zenduty_services resource to enable content-based collation, otherwise this alert grouping policy is ignored",
			serviceID, service.Collation, collationContentBased)
	}
	return nil
}

// Warns when the service left collation 3 while the policy still exists.
func warnIfNotContentBased(apiclient *client.Client, teamID, serviceID string) diag.Diagnostics {
	service, err := apiclient.Services.GetServicesByID(teamID, serviceID)
	if err != nil {
		return diag.Diagnostics{{
			Severity:      diag.Warning,
			Summary:       "Could not verify the service's collation mode",
			Detail:        fmt.Sprintf("Looking up service %s failed: %s", serviceID, err),
			AttributePath: cty.GetAttrPath("service_id"),
		}}
	}
	if service.Collation == collationContentBased {
		return nil
	}
	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "Alert grouping policy is not in effect",
		Detail: fmt.Sprintf("Service %s has collation = %d, so this policy is ignored. Set collation = %d on its zenduty_services resource to re-enable content-based collation.",
			serviceID, service.Collation, collationContentBased),
		AttributePath: cty.GetAttrPath("service_id"),
	}}
}

// Warns when a collation 3 service has no active policy (alerts not grouped).
func warnIfNoActiveGroupingPolicy(apiclient *client.Client, teamID, serviceID string) diag.Diagnostics {
	policies, err := apiclient.AlertGroupingPolicy.GetAlertGroupingPolicies(teamID, serviceID)
	if err != nil {
		return diag.Diagnostics{{
			Severity:      diag.Warning,
			Summary:       "Could not verify the service's alert grouping policy",
			Detail:        fmt.Sprintf("Listing alert grouping policies for service %s failed: %s", serviceID, err),
			AttributePath: cty.GetAttrPath("collation"),
		}}
	}
	for _, p := range policies {
		if p.IsActive {
			return nil
		}
	}
	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "Content-based collation has no active alert grouping policy",
		Detail: fmt.Sprintf("Service %s has collation = %d but no active alert grouping policy, so its alerts are not grouped. Add a zenduty_alert_grouping_policy for this service; if one is declared in this configuration, the warning clears once it is applied.",
			serviceID, collationContentBased),
		AttributePath: cty.GetAttrPath("collation"),
	}}
}

func resourceCreateAlertGroupingPolicy(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	teamID := d.Get("team_id").(string)
	serviceID := d.Get("service_id").(string)

	policy, err := validateAndCreateAlertGroupingPolicy(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := requireContentBasedCollation(apiclient, teamID, serviceID); err != nil {
		return diag.FromErr(err)
	}

	// One policy per service; look first so the error can name the one to import.
	existing, err := apiclient.AlertGroupingPolicy.GetAlertGroupingPolicies(teamID, serviceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("listing alert grouping policies for service %s: %w", serviceID, err))
	}
	if len(existing) > 0 {
		return diag.Errorf("service %s already has alert grouping policy %s (for example created from the web console); a service can have only one. Import it instead of creating another:\n  terraform import zenduty_alert_grouping_policy.<name> %s/%s/%s",
			serviceID, existing[0].UniqueID, teamID, serviceID, existing[0].UniqueID)
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

	if err := requireContentBasedCollation(apiclient, teamID, serviceID); err != nil {
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

	return warnIfNotContentBased(apiclient, teamID, serviceID)
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
