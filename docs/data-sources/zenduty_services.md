---
# generated by https://github.com/hashicorp/terraform-plugin-docs
layout: "zenduty"
page_title: "zenduty_services Data Source - terraform-provider-zenduty"
subcategory: ""
description: |- 
  "`zenduty_services` data source allows you to query services in the account."
---

# zenduty_services (Data Source)

```hcl 

data "zenduty_services" "exampleservices" {
    team_id = ""
}

```

`or` 

```hcl 

data "zenduty_services" "exampleservice" {
    team_id = ""
    service_id = ""
}

```

```hcl
output "services" {
  value = data.zenduty_services.exampleservices.services
}

```


## Argument Reference
* `team_id`(Required) - The UniqueID of the team to query.
* `service_id`(Optional) - The UniqueID of the service to query.along with team id

## Attributes Reference

The following attributes are exported as list of maps:

* `name` - The name of the service.
* `description` - The description of the service.
* `summary` - The summary of the service
* `unique_id` - The unique_id of the service.
* `team_id` - The unique_id of the team.
* `escalation_policy` - unique_id escalation policy of the service.
* `collation` - The collation of the service.which specifes collation is enabled or not.
* `collation_time` - The collation time of the service in seconds
* `sla` - Uniqueid of sla associated with the service
* `task_template` - Uniqueid of task template associated with the service
* `team_priority` - Uniqueid of team priority associated with the service
* `under_maintenance` - The under_maintenance of the service.which specifes service is under maintenance or not.





<!-- schema generated by tfplugindocs -->

### Required

- **team_id** (String)

### Optional

- **service_id** (String) The ID of this resource.

