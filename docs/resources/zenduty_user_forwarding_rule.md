---
page_title: "Zenduty User Forwarding Rule"
subcategory: ""
description: |-
    "`zenduty_user_forwarding_rule` forwards a user's notifications to another user for a period"
---
# zenduty_user_forwarding_rule (Resource)
`zenduty_user_forwarding_rule` forwards every notification of a user to another member of the account between two points in time, for example during a vacation. The API has no update for forwarding rules, so every attribute forces a new resource.

## Example Usage

```hcl
data "zenduty_user" "alice" {
  email = "alice@example.com"
}

data "zenduty_user" "bob" {
  email = "bob@example.com"
}

resource "zenduty_user_forwarding_rule" "alice_vacation" {
  username  = data.zenduty_user.alice.users[0].username
  to_user   = data.zenduty_user.bob.users[0].username
  from_date = "2027-08-01 00:00"
  to_date   = "2027-08-15 00:00"
  time_zone = "Europe/Amsterdam"
}
```

## Argument Reference

* `username` (Required, Forces new resource) - Username of the user whose notifications are forwarded.
* `to_user` (Required, Forces new resource) - Username of the member who receives them.
* `from_date` (Required, Forces new resource) - Start as `YYYY-MM-DD HH:MM` in `time_zone`.
* `to_date` (Required, Forces new resource) - End as `YYYY-MM-DD HH:MM` in `time_zone`, after `from_date`.
* `time_zone` (Optional, Forces new resource) - IANA time zone of the dates. Defaults to `UTC`.

## Attributes Reference

* `id` - The unique_id of the forwarding rule.

## Import

```hcl
resource "zenduty_user_forwarding_rule" "alice_vacation" {

}
```

`$ terraform import zenduty_user_forwarding_rule.alice_vacation username/forwarding_rule_id`

`$ terraform state show zenduty_user_forwarding_rule.alice_vacation`

`* copy the output data and paste inside zenduty_user_forwarding_rule.alice_vacation resource block and remove the id attribute`

`$ terraform plan` to verify the import
