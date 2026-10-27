---
page_title: "Zenduty: Account Setting"
subcategory: ""
description: |-
    "`zenduty_account_setting` is a data source that reads the account's settings"
---

# Data Source : zenduty_account_setting

`zenduty_account_setting` reads the account-wide settings.

## Example Usage

```hcl
data "zenduty_account_setting" "main" {}

output "incidents_deletable" {
  value = data.zenduty_account_setting.main.enable_delete_incident
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `enable_delete_incident` - Whether incidents can be deleted.
* `enable_team_manager_invite` - Whether team managers can invite users.
* `enable_regen_key` - Whether team managers can regenerate integration keys.
* `enable_first_mtta_calculation` - Whether MTTA is measured to the first acknowledgement.
