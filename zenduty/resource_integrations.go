package zenduty

import (
	"context"
	"errors"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationCreate,
		UpdateContext: resourceIntegrationUpdate,
		DeleteContext: resourceIntegrationDelete,
		ReadContext:   resourceIntegrationRead,
		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"summary": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	newIntegration := &client.IntegrationCreate{}
	var diags diag.Diagnostics
	team_id := d.Get("team_id").(string)
	service_id := d.Get("service_id").(string)
	summary := d.Get("summary").(string)
	if summary != "" {
		newIntegration.Summary = summary
	}
	if v, ok := d.GetOk("name"); ok {
		newIntegration.Name = v.(string)
	}
	if v, ok := d.GetOk("application"); ok {
		newIntegration.Application = v.(string)
	}

	integration, err := apiclient.Integrations.CreateIntegration(team_id, service_id, newIntegration)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(integration.Unique_Id)
	// added integration_key in response output
	d.Set("integration_key", integration.Integration_key)
	d.Set("is_enabled", integration.Is_Enabled)
	d.Set("webhook_url", integration.Webhook_url)
	return diags
}

func resourceIntegrationUpdate(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	return diags
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	id := d.Id()
	team_id := d.Get("team_id").(string)
	service_id := d.Get("service_id").(string)
	var diags diag.Diagnostics
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	if service_id == "" {
		return diag.FromErr(errors.New("service_id is required"))
	}
	err := apiclient.Integrations.DeleteIntegration(team_id, service_id, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	id := d.Id()
	team_id := d.Get("team_id").(string)
	service_id := d.Get("service_id").(string)
	apiclient, _ := m.(*Config).Client()

	integration, err := apiclient.Integrations.GetIntegrationByID(team_id, service_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", integration.Name)
	d.Set("application", integration.Application)
	d.Set("summary", integration.Summary)

	return diags
}
