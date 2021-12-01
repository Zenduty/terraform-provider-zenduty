package zenduty

import (
	"context"
	"time"

	"github.com/Kdheeru12/zenduty-test/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInvite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInviteCreate,
		UpdateContext: resourceInviteUpdate,
		DeleteContext: resourceInviteDelete,
		ReadContext:   resourceInviteRead,
		Schema: map[string]*schema.Schema{
			"team": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_accounts": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceInviteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	newinvite := &client.Invite{}
	if v, ok := d.GetOk("team"); ok {
		newinvite.Team = v.(string)
	}
	email_accounts := d.Get("email_accounts").([]interface{})
	var diags diag.Diagnostics
	newinvite.EmailAccounts = make([]client.EmailAccounts, len(email_accounts))
	for i, user := range email_accounts {
		email_account := user.(map[string]interface{})
		if v, ok := email_account["email"]; ok {
			newinvite.EmailAccounts[i].Email = v.(string)
		}
		if v, ok := email_account["first_name"]; ok {
			newinvite.EmailAccounts[i].First_Name = v.(string)
		}
		if v, ok := email_account["last_name"]; ok {
			newinvite.EmailAccounts[i].Last_Name = v.(string)
		}
		if v, ok := email_account["role"]; ok {
			newinvite.EmailAccounts[i].Role = v.(int)
		}

	}
	_, err := apiclient.Invite.CreateInvite(newinvite)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(time.Now().String())

	return diags
}

func resourceInviteUpdate(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceInviteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceInviteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags

}
