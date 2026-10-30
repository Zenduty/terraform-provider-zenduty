---
page_title: "Zenduty: Incident Custom Fields"
subcategory: ""
description: |-
    "`zenduty_incident_custom_fields` is a data source that lists the account's incident custom field definitions"
---

# Data Source : zenduty_incident_custom_fields

`zenduty_incident_custom_fields` lists the custom fields defined on the account.

## Example Usage

```hcl
data "zenduty_incident_custom_fields" "all" {}

output "region_field_key" {
  value = one([for f in data.zenduty_incident_custom_fields.all.custom_fields : f.field_key if f.title == "Region"])
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `custom_fields` - The list of fields, each with:
    * `unique_id` - The unique_id of the field.
    * `field_key` - The backend-generated key (`custom_field_<n>`).
    * `field_type` - `0` string, `2` integer, `3` select, `4` multiselect, `5` url.
    * `title` - The field's title.
    * `properties` - The JSON Schema property descriptor as a JSON string.
