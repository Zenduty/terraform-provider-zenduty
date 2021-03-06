---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Zenduty Incidents"
subcategory: ""
description: |-
    " `zenduty_incidents` is a resource to manage escalation policies in a team "
  
---

# zenduty_incidents (Resource)
`zenduty_incidents` is a resource to manage incidents 

<!-- schema generated by tfplugindocs -->

## Example Usage

```hcl

resource "zenduty_incidents" "incident1"{
    title = ""
    summary = ""
    service = ""
    user= ""
    escalation_policy=""

}

```

## Argument Reference

* `title` (Required) - The title of the incident.
* `summary` (Required) - The summary of the incident.
* `service` (Required) - Unique_ID of the service.
* `user` (Required) - Username of the user.
* `escalation_policy` (Required) - unique_id of the escalation policy.

## DataTypes 
### Required
- **escalation_policy** (String)
- **service** (String)
- **summary** (String)
- **title** (String)
- **user** (String)

### Optional

- **status** (Number)

* `status` (Optional) (Number) -  values are `2` to acknowledge `3` to resolve