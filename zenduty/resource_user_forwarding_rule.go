package zenduty

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserForwardingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateUserForwardingRule,
		ReadContext:   resourceReadUserForwardingRule,
		DeleteContext: resourceDeleteUserForwardingRule,
		CustomizeDiff: validateUserForwardingRuleDiff,
		Importer: &schema.ResourceImporter{
			State: resourceUserForwardingRuleImporter,
		},
		Description: "Forwards a user's notifications to another user for a period; every attribute forces a new resource.",
		Schema: map[string]*schema.Schema{
			"username": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUserName(),
				Description:      "Username of the user whose notifications are forwarded.",
			},
			"to_user": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUserName(),
				Description:      "Username of the member who receives the notifications.",
			},
			"from_date": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateLocalTime(),
				Description:      "Start as YYYY-MM-DD HH:MM in time_zone.",
			},
			"to_date": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateLocalTime(),
				Description:      "End as YYYY-MM-DD HH:MM in time_zone.",
			},
			"time_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          "UTC",
				ValidateDiagFunc: ValidateTimeZone(),
			},
		},
	}
}

func expandUserForwardingRule(d *schema.ResourceData) (*client.ForwardingRule, error) {
	zone := d.Get("time_zone").(string)
	from, to, err := localPeriodToAPI("from_date", d.Get("from_date").(string), "to_date", d.Get("to_date").(string), zone)
	if err != nil {
		return nil, err
	}
	return &client.ForwardingRule{
		ToUser:   d.Get("to_user").(string),
		FromDate: from,
		ToDate:   to,
		TimeZone: zone,
	}, nil
}

func validateUserForwardingRuleDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	for _, key := range []string{"from_date", "to_date", "time_zone"} {
		if !diff.NewValueKnown(key) {
			return nil
		}
	}
	_, _, err := localPeriodToAPI("from_date", diff.Get("from_date").(string), "to_date", diff.Get("to_date").(string), diff.Get("time_zone").(string))
	return err
}

func resourceCreateUserForwardingRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	rule, err := expandUserForwardingRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	created, err := apiclient.ForwardingRule.CreateForwardingRule(d.Get("username").(string), rule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadUserForwardingRule(ctx, d, m)
}

func resourceReadUserForwardingRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	username := d.Get("username").(string)
	rule, err := apiclient.ForwardingRule.GetForwardingRule(username, d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	zone := rule.TimeZone
	if zone == "" {
		zone = d.Get("time_zone").(string)
	}
	if rule.User != "" {
		d.Set("username", rule.User)
	}
	d.Set("to_user", rule.ToUser)
	d.Set("time_zone", zone)
	for key, value := range map[string]string{"from_date": rule.FromDate, "to_date": rule.ToDate} {
		local, convErr := apiTimeToLocal(value, zone)
		if convErr != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", key, convErr))
		}
		d.Set(key, local)
	}
	return nil
}

func resourceDeleteUserForwardingRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.ForwardingRule.DeleteForwardingRule(d.Get("username").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserForwardingRuleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <username>/<forwarding_rule_id>", d.Id())
	} else if !IsValidUsername(parts[0]) {
		return nil, fmt.Errorf("invalid username (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid forwarding_rule_id (%q)", parts[1])
	}
	d.Set("username", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
