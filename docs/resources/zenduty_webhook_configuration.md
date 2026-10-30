---
page_title: "Zenduty Webhook Configuration"
subcategory: ""
description: |-
    "`zenduty_webhook_configuration` is a resource to manage account-level outgoing webhooks"
---
# zenduty_webhook_configuration (Resource)
`zenduty_webhook_configuration` sends account events (on-call changes, incident lifecycle) to an HTTPS endpoint.

## Example Usage

```hcl
resource "zenduty_webhook_configuration" "ops_bus" {
  url   = "https://hooks.example.com/zenduty"
  event = ["incident_create", "incident_acknowledge", "incident_resolve", "oncall_change"]

  filter = jsonencode({
    team = "platform"
  })

  auth_secret_key = var.webhook_signing_secret
  api_key         = var.webhook_api_key
}
```

## Argument Reference

* `url` (Required) - HTTPS endpoint that receives the events. The API resolves the host on write and rejects private or unresolvable addresses. URLs are unique.
* `event` (Required)(Set of String) - Events to deliver, one or more of `oncall_change`, `incident_create`, `incident_acknowledge`, `incident_resolve`.
* `filter` (Optional) - JSON filter applied to events before delivery; use `jsonencode()`.
* `auth_secret_key` (Optional, Sensitive) - Secret used to sign deliveries.
* `api_key` (Optional, Sensitive) - API key sent with deliveries.

The API never returns `auth_secret_key` and `api_key` (reads return masked values), so the provider keeps the configured values in state and does not detect out-of-band changes to them. Every update sends the configured `filter`, `auth_secret_key` and `api_key`, so removing one from the configuration clears it on the next apply. After an import the secrets are absent from state and configuration: add them before the next apply that touches the webhook, or that apply clears them.

## Attributes Reference

* `id` - The unique_id of the webhook configuration.

## Import

```hcl
resource "zenduty_webhook_configuration" "ops_bus" {

}
```

`$ terraform import zenduty_webhook_configuration.ops_bus webhook_configuration_id`

`$ terraform state show zenduty_webhook_configuration.ops_bus`

`* copy the output data and paste inside zenduty_webhook_configuration.ops_bus resource block, remove the id attribute and add the secrets`

`$ terraform plan` to verify the import
