package zenduty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrderRead,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"items": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"team": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rank": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	var diags diag.Diagnostics

	team_id := d.Get("team_id").(string)

	roles, err := apiclient.Roles.GetRoles(team_id)
	if err != nil {
		return diag.FromErr(err)
	}

	items := make([]map[string]interface{}, len(roles))
	for i, role := range roles {
		item := make(map[string]interface{})
		item["team"] = role.Team
		item["unique_id"] = role.Unique_Id
		item["title"] = role.Title
		item["creation_date"] = role.Creation_Date
		item["description"] = role.Description
		item["rank"] = role.Rank
		items[i] = item
	}

	if err := d.Set("items", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(team_id)

	return diags

}
