---
page_title: "Zenduty: Workflows"
subcategory: ""
description: |-
    "`zenduty_workflows` is a data source that lists the account's incident workflows"
---

# Data Source : zenduty_workflows

`zenduty_workflows` lists the account's workflows, or reads one workflow in full.

~> **API access** The workflow endpoints currently accept browser-session authentication only, not the API token the provider uses. Until the backend grants token access to `/api/workflows/`, this data source fails with an authentication error.

## Example Usage

```hcl
data "zenduty_workflows" "all" {}

data "zenduty_workflows" "one" {
  workflow_id = zenduty_workflow.note_on_create.id
}
```

## Argument Reference

* `workflow_id` (Optional) - Return only this workflow, with its full `json_schema` and team scope. Without it the list shape is returned: `json_schema` is empty and `teams` carry no `service_ids`.

## Attributes Reference

* `workflows` - The list of workflows, each with:
    * `unique_id`, `name`, `description`
    * `is_draft`, `is_enabled`
    * `json_schema` - The workflow graph as JSON (only with `workflow_id`).
    * `teams` - The team scope, each with `team_id`, `name` and `service_ids`.
