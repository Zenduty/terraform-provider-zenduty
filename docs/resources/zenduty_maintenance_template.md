---
page_title: "Zenduty Maintenance Template"
subcategory: ""
description: |-
    "`zenduty_maintenance_template` is a resource to manage team-level maintenance rule templates"
---
# zenduty_maintenance_template (Resource)
`zenduty_maintenance_template` defines an alert-matching rule for a team. Alerts matching the rule are suppressed while a maintenance window that uses the template is active.

## Example Usage

```hcl
resource "zenduty_maintenance_template" "deploys" {
  team_id    = zenduty_teams.ops.id
  name       = "Deployment noise"
  is_enabled = true

  # same react-query-builder format as zenduty_alertrules rule_json
  rule_json = jsonencode({
    type = "group"
    id   = "9a99988a-0123-4456-b89a-b1607f326fd8"
    children1 = {
      "898aa9ba-0123-4456-b89a-b16a429f5d94" = {
        type = "rule"
        id   = "898aa9ba-0123-4456-b89a-b16a429f5d94"
        properties = {
          field           = "alert.incident_urgency"
          operator        = "select_equals"
          value           = ["low"]
          valueSrc        = ["value"]
          operatorOptions = null
          valueType       = ["select"]
        }
        path = ["9a99988a-0123-4456-b89a-b1607f326fd8", "898aa9ba-0123-4456-b89a-b16a429f5d94"]
      }
    }
    properties = { conjunction = "AND", not = false }
    path       = ["9a99988a-0123-4456-b89a-b1607f326fd8"]
  })
}
```

## Argument Reference

* `team_id` (Required, Forces new resource) - The unique_id of the team.
* `name` (Required) - Name of the template, 1 to 50 characters.
* `is_enabled` (Optional)(bool) - Whether the template is active. Defaults to `true`.
* `rule_json` (Required) - The matching rule as JSON, in the same format as `zenduty_alertrules` `rule_json`. Use `jsonencode()`; key order and whitespace differences are ignored.

## Attributes Reference

* `id` - The unique_id of the maintenance template.

## Import

```hcl
resource "zenduty_maintenance_template" "deploys" {

}
```

`$ terraform import zenduty_maintenance_template.deploys team_id/maintenance_template_id`

`$ terraform state show zenduty_maintenance_template.deploys`

`* copy the output data and paste inside zenduty_maintenance_template.deploys resource block and remove the id attribute`

`$ terraform plan` to verify the import
