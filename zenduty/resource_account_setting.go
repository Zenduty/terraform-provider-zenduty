package zenduty

import (
	"context"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// accountSettingID is the fixed id of the singleton.
const accountSettingID = "account"

func resourceAccountSetting() *schema.Resource {
	setting := func(description string) *schema.Schema {
		return &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: description,
		}
	}
	return &schema.Resource{
		CreateContext: resourceUpsertAccountSetting,
		ReadContext:   resourceReadAccountSetting,
		UpdateContext: resourceUpsertAccountSetting,
		DeleteContext: resourceDeleteAccountSetting,
		Importer: &schema.ResourceImporter{
			State: resourceAccountSettingImporter,
		},
		Description: "The account's singleton settings; destroy only removes it from state.",
		Schema: map[string]*schema.Schema{
			"enable_delete_incident":        setting("Whether incidents can be deleted."),
			"enable_team_manager_invite":    setting("Whether team managers can invite users."),
			"enable_regen_key":              setting("Whether team managers can regenerate integration keys."),
			"enable_first_mtta_calculation": setting("Whether MTTA is measured to the first acknowledgement."),
		},
	}
}

func resourceUpsertAccountSetting(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	setting := &client.AccountSetting{
		EnableDeleteIncident:       configuredBool(d, "enable_delete_incident"),
		EnableTeamManagerInvite:    configuredBool(d, "enable_team_manager_invite"),
		EnableRegenKey:             configuredBool(d, "enable_regen_key"),
		EnableFirstMTTACalculation: configuredBool(d, "enable_first_mtta_calculation"),
	}
	if setting.EnableDeleteIncident != nil || setting.EnableTeamManagerInvite != nil ||
		setting.EnableRegenKey != nil || setting.EnableFirstMTTACalculation != nil {
		if _, err := apiclient.AccountSetting.UpdateAccountSetting(setting); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(accountSettingID)
	return resourceReadAccountSetting(ctx, d, m)
}

func resourceReadAccountSetting(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	setting, err := apiclient.AccountSetting.GetAccountSetting()
	if err != nil {
		return handleReadError(d, err)
	}
	flag := func(v *bool) bool { return v != nil && *v }
	d.Set("enable_delete_incident", flag(setting.EnableDeleteIncident))
	d.Set("enable_team_manager_invite", flag(setting.EnableTeamManagerInvite))
	d.Set("enable_regen_key", flag(setting.EnableRegenKey))
	d.Set("enable_first_mtta_calculation", flag(setting.EnableFirstMTTACalculation))
	return nil
}

func resourceDeleteAccountSetting(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceAccountSettingImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.SetId(accountSettingID)
	return []*schema.ResourceData{d}, nil
}
