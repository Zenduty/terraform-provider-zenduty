package zenduty

import (
	"context"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		ReadContext:   resourceRoleRead,
		Schema: map[string]*schema.Schema{
			"team": {
				Type:     schema.TypeString,
				Required: true,
			},
			"unique_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rank": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	newrole := &client.Roles{}
	var diags diag.Diagnostics
	if v, ok := d.GetOk("team"); ok {
		newrole.Team = v.(string)

	}
	d.Set("team", newrole.Team)
	if v, ok := d.GetOk("description"); ok {
		newrole.Description = v.(string)

	}
	if v, ok := d.GetOk("title"); ok {
		newrole.Title = v.(string)
	}
	if v, ok := d.GetOk("rank"); ok {
		newrole.Rank = v.(int)
		if newrole.Rank == 0 {
			newrole.Rank = 1
		}
		if newrole.Rank <= 0 || newrole.Rank > 10 {
			return diag.Errorf("Rank should be between 1 and 10")
		}
	}

	role, err := apiclient.Roles.CreateRole(newrole.Team, newrole)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(role.Unique_Id)
	return diags
}

func resourceRoleUpdate(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	newrole := &client.Roles{}
	var team_id string
	id := d.Id()
	newrole.Unique_Id = id
	var diags diag.Diagnostics
	if v, ok := d.GetOk("description"); ok {
		newrole.Description = v.(string)

	}
	if v, ok := d.GetOk("title"); ok {
		newrole.Title = v.(string)
	}
	if v, ok := d.GetOk("team"); ok {
		team_id = v.(string)
	}
	if v, ok := d.GetOk("rank"); ok {
		newrole.Rank = v.(int)
		if newrole.Rank == 0 {
			newrole.Rank = 1
		}
		if newrole.Rank <= 0 || newrole.Rank > 10 {
			return diag.Errorf("Rank should be between 1 and 10")
		}
	}
	_, err := apiclient.Roles.UpdateRoles(team_id, newrole)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiclient, _ := m.(*Config).Client()

	id := d.Id()
	team_id := d.Get("team").(string)
	var diags diag.Diagnostics

	err := apiclient.Roles.DeleteRole(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
