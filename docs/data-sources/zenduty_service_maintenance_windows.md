---
page_title: "Zenduty: Service Maintenance Windows"
subcategory: ""
description: |-
    "`zenduty_service_maintenance_windows` is a data source that lists the maintenance windows of a service"
---

# Data Source : zenduty_service_maintenance_windows

`zenduty_service_maintenance_windows` lists the service-level maintenance windows of one service. Team-level windows are listed by `zenduty_maintenance_window`.

## Example Usage

```hcl
data "zenduty_service_maintenance_windows" "api" {
  team_id    = zenduty_teams.ops.id
  service_id = zenduty_services.api.id
}
```

## Argument Reference

* `team_id` (Required) - The unique_id of the team.
* `service_id` (Required) - The unique_id of the service.

## Attributes Reference

* `maintenance_windows` - The list of windows, each with:
    * `unique_id` - The unique_id of the window.
    * `name` - The name.
    * `start_time`, `end_time`, `repeat_until` - As `YYYY-MM-DD HH:MM` in `timezone`.
    * `timezone` - The IANA time zone of the window.
    * `repeat_interval` - The repeat interval, `0` when the window does not repeat.
    * `creation_date` - When the window was created.
