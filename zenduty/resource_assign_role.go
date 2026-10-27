package zenduty

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
// is planned for re-creation. The API answers 403 for account admins and
// owners, who cannot hold a custom role: that assignment no longer applies
// either, so it leaves state with a warning instead of failing every plan.
func resourceReadRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	username := d.Get("username").(string)
	if username == "" {
		username = d.Id()
	}
	userRole, err := apiclient.AccountRole.GetUserRole(username)
	if client.StatusCode(err) == http.StatusForbidden {
		log.Printf("[WARN] Removing zenduty_assign_account_role for %s from state: the API refused to serve the user's custom role (403)", username)
		d.SetId("")
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Custom role assignment no longer applies",
			Detail: fmt.Sprintf("The API refuses to serve the custom role of user %s (HTTP 403). Only members with the regular user role can hold one, so the user was most likely promoted to account admin or owner (or the token lacks admin rights). The assignment was removed from state; remove the resource from the configuration or demote the user before the next apply.",
				username),
		}}
	}
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
