---
page_title: "Zenduty Alert Rules"
subcategory: ""
description: |-
    " `zenduty_alertrules` is a resource to manage alert rules in a integration"
---
# zenduty_alertrules (Resource)
`zenduty_alertrules` is a resource to manage alert rules in a integration


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
resource "zenduty_alertrules" "example_alertrules" {
  
    description = "This is the description for the new alert rules"
    team_id = zenduty_teams.exampleteam.id
    service_id = zenduty_services.exampleservice.id
    integration_id = zenduty_integrations.exampleintegration.id
    rule_json = "" 
    #actions
}

```

## Argument Reference

* `team_id` (Required) - The unique_id of the team to create the alert rule in.
* `service_id` (Required) - The unique_id of the service to create the alert rule in.
* `integration_id` (Required) - The unique_id of the integration to create the alert rule in.
* `description` (Required) - The description of the alert rule.
* `rule_json` (Required)(string) - The rule json of the alert rule.You cannot construct the rule json in terraform as of now.One can construct the rule json in Zenduty's UI.Create an dummy alert rule in Zenduty and copy the rule_json from the UI.
* `actions` (Optional) - The rules of the escalation policy. (see [below for nested schema](#nestedblock--actions))


<a id="nestedblock--actions"></a>

## Actions
```hcl
    actions {
        action_type = ""
        value = ""
        #key
    }

```
* `action_type` (Required) (Number):
    * `1` - change the alert type , value should be one of the following: `0` for info, `1` for warning, `2` for error, `3` for critical , `4` for acknowledged , `5` for resolved
    * `2` - add note , value will be the note summary to add
    * `3` - supress alert , value is not required
    * `4` - add escalation policy , value should be the unique_id of the escalation policy
    * `6` - assign user , value should be the username of the user
    * `7` - change urgency  , value should be one of the following: `0` for low, `1` for high
    * `8` - change message , value should be the message to change to 
    * `9` - change summary , value should be the summary to change to
    * `10` - change entry_id , value should be the entity to change to
    * `11` - assign role to user , `key` should be unique_id of the role , value should be the username of the user
    * `12` - add tag, value should be the unique_id of the tag
    * `14` - add sla , value should be the unique_id of the sla
    * `15` - add team priority , value should be the unique_id of the team priority

* `value` (Required)(string) - The value of the action. (not required for `3`)
* `key`  (Optional)(string) - The key of the action. (required for `11`)


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Alert Rule.

## Import

Integrations can be imported using the `team_id`(ie. unique_id of the team), `service_id`(ie. unique_id of the service),`integration_id`(ie. unique_id of the integration) and `alertrule_id`(ie. unique_id of the alert rule).

```hcl
resource "zenduty_alertrules" "rule1" {

}

```

`$ terraform import zenduty_alertrules.rule1 team_id/service_id/integration_id/alertrule_id` 

`$ terraform state show zenduty_alertrules.rule1`

`* copy the output data and paste inside zenduty_alertrules.rule1 resource block and remove the id attribute`

`$ terraform plan` to verify the import



    

