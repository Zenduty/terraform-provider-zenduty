package zenduty

import (
	"context"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAssignAccountRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssignRole,
		ReadContext:   resourceReadRole,
		UpdateContext: resourceUpdateAssignRole,
		DeleteContext: resourceRemoveAssignedRole,
		Importer: &schema.ResourceImporter{
			State: resourceAssignAccountRoleImporter,
		},
		Schema: map[string]*schema.Schema{
			"account_role": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"username": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUserName(),
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
	return resourceReadRole(ctx, d, m)
}

func resourceUpdateAssignRole(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	role := &client.AddRoleToUser{}

	if v, ok := d.GetOk("account_role"); ok {
		AccountRole := v.(string)
		role.AccountRole = &AccountRole
	}
	username := d.Get("username").(string)

	_, err := apiclient.AccountRole.AssignRoleToUser(username, role)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceReadRole(Ctx, d, m)

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

// resourceReadRole refreshes the assignment from the API. A user without a
// custom role no longer has an assignment, so the resource leaves state and
// is planned for re-creation.
func resourceReadRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	username := d.Get("username").(string)
	if username == "" {
		username = d.Id()
	}
	userRole, err := apiclient.AccountRole.GetUserRole(username)
	if err != nil {
		return handleReadError(d, err)
	}
	if userRole.CustomRole == nil || *userRole.CustomRole == "" {
		d.SetId("")
		return nil
	}
	d.Set("username", username)
	d.Set("account_role", *userRole.CustomRole)
	return nil
}

func resourceAssignAccountRoleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if !IsValidUsername(d.Id()) {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <username>", d.Id())
	}
	d.Set("username", d.Id())
	return []*schema.ResourceData{d}, nil
}
