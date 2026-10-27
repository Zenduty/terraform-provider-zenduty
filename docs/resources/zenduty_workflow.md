---
page_title: "Zenduty Workflow"
subcategory: ""
description: |-
    "`zenduty_workflow` is a resource to manage incident workflows"
---
# zenduty_workflow (Resource)
`zenduty_workflow` manages an incident workflow: a graph of a trigger node followed by action, condition and delay nodes, scoped to one or more teams.

~> **API access** The workflow endpoints currently accept browser-session authentication only, not the API token the provider uses. Until the backend grants token access to `/api/workflows/`, every operation of this resource fails with an authentication error.

## Example Usage

```hcl
resource "zenduty_workflow" "note_on_create" {
  name        = "Note on create"
  description = "Adds a note to every new incident"
  is_draft    = false
  is_enabled  = true

  json_schema = jsonencode({
    nodes = [
      {
        id   = "trigger-1"
        type = "trigger"
        data = { name = "Incident Created", action_name = "incident_create", group = "Incident" }
      },
      {
        id   = "action-1"
        type = "action"
        data = {
          name        = "Add Note"
          action_name = "add_note"
          group       = "Incident"
          input       = { note = "Triage checklist: https://runbooks.example.com/triage" }
        }
      },
    ]
    edges = [
      { id = "trigger-1->action-1", source = "trigger-1", target = "action-1" },
    ]
  })

  teams {
    team_id     = zenduty_teams.ops.id
    service_ids = [zenduty_services.api.id]
  }

  teams {
    team_id = zenduty_teams.platform.id # every service of the team
  }
}
```

## Argument Reference

* `name` (Required) - Name of the workflow.
* `description` (Required) - Description of the workflow.
* `is_draft` (Optional)(bool) - Whether the workflow is a draft. Defaults to `true`, like the API; set `false` to publish it.
* `is_enabled` (Optional)(bool) - Whether the workflow runs. Defaults to `true`.
* `json_schema` (Required) - The workflow graph as JSON with `nodes` and `edges`, in the format the web console's workflow editor produces. The graph must start with exactly one `trigger` node, have no cycles and no node with several successors other than condition nodes. The trigger count is checked at plan time, the rest by the API.
* `teams` (Optional) - Teams the workflow runs for; without a block it runs for no team.
    * `team_id` (Required) - The unique_id of the team.
    * `service_ids` (Optional)(Set of String) - Restrict the workflow to these services of the team. Empty means every service.

## Attributes Reference

* `id` - The unique_id of the workflow.

## Import

```hcl
resource "zenduty_workflow" "note_on_create" {

}
```

`$ terraform import zenduty_workflow.note_on_create workflow_id`

`$ terraform state show zenduty_workflow.note_on_create`

`* copy the output data and paste inside zenduty_workflow.note_on_create resource block and remove the id attribute`

`$ terraform plan` to verify the import
