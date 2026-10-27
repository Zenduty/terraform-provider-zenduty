---
page_title: "Zenduty Incident Type"
subcategory: ""
description: |-
    "`zenduty_incident_type` is a resource to manage incident types and the custom fields on their incident form"
---
# zenduty_incident_type (Resource)
`zenduty_incident_type` manages an incident type: its name, whether it is the account default, which built-in form sections it shows and which custom fields (from `zenduty_incident_custom_field`) it carries. An account can have up to 10 active incident types.

## Example Usage

```hcl
resource "zenduty_incident_type" "outage" {
  name       = "Outage"
  is_default = true

  form_toggles = {
    impacted_services = true
    responders        = true
  }

  custom_fields {
    custom_field  = zenduty_incident_custom_field.region.id
    required      = true
    default_value = jsonencode("eu-west")
  }

  custom_fields {
    custom_field = zenduty_incident_custom_field.ticket.id
  }
}
```

## Argument Reference

* `name` (Required) - Name of the type, unique within the account.
* `is_default` (Optional)(bool) - Whether this is the account's default incident type. Setting it clears the flag on the previous default. Defaults to `false`.
* `is_active` (Optional)(bool) - Whether the type can be picked for new incidents. Defaults to `true`. The API always creates a type active; with `false` the provider deactivates it right after creation.
* `form_toggles` (Optional)(Map of bool) - Built-in incident form sections shown for this type, e.g. `impacted_services`, `responders`.
* `custom_fields` (Optional) - Custom fields on the type's incident form, one block per field. The block order is the display order.
    * `custom_field` (Required) - The unique_id of the `zenduty_incident_custom_field`.
    * `required` (Optional)(bool) - Whether the field must be filled. Defaults to `false`.
    * `default_value` (Optional) - Default value as JSON: `jsonencode("eu-west")`, `jsonencode(3)` or `jsonencode(["a", "b"])` for a multiselect. Validated against the field's type. Unset inherits the default from the field definition's `properties`, if it has one.

## Deletion

Deleting an incident type deactivates it (`is_active = false`) instead of removing it, so that existing incidents keep their type. The name stays reserved on the account; re-creating a type with the same name fails until the deactivated one is renamed in the web console.

## Attributes Reference

* `id` - The unique_id of the incident type.

## Import

```hcl
resource "zenduty_incident_type" "outage" {

}
```

`$ terraform import zenduty_incident_type.outage incident_type_id`

`$ terraform state show zenduty_incident_type.outage`

`* copy the output data and paste inside zenduty_incident_type.outage resource block and remove the id attribute`

`$ terraform plan` to verify the import
