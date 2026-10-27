package zenduty

import (
	"context"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWebhookConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateWebhookConfiguration,
		ReadContext:   resourceReadWebhookConfiguration,
		UpdateContext: resourceUpdateWebhookConfiguration,
		DeleteContext: resourceDeleteWebhookConfiguration,
		Importer: &schema.ResourceImporter{
			State: resourceWebhookConfigurationImporter,
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
				Description:      "HTTPS endpoint that receives the events. The API resolves it on write and rejects private or unresolvable hosts.",
			},
			"event": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(client.WebhookConfigurationEvents, false),
				},
				Description: "Events to send: oncall_change, incident_create, incident_acknowledge, incident_resolve.",
			},
			"filter": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ValidateJSONString(),
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
				Description:      "Optional JSON filter applied to events before delivery.",
			},
			"auth_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Secret used to sign deliveries. The API never returns it, so the configured value is kept as-is in state.",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "API key sent with deliveries. The API never returns it, so the configured value is kept as-is in state.",
			},
		},
	}
}

func expandWebhookConfiguration(d *schema.ResourceData) *client.WebhookConfiguration {
	return &client.WebhookConfiguration{
		URL:           d.Get("url").(string),
		Event:         expandStringSet(d.Get("event")),
		Filter:        d.Get("filter").(string),
		AuthSecretKey: d.Get("auth_secret_key").(string),
		APIKey:        d.Get("api_key").(string),
	}
}

func resourceCreateWebhookConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.WebhookConfiguration.CreateWebhookConfiguration(expandWebhookConfiguration(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadWebhookConfiguration(ctx, d, m)
}

func resourceUpdateWebhookConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if _, err := apiclient.WebhookConfiguration.UpdateWebhookConfiguration(d.Id(), expandWebhookConfiguration(d)); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadWebhookConfiguration(ctx, d, m)
}

// The API masks auth_secret_key and api_key on read, so Read leaves the
// configured secrets in state untouched and only refreshes the other fields.
func resourceReadWebhookConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	config, err := apiclient.WebhookConfiguration.GetWebhookConfiguration(d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	d.Set("url", config.URL)
	d.Set("event", config.Event)
	if config.Filter == "" {
		d.Set("filter", "")
	} else if normalized, normErr := normalizeJSON(config.Filter); normErr == nil {
		d.Set("filter", normalized)
	} else {
		d.Set("filter", config.Filter)
	}
	return nil
}

func resourceDeleteWebhookConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.WebhookConfiguration.DeleteWebhookConfiguration(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceWebhookConfigurationImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if !IsValidUUID(d.Id()) {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <webhook_configuration_id>", d.Id())
	}
	return []*schema.ResourceData{d}, nil
}
