---
page_title: "Zenduty Incident Custom Field"
subcategory: ""
description: |-
    "`zenduty_incident_custom_field` is a resource to manage account-level incident custom field definitions"
---
# zenduty_incident_custom_field (Resource)
`zenduty_incident_custom_field` defines a custom field that incident types can place on their incident form (see `zenduty_incident_type`). An account can have up to 20 custom fields.

## Example Usage

```hcl
resource "zenduty_incident_custom_field" "region" {
  field_type = 3 # select
  properties = jsonencode({
    type  = "string"
    title = "Region"
    enum  = ["eu-west", "eu-central", "us-east"]
  })
}

resource "zenduty_incident_custom_field" "ticket" {
  field_type = 5 # url
  properties = jsonencode({
    type  = "string"
    title = "Ticket"
  })
}
```

## Argument Reference

* `field_type` (Required)(Number) - `0` string, `2` integer, `3` select, `4` multiselect, `5` url.
* `properties` (Required) - The field's JSON Schema property descriptor, as JSON. `title` is required (unique within the account, at most 256 characters) and `type` must match the field type: `string` for string, select and url, `integer` for integer, `array` for multiselect. A select needs a non-empty `enum`; a multiselect needs `items.enum`. These rules are checked at plan time.

## Attributes Reference

* `id` - The unique_id of the custom field.
* `field_key` - Backend-generated key (`custom_field_<n>`) under which the field's value appears in incident payloads.

## Import

```hcl
resource "zenduty_incident_custom_field" "region" {

}
```

`$ terraform import zenduty_incident_custom_field.region custom_field_id`

`$ terraform state show zenduty_incident_custom_field.region`

`* copy the output data and paste inside zenduty_incident_custom_field.region resource block and remove the id and field_key attributes`

`$ terraform plan` to verify the import
