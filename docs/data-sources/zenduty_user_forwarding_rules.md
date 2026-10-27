---
page_title: "Zenduty: User Forwarding Rules"
subcategory: ""
description: |-
    "`zenduty_user_forwarding_rules` is a data source that lists the forwarding rules involving a user"
---

# Data Source : zenduty_user_forwarding_rules

`zenduty_user_forwarding_rules` lists the forwarding rules involving a user: the rules forwarding their notifications and the rules forwarding other users' notifications to them. Compare `username` with the queried user to tell them apart.

## Example Usage

```hcl
data "zenduty_user_forwarding_rules" "alice" {
  username = data.zenduty_user.alice.users[0].username
}

output "alice_forwards_to" {
  value = [for r in data.zenduty_user_forwarding_rules.alice.forwarding_rules : r.to_user if r.username == data.zenduty_user.alice.users[0].username]
}
```

## Argument Reference

* `username` (Required) - Username of the user.

## Attributes Reference

* `forwarding_rules` - The list of rules, each with:
    * `unique_id` - The unique_id of the rule.
    * `username` - Username of the user whose notifications are forwarded.
    * `to_user` - Username of the receiving user.
    * `from_date`, `to_date` - As `YYYY-MM-DD HH:MM` in `time_zone`.
    * `time_zone` - The IANA time zone of the rule.
