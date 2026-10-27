---
page_title: "Zenduty: Service Contexts"
subcategory: ""
description: |-
    "`zenduty_service_contexts` is a data source that lists the context integrations linked to a service"
---

# Data Source : zenduty_service_contexts

`zenduty_service_contexts` lists the context integrations attached to a service.

## Example Usage

```hcl
data "zenduty_service_contexts" "api" {
  team_id    = zenduty_teams.ops.id
  service_id = zenduty_services.api.id
}
```

## Argument Reference

* `team_id` (Required) - The unique_id of the team.
* `service_id` (Required) - The unique_id of the service.

## Attributes Reference

* `service_contexts` - The list of links, each with:
    * `unique_id` - The unique_id of the service context.
    * `integration` - The unique_id of the context integration.
    * `created_by` - Username of the creator.
    * `creation_date` - When the link was created.
