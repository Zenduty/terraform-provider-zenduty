---
page_title: "Zenduty: User Custom Role"
subcategory: ""
description: |-
    "`zenduty_user_custom_role` is a data source that reads the custom account role assigned to a user"
---

# Data Source : zenduty_user_custom_role

`zenduty_user_custom_role` reads the custom account role (see `zenduty_account_role`) assigned to a user. The API only serves this for members with the regular user role; account admins and owners cannot hold a custom role.

## Example Usage

```hcl
data "zenduty_user_custom_role" "alice" {
  username = data.zenduty_user.alice.users[0].username
}

output "alice_role" {
  value = data.zenduty_user_custom_role.alice.role_name
}
```

## Argument Reference

* `username` (Required) - Username of the user.

## Attributes Reference

* `account_role` - The unique_id of the assigned custom role, empty when the user has none.
* `role_name` - The name of the assigned role.
