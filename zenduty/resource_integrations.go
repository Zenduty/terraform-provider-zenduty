package zenduty

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationCreate,
		UpdateContext: resourceIntegrationUpdate,
		DeleteContext: resourceIntegrationDelete,
		ReadContext:   resourceIntegrationRead,
		Importer: &schema.ResourceImporter{
			State: resourceIntegrationImporter,
		},
		Schema: map[string]*schema.Schema{
			"application": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"service_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateUUID(),
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
				Default:  true,
			},
			"create_incident_for": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 3),
				Default:      1,
			},
			"default_urgency": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Default:      1,
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
	if v, ok := d.GetOk("is_enabled"); ok {
		newIntegration.Is_Enabled = v.(bool)
	}
	if v, ok := d.GetOk("create_incident_for"); ok {
		newIntegration.Create_Incident_For = v.(int)
	}
	if v, ok := d.GetOk("default_urgency"); ok {
		newIntegration.Default_Urgency = v.(int)
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
	apiclient, _ := m.(*Config).Client()
	id := d.Id()

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
	if v, ok := d.GetOk("is_enabled"); ok {
		newIntegration.Is_Enabled = v.(bool)
	}
	if v, ok := d.GetOk("create_incident_for"); ok {
		newIntegration.Create_Incident_For = v.(int)
	}
	if v, ok := d.GetOk("default_urgency"); ok {
		newIntegration.Default_Urgency = v.(int)
	}

	_, err := apiclient.Integrations.UpdateIntegration(team_id, service_id, id, newIntegration)
	if err != nil {
		return diag.FromErr(err)
	}

	// added integration_key in response output

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
	d.Set("integration_key", integration.Integration_key)
	d.Set("webhook_url", integration.Webhook_url)
	d.Set("is_enabled", integration.Is_Enabled)
	d.Set("create_incident_for", integration.Create_Incident_For)
	d.Set("default_urgency", integration.Default_Urgency)

	return diags
}

func resourceIntegrationImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<service_id>/<integration_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid serviceid (%q)", parts[1])
	} else if !IsValidUUID(parts[2]) {
		return nil, fmt.Errorf("invalid integration_id (%q)", parts[2])
	}
	d.SetId(parts[2])
	d.Set("team_id", parts[0])
	d.Set("service_id", parts[1])

	return []*schema.ResourceData{d}, nil
}
