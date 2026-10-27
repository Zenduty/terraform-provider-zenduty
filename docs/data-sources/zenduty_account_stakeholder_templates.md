---
page_title: "Zenduty: Account Stakeholder Templates"
subcategory: ""
description: |-
    "`zenduty_account_stakeholder_templates` is a data source that lists the account's stakeholder update templates"
---

# Data Source : zenduty_account_stakeholder_templates

`zenduty_account_stakeholder_templates` lists the account-level stakeholder update templates.

## Example Usage

```hcl
data "zenduty_account_stakeholder_templates" "all" {}

output "template_titles" {
  value = data.zenduty_account_stakeholder_templates.all.templates[*].title
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `templates` - The list of templates, each with:
    * `unique_id` - The unique_id of the template.
    * `title` - The title.
    * `subject` - The subject line.
    * `template` - The body.
