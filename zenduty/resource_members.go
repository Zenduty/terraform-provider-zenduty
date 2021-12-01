package zenduty

import (
	"context"

	"github.com/Kdheeru12/zenduty-test/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMembers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberCreate,
		ReadContext:   resourceMemberRead,
		UpdateContext: resourceMemberUpdate,
		DeleteContext: resourceMemberDelete,
		Schema: map[string]*schema.Schema{
			"team": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	new_members := &client.Member{}
	role := d.Get("role").(int)
	if role == 0 {
		new_members.Role = 2
	} else {
		new_members.Role = role
	}
	var diags diag.Diagnostics
	if v, ok := d.GetOk("team"); ok {
		new_members.Team = v.(string)

	}
	if v, ok := d.GetOk("user"); ok {
		new_members.User = v.(string)
	}

	member, err := apiclient.Members.CreateTeamMember(new_members.Team, new_members)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(member.Unique_Id)
	return diags
}

func resourceMemberUpdate(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	new_members := &client.Member{}
	id := d.Id()
	new_members.Unique_Id = id
	var diags diag.Diagnostics
	if v, ok := d.GetOk("user"); ok {
		new_members.User = v.(string)
	}
	if v, ok := d.GetOk("role"); ok {
		new_members.Role = v.(int)
	}
	if v, ok := d.GetOk("team"); ok {
		new_members.Team = v.(string)
	}
	_, err := apiclient.Members.UpdateTeamMember(new_members)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags

}

func resourceMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	id := d.Id()
	team := d.Get("team").(string)
	var diags diag.Diagnostics
	err := apiclient.Members.DeleteTeamMember(team, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	id := d.Id()
	team := d.Get("team").(string)
	var diags diag.Diagnostics
	member, err := apiclient.Members.GetTeamMembersByID(team, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(member.Unique_Id)
	d.Set("team", member.Team)
	d.Set("user", member.User)
	d.Set("role", member.Role)
	return diags
}
