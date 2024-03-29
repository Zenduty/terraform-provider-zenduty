---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Zenduty Integration"
subcategory: ""
description: |-
    "`zenduty_integrations` is a resource to manage integrations of a service."
  
---

# Resource : zenduty_integrations

`zenduty_integrations` is a resource to manage integrations of a service.

## Example Usage

```hcl
resource "zenduty_teams" "exampleteam" {
  name = "exmaple team"
}

resource "zenduty_esp" "example_esp" {
  name = "example esp"
  team_id = zenduty_teams.exampleteam.id
  description = "this is an example esp"
}

resource "zenduty_services" "exampleservice" {
  name = "example service"
  team_id = zenduty_teams.exampleteam.id 
  escalation_policy = zenduty_esp.example_esp.id e
}
```


```hcl
resource "zenduty_integrations" "exampleintegration" {
    team_id = zenduty_teams.exampleteam.id
    service_id = zenduty_services.exampleservice.id
    application = ""
    name = "exampleintegration"
    summary = "This is the summary for the example integration"
}

```
## Argument Reference
* `name` (Required) - The name of the integration.
* `team_id` (Required)- The Unique_ID of the team to which the integration belongs.
* `service_id` (Required)- The Unique_ID of the service to which the integration belongs.
* `summary` (Required)- The summary for the integration.
* `application` (Required)- The application id you want to be integrated.
* `is_enabled` (Optional)(Boolean) - Whether the integration is enabled or not.
* `create_incident_for` (Optional)(Int) - Type of Alerts to create an incident.`0`:Don't create incidents,`1`:critical alerts, `2`:critical and error alerts,
`3`:critical, error and warning alerts
* `default_urgency` (Optional)(Int) - The default urgency of the incident. values are `1` for high `0` for low.

* To get application id, vist https://www.zenduty.com/api/account/applications/ and get unique_id of the application.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Integration.
* `integration_key` - The integration key of the Integration.
* `is_enabled` - Whether the Integration is enabled or not.
* `webhook_url` - The webhook url of the Integration.

## Import

Integrations can be imported using the `team_id`(ie. unique_id of the team) and `service_id`(ie. unique_id of the service)`integration_id`(ie. unique_id of the integration).

```hcl
  resource "zenduty_integrations" "integration1" {
  
  }
  
```

`$ terraform import zenduty_integrations.integration1 team_id/service_id/integration_id` 

`$ terraform state show zenduty_integrations.integration1`

`* copy the output data and paste inside zenduty_integrations.integration1 resource block and remove the id,integration_key,webhook_url attributes`

`$ terraform plan` to verify the import




<!-- schema generated by tfplugindocs -->
## DataTypes 

 Required fields:

- **application** (String)
- **name** (String)
- **service_id** (String)
- **summary** (String)
- **team_id** (String)


