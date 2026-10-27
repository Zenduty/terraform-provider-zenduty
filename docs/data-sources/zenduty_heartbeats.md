---
page_title: "Zenduty: Heartbeats"
subcategory: ""
description: |-
    "`zenduty_heartbeats` is a data source that lists the heartbeats of the account"
---

# Data Source : zenduty_heartbeats

`zenduty_heartbeats` lists the account's heartbeats, optionally filtered to one service.

## Example Usage

```hcl
data "zenduty_heartbeats" "backup" {
  service_id = zenduty_services.backup.id
}

output "expired_heartbeats" {
  value = [for h in data.zenduty_heartbeats.backup.heartbeats : h.name if h.status == "expired"]
}
```

## Argument Reference

* `service_id` (Optional) - Only return heartbeats of this service.

## Attributes Reference

* `heartbeats` - The list of heartbeats, each with:
    * `unique_id`, `name`, `description`
    * `service_id`, `service_name`, `team_id`, `team_name`
    * `interval_seconds`, `is_enabled`, `status` (`disabled`, `pending`, `healthy`, `expired`)
    * `alert_message`, `alert_urgency` (`0` low, `1` high, `-1` when unset), `ping_type`, `email_subject_contains`
    * `ping_url`, `ping_email`, `last_ping_at`, `expires_at`
