package zenduty

import (
	"context"
	"errors"

	"github.com/Zenduty/zenduty-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServices() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateServices,
		UpdateContext: resourceUpdateServices,
		DeleteContext: resourceDeleteServices,
		ReadContext:   resourceReadServices,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"escalation_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeInt,
				Optional: true,
			},
			"collation_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sla": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func CreateServices(Ctx context.Context, d *schema.ResourceData, m interface{}) *client.Services {
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

	return new_service
}

func resourceCreateServices(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	new_service := CreateServices(Ctx, d, m)
	var diags diag.Diagnostics
	service, err := apiclient.Services.CreateService(team_id, new_service)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(service.Unique_Id)
	return diags
}

func resourceUpdateServices(Ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()

	team_id := d.Get("team_id").(string)
	id := d.Id()
	if team_id == "" {
		return diag.FromErr(errors.New("team_id is required"))
	}
	new_service := CreateServices(Ctx, d, m)

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
	if team_id == "" {
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

	return diags
}
