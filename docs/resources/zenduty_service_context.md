---
page_title: "Zenduty Service Context"
subcategory: ""
description: |-
    "`zenduty_service_context` links a context integration to a service"
---
# zenduty_service_context (Resource)
`zenduty_service_context` attaches a context integration (one that enriches incidents rather than creating them) to a service. The API has no update for this object, so every attribute forces a new resource.

## Example Usage

```hcl
resource "zenduty_service_context" "runbooks" {
  team_id     = zenduty_teams.ops.id
  service_id  = zenduty_services.api.id
  integration = zenduty_integrations.runbook_links.id
}
```

## Argument Reference

* `team_id` (Required, Forces new resource) - The unique_id of the team.
* `service_id` (Required, Forces new resource) - The unique_id of the service.
* `integration` (Required, Forces new resource) - The unique_id of the context integration. A service can link each integration once.

## Attributes Reference

* `id` - The unique_id of the service context.

## Import

```hcl
resource "zenduty_service_context" "runbooks" {

}
```

`$ terraform import zenduty_service_context.runbooks team_id/service_id/service_context_id`

`$ terraform state show zenduty_service_context.runbooks`

`* copy the output data and paste inside zenduty_service_context.runbooks resource block and remove the id attribute`

`$ terraform plan` to verify the import
