package zenduty

import (
	"context"
	"errors"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateTags,
		UpdateContext: resourceUpdateTags,
		DeleteContext: resourceDeleteTags,
		ReadContext:   resourceReadTag,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func validateTags(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.Tag, diag.Diagnostics) {
	name := d.Get("name").(string)
	color := d.Get("color").(string)
	team := d.Get("team_id").(string)
	new_tag := &client.Tag{}
	if !IsValidUUID(team) {
		return nil, diag.FromErr(errors.New("team_id must be a valid UUID"))
	}
	if color != "" && !checkList(color, []string{"magenta", "red", "volcano", "orange", "gold", "lime", "green", "cyan", "blue", "geekblue", "purple"}) {
		return nil, diag.FromErr(errors.New("color must be one of the following: magenta, red, volcano, orange, gold, lime, green, cyan, blue, geekblue, purple"))
	}

	new_tag.Name = name
	new_tag.Color = color

	return new_tag, nil
}

func resourceCreateTags(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := d.Get("team_id").(string)
	apiclient, _ := m.(*Config).Client()
	newtag, validationerr := validateTags(ctx, d, m)
	if validationerr != nil {
		return validationerr
	}
	tag, err := apiclient.Tags.CreateTag(team, newtag)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tag.Unique_Id)
	return nil
}

func resourceUpdateTags(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := d.Get("team_id").(string)
	apiclient, _ := m.(*Config).Client()
	newtag, validationerr := validateTags(ctx, d, m)
	if validationerr != nil {
		return validationerr
	}
	tag, err := apiclient.Tags.UpdateTag(team, d.Id(), newtag)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tag.Unique_Id)
	return nil
}

func resourceDeleteTags(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := d.Get("team_id").(string)
	apiclient, _ := m.(*Config).Client()

	err := apiclient.Tags.DeleteTag(team, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceReadTag(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	team := d.Get("team_id").(string)
	tag, err := apiclient.Tags.GetTagId(team, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tag.Unique_Id)
	d.Set("name", tag.Name)
	d.Set("color", tag.Color)
	d.Set("team_id", tag.Team)
	return nil
}
