---
page_title: "Zenduty Alert Grouping Policy"
subcategory: ""
description: |-
    "`zenduty_alert_grouping_policy` is a resource to manage content-based alert collation for a service"
---
# zenduty_alert_grouping_policy (Resource)
`zenduty_alert_grouping_policy` configures content-based alert collation (noise reduction) for a service: alerts are grouped into one incident when they match on the fields you list here, instead of purely on arrival time.

Set the service's `collation` attribute to `3` (with a `collation_time`) to enable content-based collation; `collation` mode `1` is time-based collation, a separate mode that does not use this resource.

## Example Usage

```hcl
resource "zenduty_teams" "exampleteam" {
  name = "example team"
}

resource "zenduty_services" "exampleservice" {
  name              = "example service"
  team_id           = zenduty_teams.exampleteam.id
  escalation_policy = zenduty_esp.example_esp.id
  collation         = 3
  collation_time    = 10
}

resource "zenduty_alert_grouping_policy" "example_policy" {
  team_id       = zenduty_teams.exampleteam.id
  service_id    = zenduty_services.exampleservice.id
  match_mode    = 2
  static_fields = ["incident_title", "integration"]
  custom_fields = ["message"]
  time_window   = 5
}
```

## Argument Reference

* `team_id` (Required, Forces new resource) - The unique_id of the team the service belongs to.
* `service_id` (Required, Forces new resource) - The unique_id of the service this policy applies to.
* `match_mode` (Optional)(Number) - `1` requires only one of the listed fields (`static_fields` and `custom_fields` combined) to match, `2` requires every listed field to match. Defaults to `2`.
* `static_fields` (Optional)(List of String) - Built-in alert attributes to compare. One or more of `incident_title`, `incident_summary`, `integration`.
* `custom_fields` (Optional)(List of String) - JSON keys from the integration's alert payload to compare (e.g. `host_id`).
* `time_window` (Required)(Number) - Minutes to wait for matching alerts before closing correlation eligibility, `1` to `20`.
* `is_active` (Optional)(bool) - Whether the policy is enforced. Defaults to `true`.

At least one of `static_fields` or `custom_fields` must be set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alert grouping policy.

## Import

Alert grouping policies can be imported using the `team_id` (unique_id of the team), `service_id` (unique_id of the service), and `alert_grouping_policy_id` (unique_id of the policy).

```hcl
resource "zenduty_alert_grouping_policy" "policy1" {

}
```

`$ terraform import zenduty_alert_grouping_policy.policy1 team_id/service_id/alert_grouping_policy_id`

`$ terraform state show zenduty_alert_grouping_policy.policy1`

`* copy the output data and paste inside zenduty_alert_grouping_policy.policy1 resource block and remove the id attribute`

`$ terraform plan` to verify the import
