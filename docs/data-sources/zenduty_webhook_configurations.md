---
page_title: "Zenduty: Webhook Configurations"
subcategory: ""
description: |-
    "`zenduty_webhook_configurations` is a data source that lists the account's outgoing webhooks"
---

# Data Source : zenduty_webhook_configurations

`zenduty_webhook_configurations` lists the account-level webhook configurations. Secrets are not returned by the API and are not part of this data source.

## Example Usage

```hcl
data "zenduty_webhook_configurations" "all" {}

output "webhook_urls" {
  value = data.zenduty_webhook_configurations.all.webhook_configurations[*].url
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `webhook_configurations` - The list of webhooks, each with:
    * `unique_id` - The unique_id of the webhook configuration.
    * `url` - The endpoint.
    * `event` - The list of subscribed events.
    * `filter` - The JSON filter, if any.
