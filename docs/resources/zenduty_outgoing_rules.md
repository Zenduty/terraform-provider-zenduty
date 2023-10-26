---
page_title: "Zenduty Outgoing Rules"
subcategory: ""
description: |-
    " `zenduty_outgoing_rules` is a resource to manage outgoing rules in a integration"
---
# zenduty_outgoing_rules (Resource)
`zenduty_outgoing_rules` is a resource to manage outgoing rules in a integration


## Example Usage

```hcl
resource "zenduty_teams" "exampleteam" {
    name = "exmaple team"
}

resource "zenduty_services" "exampleservice" {
    name = "example service"
    team_id = zenduty_teams.exampleteam.id 
    escalation_policy = zenduty_esp.example_esp.id e
}

resource "zenduty_integrations" "exampleintegration" {
    team_id = zenduty_teams.exampleteam.id
    service_id = zenduty_services.exampleservice.id
    application = ""
    name = "exampleintegration"
    summary = "This is the summary for the example integration"
}

```

```hcl 
resource "zenduty_outgoing_rules" "example_outgoingrules" {
    team_id = zenduty_teams.exampleteam.id
    service_id = zenduty_services.exampleservice.id
    integration_id = zenduty_integrations.exampleintegration.id
    rule_json = ""    
}

```

## Argument Reference

* `team_id` (Required) - The unique_id of the team to create the alert rule in.
* `service_id` (Required) - The unique_id of the service to create the alert rule in.
* `integration_id` (Required) - The unique_id of the integration to create the alert rule in.
* `rule_json` (Required)(string) - The rule json of the alert rule.You cannot construct the rule json in terraform as of now.One can construct the rule json in Zenduty's UI.Create an outgoing rule in Zenduty and copy the rule_json from the UI.
* `enabled` (Optional) - boolean value to enable or disabled 


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Outgoing Rule.

## Import

Integrations can be imported using the `team_id`(ie. unique_id of the team), `service_id`(ie. unique_id of the service),`integration_id`(ie. unique_id of the integration) and `outgoing_rule_id`(ie. unique_id of the alert rule).

```hcl
resource "zenduty_outgoing_rules" "demorule" {

}

```

`$ terraform import zenduty_outgoing_rules.demorule team_id/service_id/integration_id/outgoing_rule_id` 

`$ terraform state show zenduty_outgoing_rules.demorule`

`* copy the output data and paste inside zenduty_outgoing_rules.demorule resource block and remove the id attribute`

`$ terraform plan` to verify the import



    

