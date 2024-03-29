---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Zenduty: GlobalRoutingRule"
subcategory: ""
description: |-
  Provides a Zenduty GlobalRoutingRule Resource. This allows you to  users to create routing rules.which allows you to route alerts based on the rules
---

# Resource : zenduty_globalrouting_rule

An Global Router allows users to create routing rules. The Router evaluates events sent to this Router against each of its rules, one at a time, and routes the event to a specific Integration based on the first rule that matches.

---

## Example Usage
```hcl

resource "zenduty_globalrouter" "router" {
  name        = "demorouter"
  description = "This is a demo router" 
}

resource "zenduty_globalrouting_rule" "demorules" {
    router_id = zenduty_globalrouter.router.id
    name      = "demorule"
    rule_json = "" 
    actions {
        action_type = 0
        integration = "unique_id of integration"
    }
}

resource "zenduty_globalrouting_rule" "supressrule" {
    router_id = zenduty_globalrouter.router.id
    name      = "supress"
    rule_json = "" 
    actions {
        action_type = 1
    }
}
```


<!-- schema generated by tfplugindocs -->
## Argument Reference

* `name` (Required) - Name of the Routing Rule
* `router_id` - UniqueID of the GlobalRouter
*  `rule_json` (Required)(string) - The rule json of the routing rule.You cannot construct the rule json in terraform as of now.One can construct the rule json in Zenduty's UI.Create an dummy alert rule in Zenduty and copy the rule_json from the UI.
* `actions` (Optional) - The actions to be performed when the rule matches.values are `0` route to integration `1` supress the alert

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the GlobalRouterRule.


## Import

GlobalRouterRule can be imported using the `router_id/rule_id`(ie. UniqueID of the router,rule), e.g.

```hcl
  resource "zenduty_globalrouting_rule" "rule1" {
  
  }

```

`$ terraform import zenduty_globalrouting_rule.rule1 router_id/rule_id` 

`$ terraform state show zenduty_globalrouting_rule.rule1`

`copy the output data and paste inside zenduty_globalrouting_rule.rule1 resource block and remove the id attribute`

`$ terraform plan` to verify the import
