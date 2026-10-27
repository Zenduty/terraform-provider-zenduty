package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccountSetting() *schema.Resource {
	flag := func() *schema.Schema { return &schema.Schema{Type: schema.TypeBool, Computed: true} }
	return &schema.Resource{
		ReadContext: dataSourceAccountSettingRead,
		Schema: map[string]*schema.Schema{
			"enable_delete_incident":        flag(),
			"enable_team_manager_invite":    flag(),
			"enable_regen_key":              flag(),
			"enable_first_mtta_calculation": flag(),
		},
	}
}

func dataSourceAccountSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	setting, err := apiclient.AccountSetting.GetAccountSetting()
	if err != nil {
		return diag.FromErr(err)
	}
	flag := func(v *bool) bool { return v != nil && *v }
	d.Set("enable_delete_incident", flag(setting.EnableDeleteIncident))
	d.Set("enable_team_manager_invite", flag(setting.EnableTeamManagerInvite))
	d.Set("enable_regen_key", flag(setting.EnableRegenKey))
	d.Set("enable_first_mtta_calculation", flag(setting.EnableFirstMTTACalculation))
	d.SetId(accountSettingID)
	return nil
}
