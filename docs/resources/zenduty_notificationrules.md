---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Zenduty: Notification Rules"
subcategory: ""
description: |-
  Provides a Notification Rules Resource. This allows Notification Rules to be created, updated, and deleted.
---

# Resource : zenduty_notification_rules

Provides a Notification Rules Resource. This allows Notification Rules to be created, updated, and deleted.

## Example Usage

```hcl

data "zenduty_user" "user1" {
  email = "demouser@gmail.com"
}

data "zenduty_usercontact" "exampleusercontact" {
  user_id = data.zenduty_user.user1.users[0].username
  contact_type=1
}

```

```hcl 

resource "zenduty_notification_rules" "notification_rules" {
  username = data.zenduty_user.user1.users[0].username
  contact = data.zenduty_usercontact.exampleusercontact.id
  urgency = 1
  delay = 1
}

```

## Argument Reference

* `username` (Required) - The username of the user. 
* `contact` (Required) - The contact ID of the user.
* `urgency` (Required) - The urgency of the notification.
   values are `1` High, `0` Low
* `delay` (Required) - The delay of the notification.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Zenduty Notification Rule.

## Import

Notification Rules can be imported using the `username`(username of user) and `notification_rule_id`(ie. unique_id of the notification rule), e.g.

```hcl
resource "zenduty_notification_rules" "rule1" {

}
```

`$ terraform import zenduty_notification_rules.rule1 username/notification_rule_id` 

`$ terraform state show zenduty_notification_rules.rule1`

`copy the output data and paste inside zenduty_notification_rules.rule1 resource block and remove the id attribute`
`$ terraform plan` to verify the import