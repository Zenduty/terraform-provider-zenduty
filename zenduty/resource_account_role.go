package zenduty

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccountRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAccountRole,
		ReadContext:   wrapReadWith404(resourceReadAccountRole),
		UpdateContext: resourceUpdateAccountRole,
		DeleteContext: resourceDeleteAccountRole,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func validateAccountRoles(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.AccountRole, diag.Diagnostics) {
	permissionsList := []string{
		"analytics_read",
		"escalation_policy_read",
		"escalation_policy_write",
		"incident_read",
		"incident_write",
		"incident_role_read",
		"incident_role_write",
		"integration_read",
		"integration_write",
		"maintenance_read",
		"maintenance_write",
		"member_read",
		"member_write",
		"postmortem_read",
		"postmortem_write",
		"priority_read",
		"priority_write",
		"schedule_read",
		"schedule_write",
		"service_read",
		"service_write",
		"sla_read",
		"sla_write",
		"stakeholder_template_read",
		"stakeholder_template_write",
		"tag_read",
		"tag_write",
		"task_template_read",
		"task_template_write",
		"team_read",
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	permissions := d.Get("permissions").([]interface{})
	newRole := &client.AccountRole{}

	newRole.Name = name
	newRole.Description = description
	for _, permission := range permissions {
		if permission.(string) == "" {
			return nil, diag.FromErr(errors.New("permission must not be empty"))
		}
		if !checkList(permission.(string), permissionsList) {
			return nil, diag.FromErr(fmt.Errorf("invalid permission received %s", permission.(string)))
		}
		newRole.Permissions = append(newRole.Permissions, permission.(string))
	}

	return newRole, nil
}

func resourceCreateAccountRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	newrole, validationerr := validateAccountRoles(ctx, d, m)
	if validationerr != nil {
		return validationerr
	}
	role, err := apiclient.AccountRole.CreateAccountRole(newrole)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(role.UniqueID)
	return nil
}

func resourceUpdateAccountRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	newrole, validationerr := validateAccountRoles(ctx, d, m)
	if validationerr != nil {
		return validationerr
	}
	role, err := apiclient.AccountRole.UpdateAccountRole(d.Id(), newrole)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(role.UniqueID)
	return nil
}

func resourceDeleteAccountRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	err := apiclient.AccountRole.DeleteAccountRole(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceReadAccountRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	role, err := apiclient.AccountRole.GetAccountRoleByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(role.UniqueID)
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("permissions", flattenPermissions(role.Permissions))
	return nil
}

func flattenPermissions(permissions []string) []string {
	var permissionList []string
	permissionList = append(permissionList, permissions...)
	return permissionList
}
