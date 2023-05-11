package zenduty

import (
	"context"
	"errors"
	"regexp"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAssignAccountRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssignRole,
		ReadContext:   resourceReadRole,
		UpdateContext: resourceUpdateAssignRole,
		DeleteContext: resourceRemoveAssignedRole,
		Schema: map[string]*schema.Schema{
			"account_role": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{1}$`), "must be a valid username"),
			},
		},
	}
}

func resourceAssignRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	role := &client.AddRoleToUser{}

	if v, ok := d.GetOk("account_role"); ok {
		AccountRole := v.(string)
		role.AccountRole = &AccountRole
	}
	username := d.Get("username").(string)

	userrole, err := apiclient.AccountRole.AssignRoleToUser(username, role)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(userrole.Username)
	return nil
}

func resourceUpdateAssignRole(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	role := &client.AddRoleToUser{}

	var diags diag.Diagnostics

	if v, ok := d.GetOk("account_role"); ok {
		AccountRole := v.(string)
		role.AccountRole = &AccountRole
	}
	username := d.Get("username").(string)
	if username != d.Id() {
		return diag.FromErr(errors.New("cannot update username"))
	}

	_, err := apiclient.AccountRole.AssignRoleToUser(username, role)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags

}

func resourceRemoveAssignedRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	role := &client.AddRoleToUser{AccountRole: nil}

	var diags diag.Diagnostics
	username := d.Get("username").(string)
	_, err := apiclient.AccountRole.AssignRoleToUser(username, role)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceReadRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
