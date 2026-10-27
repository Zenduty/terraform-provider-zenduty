---
page_title: "Zenduty: Maintenance Templates"
subcategory: ""
description: |-
    "`zenduty_maintenance_templates` is a data source that lists the maintenance templates of a team"
---

# Data Source : zenduty_maintenance_templates

`zenduty_maintenance_templates` lists a team's maintenance rule templates.

## Example Usage

```hcl
data "zenduty_maintenance_templates" "ops" {
  team_id = zenduty_teams.ops.id
}
```

## Argument Reference

* `team_id` (Required) - The unique_id of the team.

## Attributes Reference

* `maintenance_templates` - The list of templates, each with:
    * `unique_id` - The unique_id of the template.
    * `name` - The name.
    * `is_enabled` - Whether the template is active.
    * `rule_json` - The matching rule as JSON.
