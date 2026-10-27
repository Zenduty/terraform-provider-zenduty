package zenduty

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateWorkflow,
		ReadContext:   resourceReadWorkflow,
		UpdateContext: resourceUpdateWorkflow,
		DeleteContext: resourceDeleteWorkflow,
		CustomizeDiff: validateWorkflowDiff,
		Importer: &schema.ResourceImporter{
			State: resourceWorkflowImporter,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
			},
			"description": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
			},
			"is_draft": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"json_schema": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateJSONObjectString(),
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
				Description:      "The workflow graph as JSON: {nodes = [...], edges = [...]}, starting with exactly one trigger node, as the web console's workflow editor produces it.",
			},
			"teams": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Teams the workflow runs for. Without a block the workflow runs for no team.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"team_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ValidateUUID(),
						},
						"service_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Restrict the workflow to these services of the team. Empty means every service.",
						},
					},
				},
			},
		},
	}
}

// validateWorkflowSchema checks the shape the API validates first: a JSON
// object with a nodes list that contains exactly one trigger node.
func validateWorkflowSchema(raw string) error {
	var graph struct {
		Nodes []struct {
			Type string `json:"type"`
		} `json:"nodes"`
		Edges []interface{} `json:"edges"`
	}
	if err := json.Unmarshal([]byte(raw), &graph); err != nil {
		return errors.New("json_schema must be a JSON object with nodes and edges lists")
	}
	if graph.Nodes == nil {
		return errors.New("json_schema needs a nodes list")
	}
	triggers := 0
	for _, node := range graph.Nodes {
		if node.Type == "trigger" {
			triggers++
		}
	}
	if triggers != 1 {
		return fmt.Errorf("json_schema must contain exactly one node of type \"trigger\", found %d", triggers)
	}
	return nil
}

func validateWorkflowDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	if !diff.NewValueKnown("json_schema") {
		return nil
	}
	return validateWorkflowSchema(diff.Get("json_schema").(string))
}

func expandWorkflowInput(d *schema.ResourceData) *client.WorkflowInput {
	input := &client.WorkflowInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IsDraft:     d.Get("is_draft").(bool),
		IsEnabled:   d.Get("is_enabled").(bool),
		JSONSchema:  d.Get("json_schema").(string),
		Teams:       []client.WorkflowTeamInput{},
	}
	for _, raw := range d.Get("teams").(*schema.Set).List() {
		block, _ := raw.(map[string]interface{})
		input.Teams = append(input.Teams, client.WorkflowTeamInput{
			UniqueID: block["team_id"].(string),
			Services: expandStringSet(block["service_ids"]),
		})
	}
	return input
}

func flattenWorkflowTeams(teams []client.WorkflowTeam) []interface{} {
	out := make([]interface{}, 0, len(teams))
	for _, team := range teams {
		services := make([]interface{}, 0, len(team.Services))
		for _, service := range team.Services {
			services = append(services, service.UniqueID)
		}
		out = append(out, map[string]interface{}{
			"team_id":     team.UniqueID,
			"service_ids": services,
		})
	}
	return out
}

func resourceCreateWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.Workflow.CreateWorkflow(expandWorkflowInput(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadWorkflow(ctx, d, m)
}

func resourceUpdateWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if _, err := apiclient.Workflow.UpdateWorkflow(d.Id(), expandWorkflowInput(d)); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadWorkflow(ctx, d, m)
}

func resourceReadWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	workflow, err := apiclient.Workflow.GetWorkflow(d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	d.Set("name", workflow.Name)
	d.Set("description", workflow.Description)
	d.Set("is_draft", workflow.IsDraft)
	d.Set("is_enabled", workflow.IsEnabled)
	if normalized, normErr := normalizeJSON(workflow.JSONSchema); normErr == nil {
		d.Set("json_schema", normalized)
	} else {
		d.Set("json_schema", workflow.JSONSchema)
	}
	if err := d.Set("teams", flattenWorkflowTeams(workflow.Teams)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDeleteWorkflow(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.Workflow.DeleteWorkflow(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceWorkflowImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if !IsValidUUID(d.Id()) {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <workflow_id>", d.Id())
	}
	return []*schema.ResourceData{d}, nil
}
