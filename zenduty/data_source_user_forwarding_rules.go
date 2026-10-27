package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserForwardingRules() *schema.Resource {
	str := func() *schema.Schema { return &schema.Schema{Type: schema.TypeString, Computed: true} }
	return &schema.Resource{
		ReadContext: dataSourceUserForwardingRulesRead,
		Description: "Lists the forwarding rules involving a user: the rules forwarding their notifications and the rules forwarding other users' notifications to them. Compare username with the queried user to tell them apart.",
		Schema: map[string]*schema.Schema{
			"username": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUserName(),
			},
			"forwarding_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": str(),
						"username":  str(),
						"to_user":   str(),
						"from_date": str(),
						"to_date":   str(),
						"time_zone": str(),
					},
				},
			},
		},
	}
}

func dataSourceUserForwardingRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	username := d.Get("username").(string)

	rules, err := apiclient.ForwardingRule.GetForwardingRules(username)
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(rules))
	for i, r := range rules {
		zone := r.TimeZone
		if zone == "" {
			zone = "UTC"
		}
		items[i] = map[string]interface{}{
			"unique_id": r.UniqueID,
			"username":  r.User,
			"to_user":   r.ToUser,
			"from_date": localOrRaw(r.FromDate, zone),
			"to_date":   localOrRaw(r.ToDate, zone),
			"time_zone": zone,
		}
	}
	if err := d.Set("forwarding_rules", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(username)
	return nil
}
