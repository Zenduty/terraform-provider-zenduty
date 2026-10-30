---
page_title: "Zenduty Account Setting"
subcategory: ""
description: |-
    "`zenduty_account_setting` manages the account's singleton settings"
---
# zenduty_account_setting (Resource)
`zenduty_account_setting` manages the account-wide settings. Every account has exactly one settings object, so the resource does not create anything: `terraform apply` updates the existing settings, and `terraform destroy` only removes the resource from state and leaves the settings as they are.

Only the settings written in the configuration are sent; the others keep their current value and are read into state.

## Example Usage

```hcl
resource "zenduty_account_setting" "main" {
  enable_delete_incident        = false
  enable_team_manager_invite    = true
  enable_regen_key              = true
  enable_first_mtta_calculation = true
}
```

## Argument Reference

* `enable_delete_incident` (Optional)(bool) - Whether incidents can be deleted.
* `enable_team_manager_invite` (Optional)(bool) - Whether team managers can invite users.
* `enable_regen_key` (Optional)(bool) - Whether team managers can regenerate integration keys.
* `enable_first_mtta_calculation` (Optional)(bool) - Whether MTTA is measured to the first acknowledgement.

## Attributes Reference

* `id` - Always `account`.

## Import

```hcl
resource "zenduty_account_setting" "main" {

}
```

`$ terraform import zenduty_account_setting.main account`

`$ terraform state show zenduty_account_setting.main`

`$ terraform plan` to verify the import
