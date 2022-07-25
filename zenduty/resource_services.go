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

func resourceServices() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateServices,
		UpdateContext: resourceUpdateServices,
		DeleteContext: resourceDeleteServices,
		ReadContext:   resourceReadServices,
		Importer: &schema.ResourceImporter{
			State: resourceServiceImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"escalation_policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateUUID(),
			},
			"team_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateUUID(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"summary": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"collation": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"collation_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sla": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateUUID(),
			},
			"task_template": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateUUID(),
			},
			"team_priority": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateUUID(),
			},
		},
	}
}

func CreateServices(Ctx context.Context, d *schema.ResourceData, m interface{}) (*client.Services, error) {
	new_service := &client.Services{}

	if v, ok := d.GetOk("name"); ok {
		new_service.Name = v.(string)
	}
	if v, ok := d.GetOk("escalation_policy"); ok {
		new_service.Escalation_Policy = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		new_service.Description = v.(string)
	}
	if v, ok := d.GetOk("summary"); ok {
		new_service.Summary = v.(string)
	}
	if v, ok := d.GetOk("collation"); ok {
		new_service.Collation = v.(int)
	}
	if v, ok := d.GetOk("collation_time"); ok {

		new_service.Collation_Time = v.(int)

	}
	if v, ok := d.GetOk("sla"); ok {
		new_service.Sla = v.(string)
	}
	if v, ok := d.GetOk("task_template"); ok {
		new_service.Task_Template = v.(string)
	}
	if v, ok := d.GetOk("team_priority"); ok {
		new_service.Team_Priority = v.(string)
	}
	if new_service.Collation == 1 && new_service.Collation_Time == 0 {
		return nil, fmt.Errorf("collation_time is required when collation is enabled")

	}
	return new_service, nil
}

func resourceCreateServices(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}

	new_service, serviceErr := CreateServices(Ctx, d, m)
	if serviceErr != nil {
		return diag.FromErr(serviceErr)
	}

	var diags diag.Diagnostics
	service, err := apiclient.Services.CreateService(team_id, new_service)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(service.Unique_Id)
	d.Set("team_id", team_id)

	return diags
}

func resourceUpdateServices(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	new_service, serviceErr := CreateServices(Ctx, d, m)
	if serviceErr != nil {
		return diag.FromErr(serviceErr)

	}

	_, err := apiclient.Services.UpdateService(team_id, id, new_service)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceReadServices(Ctx, d, m)
}

func resourceDeleteServices(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	err := apiclient.Services.DeleteService(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceReadServices(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	id := d.Id()
	if emptyString(team_id) {
		return diag.FromErr(errors.New("team_id is required"))
	}
	var diags diag.Diagnostics
	service, err := apiclient.Services.GetServicesById(team_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", service.Name)
	d.Set("escalation_policy", service.Escalation_Policy)
	d.Set("description", service.Description)
	d.Set("summary", service.Summary)
	d.Set("collation", service.Collation)
	d.Set("collation_time", service.Collation_Time)
	d.Set("sla", service.Sla)
	d.Set("task_template", service.Task_Template)
	d.Set("team_priority", service.Team_Priority)
	d.Set("team_id", team_id)

	return diags
}

func resourceServiceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<service_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid service_id (%q)", parts[1])
	}
	d.Set("team_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
