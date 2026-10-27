package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserCustomRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserCustomRoleRead,
		Description: "The custom account role assigned to a user, if any. The API only serves this for members with the regular user role.",
		Schema: map[string]*schema.Schema{
			"username": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUserName(),
			},
			"account_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "unique_id of the assigned custom role, empty when the user has none.",
			},
			"role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserCustomRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	username := d.Get("username").(string)

	role, err := apiclient.AccountRole.GetUserRole(username)
	if err != nil {
		return diag.FromErr(err)
	}
	accountRole, roleName := "", ""
	if role.CustomRole != nil {
		accountRole = *role.CustomRole
	}
	if role.CustomRoleDetails != nil {
		roleName = role.CustomRoleDetails.Name
	}
	d.Set("account_role", accountRole)
	d.Set("role_name", roleName)
	d.SetId(username)
	return nil
}
