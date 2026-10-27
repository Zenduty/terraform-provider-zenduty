---
page_title: "Zenduty Service Maintenance Window"
subcategory: ""
description: |-
    "`zenduty_service_maintenance_window` is a resource to manage a maintenance window on a single service"
---
# zenduty_service_maintenance_window (Resource)
`zenduty_service_maintenance_window` puts one service into maintenance for a period, optionally repeating. Use `zenduty_maintenance_window` for a team-level window that covers several services at once.

## Example Usage

```hcl
resource "zenduty_service_maintenance_window" "patching" {
  team_id    = zenduty_teams.ops.id
  service_id = zenduty_services.api.id
  name       = "Monthly patching"
  start_time = "2027-01-15 22:00"
  end_time   = "2027-01-16 02:00"
  timezone   = "Europe/Berlin"

  repeat_interval = 1
  repeat_until    = "2027-12-31 02:00"
}
```

## Argument Reference

* `team_id` (Required, Forces new resource) - The unique_id of the team.
* `service_id` (Required, Forces new resource) - The unique_id of the service.
* `name` (Required) - Name of the window.
* `start_time` (Required) - Start as `YYYY-MM-DD HH:MM` in `timezone`.
* `end_time` (Required) - End as `YYYY-MM-DD HH:MM` in `timezone`, after `start_time`.
* `timezone` (Required) - IANA time zone the times are expressed in, e.g. `UTC` or `Europe/Berlin`.
* `repeat_interval` (Optional)(Number) - Repeat interval of the window. `0` (default) means the window does not repeat.
* `repeat_until` (Optional) - Last repetition as `YYYY-MM-DD HH:MM` in `timezone`.

## Attributes Reference

* `id` - The unique_id of the maintenance window.

## Import

```hcl
resource "zenduty_service_maintenance_window" "patching" {

}
```

`$ terraform import zenduty_service_maintenance_window.patching team_id/service_id/maintenance_window_id`

`$ terraform state show zenduty_service_maintenance_window.patching`

`* copy the output data and paste inside zenduty_service_maintenance_window.patching resource block and remove the id attribute`

`$ terraform plan` to verify the import
