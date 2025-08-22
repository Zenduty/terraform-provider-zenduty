package zenduty

import (
	"context"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMembers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberCreate,
		ReadContext:   wrapReadWith404(resourceMemberRead),
		UpdateContext: resourceMemberUpdate,
		DeleteContext: resourceMemberDelete,
		Schema: map[string]*schema.Schema{
			"team": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},
	}
}

func resourceMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	newMembers := &client.Member{}
	role := d.Get("role").(int)
	if role == 0 {
		newMembers.Role = 2
	} else {
		newMembers.Role = role
	}
	var diags diag.Diagnostics
	if v, ok := d.GetOk("team"); ok {
		newMembers.Team = v.(string)

	}
	if v, ok := d.GetOk("user"); ok {
		newMembers.User = v.(string)
	}

	member, err := apiclient.Members.CreateTeamMember(newMembers.Team, newMembers)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(member.UniqueID)
	return diags
}

func resourceMemberUpdate(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	newMembers := &client.Member{}
	id := d.Id()
	newMembers.UniqueID = id
	var diags diag.Diagnostics
	if v, ok := d.GetOk("user"); ok {
		newMembers.User = v.(string)
	}
	if v, ok := d.GetOk("role"); ok {

		if v.(int) == 0 {
			newMembers.Role = 2
		} else {
			newMembers.Role = v.(int)
		}
	}
	if v, ok := d.GetOk("team"); ok {
		newMembers.Team = v.(string)
	}
	_, err := apiclient.Members.UpdateTeamMember(newMembers)
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
	d.SetId(member.UniqueID)
	d.Set("team", member.Team)
	d.Set("user", member.User)
	d.Set("role", member.Role)

	return diags
}
