package zenduty

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamReads,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"teams": {
				Type:        schema.TypeList,
				Description: "List of teams",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"members": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unique_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"team": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"joining_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"roles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unique_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"team": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_date": {
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
				},
			},
		},
	}
}

func dataSourceTeamReads(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	if team_id != "" {
		var diags diag.Diagnostics
		team, err := apiclient.Teams.GetTeamById(team_id)
		if err != nil {
			return diag.FromErr(err)
		}
		items := make([]map[string]interface{}, 1)
		item := make(map[string]interface{})
		item["unique_id"] = team.Unique_Id
		item["name"] = team.Name
		item["owner"] = team.Owner
		item["account"] = team.Account
		item["creation_date"] = team.Creation_Date
		roles := make([]map[string]interface{}, len(team.Roles))
		for j, role := range team.Roles {
			roles[j] = map[string]interface{}{
				"unique_id":     role.Unique_Id,
				"title":         role.Title,
				"description":   role.Description,
				"creation_date": role.Creation_Date,
				"rank":          role.Rank,
			}
		}
		item["roles"] = roles
		item["roles"] = roles

		members := make([]map[string]interface{}, len(team.Members))
		for j, member := range team.Members {
			members[j] = map[string]interface{}{
				"unique_id":    member.Unique_Id,
				"team":         member.Team,
				"joining_date": member.Joining_Date,
				"role":         member.Role,
				"user": map[string]interface{}{
					"username":   member.User.Username,
					"first_name": member.User.First_Name,
					"last_name":  member.User.Last_Name,
					"email":      member.User.Email,
				},
			}

		}
		item["members"] = members
		items[0] = item
		if err := d.Set("teams", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())

		return diags

	} else {

		var diags diag.Diagnostics

		teams, err := apiclient.Teams.GetTeams()
		if err != nil {
			return diag.FromErr(err)
		}
		items := make([]map[string]interface{}, len(teams))
		for i, team := range teams {
			item := make(map[string]interface{})
			item["unique_id"] = team.Unique_Id
			item["name"] = team.Name
			item["owner"] = team.Owner
			roles := make([]map[string]interface{}, len(team.Roles))
			for j, role := range team.Roles {
				roles[j] = map[string]interface{}{
					"unique_id":     role.Unique_Id,
					"title":         role.Title,
					"description":   role.Description,
					"creation_date": role.Creation_Date,
					"rank":          role.Rank,
				}
			}
			item["roles"] = roles

			members := make([]map[string]interface{}, len(team.Members))
			for j, member := range team.Members {
				members[j] = map[string]interface{}{
					"unique_id":    member.Unique_Id,
					"team":         member.Team,
					"joining_date": member.Joining_Date,
					"role":         member.Role,
					"user": map[string]interface{}{
						"username":   member.User.Username,
						"first_name": member.User.First_Name,
						"last_name":  member.User.Last_Name,
						"email":      member.User.Email,
					},
				}

			}
			item["members"] = members
			items[i] = item

		}
		if err := d.Set("teams", items); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(time.Now().String())

		return diags
	}

}
