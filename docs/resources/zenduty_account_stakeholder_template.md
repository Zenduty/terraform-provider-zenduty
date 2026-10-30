---
page_title: "Zenduty Account Stakeholder Template"
subcategory: ""
description: |-
    "`zenduty_account_stakeholder_template` is a resource to manage account-level stakeholder update templates"
---
# zenduty_account_stakeholder_template (Resource)
`zenduty_account_stakeholder_template` manages a reusable stakeholder update (subject and body) available to every team of the account.

## Example Usage

```hcl
resource "zenduty_account_stakeholder_template" "outage" {
  title    = "Major outage"
  subject  = "[Zenduty] Service disruption update"
  template = "We are investigating a service disruption. Next update in 30 minutes."
}
```

## Argument Reference

* `title` (Required) - Title of the template, unique within the account.
* `subject` (Required) - Subject line of the stakeholder update.
* `template` (Optional) - Body of the stakeholder update.

## Attributes Reference

* `id` - The unique_id of the template.

## Import

```hcl
resource "zenduty_account_stakeholder_template" "outage" {

}
```

`$ terraform import zenduty_account_stakeholder_template.outage stakeholder_template_id`

`$ terraform state show zenduty_account_stakeholder_template.outage`

`* copy the output data and paste inside zenduty_account_stakeholder_template.outage resource block and remove the id attribute`

`$ terraform plan` to verify the import
