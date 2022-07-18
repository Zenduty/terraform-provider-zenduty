package zenduty

import (
	"context"
	"errors"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePriority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreatePriority,
		UpdateContext: resourceUpdatePriority,
		DeleteContext: resourceDeletePriority,
		ReadContext:   resourceReadPriority,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"color": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func validatePriority(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.Priority, diag.Diagnostics) {
	name := d.Get("name").(string)
	color := d.Get("color").(string)
	team := d.Get("team_id").(string)
	description := d.Get("description").(string)

	new_priority := &client.Priority{}
	if !IsValidUUID(team) {
		return nil, diag.FromErr(errors.New("team_id must be a valid UUID"))
	}
	if color != "" && !checkList(color, []string{"magenta", "red", "volcano", "orange", "gold", "lime", "green", "cyan", "blue", "geekblue", "purple"}) {
		return nil, diag.FromErr(errors.New("color must be one of the following: magenta, red, volcano, orange, gold, lime, green, cyan, blue, geekblue, purple"))
	}

	new_priority.Name = name
	new_priority.Color = color
	new_priority.Description = description

	return new_priority, nil
}

func resourceCreatePriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := d.Get("team_id").(string)
	apiclient, _ := m.(*Config).Client()
	newpriority, validationerr := validatePriority(ctx, d, m)
	if validationerr != nil {
		return validationerr
	}
	tag, err := apiclient.Priority.CreatePriority(team, newpriority)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tag.Unique_Id)
	return nil
}

func resourceUpdatePriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := d.Get("team_id").(string)
	apiclient, _ := m.(*Config).Client()
	newpriority, validationerr := validatePriority(ctx, d, m)
	if validationerr != nil {
		return validationerr
	}
	tag, err := apiclient.Priority.UpdatePriority(team, d.Id(), newpriority)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tag.Unique_Id)
	return nil
}

func resourceDeletePriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := d.Get("team_id").(string)
	apiclient, _ := m.(*Config).Client()

	err := apiclient.Priority.DeletePriority(team, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceReadPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	team := d.Get("team_id").(string)
	tag, err := apiclient.Priority.GetPriorityById(team, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tag.Unique_Id)
	d.Set("name", tag.Name)
	d.Set("color", tag.Color)
	d.Set("team_id", tag.Team)
	return nil
}
