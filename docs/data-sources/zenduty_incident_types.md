---
page_title: "Zenduty: Incident Types"
subcategory: ""
description: |-
    "`zenduty_incident_types` is a data source that lists the account's active incident types"
---

# Data Source : zenduty_incident_types

`zenduty_incident_types` lists the active incident types of the account with their form sections and custom fields.

## Example Usage

```hcl
data "zenduty_incident_types" "all" {}

output "default_incident_type" {
  value = one([for t in data.zenduty_incident_types.all.incident_types : t.unique_id if t.is_default])
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `incident_types` - The list of types, each with:
    * `unique_id` - The unique_id of the type.
    * `name` - The name.
    * `is_default` - Whether it is the account default.
    * `is_active` - Always `true` here; deactivated types are not listed.
    * `form_toggles` - Map of built-in form sections to whether they are shown.
    * `custom_fields` - The mapped custom fields in display order, each with `custom_field` (the field's unique_id), `field_key`, `position`, `required` and `default_value` (the effective default as JSON, the field definition's own when the type sets none).
