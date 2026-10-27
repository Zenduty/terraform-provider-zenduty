---
page_title: "Zenduty Heartbeat"
subcategory: ""
description: |-
    "`zenduty_heartbeat` is a resource to manage a heartbeat monitor that pages a service when it stops receiving pings"
---
# zenduty_heartbeat (Resource)
`zenduty_heartbeat` creates an account-level heartbeat. An external job pings the heartbeat's `ping_url` (or `ping_email`); when no ping arrives within `interval_seconds`, an incident is raised on the service.

## Example Usage

```hcl
resource "zenduty_services" "backup" {
  name              = "nightly backup"
  team_id           = zenduty_teams.ops.id
  escalation_policy = zenduty_esp.ops.id
}

resource "zenduty_heartbeat" "backup" {
  service_id       = zenduty_services.backup.id
  name             = "nightly-backup"
  description      = "Pinged at the end of the backup job"
  interval_seconds = 86400
  alert_message    = "The nightly backup did not report in"
  alert_urgency    = 1
}

output "backup_ping_url" {
  value = zenduty_heartbeat.backup.ping_url
}
```

## Argument Reference

* `service_id` (Required, Forces new resource) - The unique_id of the service paged when the heartbeat expires. The service must have an escalation policy.
* `name` (Required) - Name of the heartbeat, unique within the service.
* `description` (Optional) - Free text.
* `interval_seconds` (Required)(Number) - Seconds without a ping before the heartbeat expires, `60` to `7776000` (90 days). The timer starts at creation, so a heartbeat that is never pinged fires once.
* `is_enabled` (Optional)(bool) - Whether the heartbeat is monitored. Defaults to `true`. Re-enabling resets it to `pending` until the next ping.
* `alert_message` (Optional) - Message of the incident raised when the heartbeat expires.
* `alert_urgency` (Optional)(Number) - Urgency of the incident: `0` low, `1` high. Unset leaves the API default.
* `ping_type` (Optional)(Number) - `0` for webhook pings to `ping_url` (default), `1` for email pings to `ping_email`.
* `email_subject_contains` (Optional) - For email pings, only emails whose subject contains this text count as a ping.

## Attributes Reference

* `id` - The unique_id of the heartbeat.
* `team_id` - The unique_id of the service's team.
* `ping_url` - URL to request to register a ping.
* `ping_email` - Address to email to register a ping.
* `status` - `disabled`, `pending` (no ping yet), `healthy` or `expired`.

## Import

```hcl
resource "zenduty_heartbeat" "backup" {

}
```

`$ terraform import zenduty_heartbeat.backup heartbeat_id`

`$ terraform state show zenduty_heartbeat.backup`

`* copy the output data and paste inside zenduty_heartbeat.backup resource block and remove the id and computed attributes`

`$ terraform plan` to verify the import
