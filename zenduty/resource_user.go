package zenduty

import (
	"context"
	"errors"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateUser,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUpdateUser,
		DeleteContext: resourceDeleteUser,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"team": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: ValidateRequired(),
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: ValidateEmail(),
			},
			"role": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntBetween(2, 3),
				Default: 3,
			},
		},
	}
}

func resourceCreateUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	team := d.Get("team").(string)
	first_name := d.Get("first_name").(string)
	last_name := d.Get("last_name").(string)
	if emptyString(team) {
		return diag.FromErr(errors.New("team is required"))
	}
	if emptyString(last_name){
		return diag.FromErr(errors.New("last_name is required"))
	}
	email := d.Get("email").(string)
	apiclient, _ := m.(*Config).Client()
	new_user := &client.UserObj{First_Name: first_name, Last_Name: last_name, Email: email, Role: 3}
	new_userobj := &client.CreateUser{Team: team, User: *new_user}

	user, err := apiclient.Users.CreateUser(new_userobj)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(user.User.Username)
	d.Set("role", 3)
	return nil
}

func resourceUpdateUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	role := d.Get("role").(int)
	first_name := d.Get("first_name").(string)
	last_name := d.Get("last_name").(string)
	email := d.Get("email").(string)
	apiclient, _ := m.(*Config).Client()

	new_user := &client.UserObj{First_Name: first_name, Last_Name: last_name, Email: email, Role: role}
	
	user, err := apiclient.Users.UpdateUser(d.Id(),new_user)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(user.User.Username)
	d.Set("role", user.Role)
	d.Set("first_name", user.User.First_Name)
	d.Set("last_name", user.User.Last_Name)
	d.Set("email", user.User.Email)
	return nil
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiclient, _ := m.(*Config).Client()	
	user, err := apiclient.Users.GetUser(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(user.User.Username)
	d.Set("role", user.Role)
	d.Set("first_name", user.User.First_Name)
	d.Set("last_name", user.User.Last_Name)
	d.Set("email", user.User.Email)
	return nil
}

func resourceDeleteUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
