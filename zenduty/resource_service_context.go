package zenduty

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceContext() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateServiceContext,
		ReadContext:   resourceReadServiceContext,
		DeleteContext: resourceDeleteServiceContext,
		Importer: &schema.ResourceImporter{
			State: resourceServiceContextImporter,
		},
		Description: "Links a context integration to a service. Every attribute forces a new resource; the API has no update.",
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"service_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
			},
			"integration": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateUUID(),
				Description:      "unique_id of the context integration. A service can link each integration once.",
			},
		},
	}
}

func resourceCreateServiceContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.ServiceContext.CreateServiceContext(d.Get("team_id").(string), d.Get("service_id").(string), &client.ServiceContext{
		Integration: d.Get("integration").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadServiceContext(ctx, d, m)
}

func resourceReadServiceContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	link, err := apiclient.ServiceContext.GetServiceContext(d.Get("team_id").(string), d.Get("service_id").(string), d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	d.Set("integration", link.Integration)
	if link.Service != "" {
		d.Set("service_id", link.Service)
	}
	return nil
}

func resourceDeleteServiceContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.ServiceContext.DeleteServiceContext(d.Get("team_id").(string), d.Get("service_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceServiceContextImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <team_id>/<service_id>/<service_context_id>", d.Id())
	} else if !IsValidUUID(parts[0]) {
		return nil, fmt.Errorf("invalid team_id (%q)", parts[0])
	} else if !IsValidUUID(parts[1]) {
		return nil, fmt.Errorf("invalid service_id (%q)", parts[1])
	} else if !IsValidUUID(parts[2]) {
		return nil, fmt.Errorf("invalid service_context_id (%q)", parts[2])
	}
	d.Set("team_id", parts[0])
	d.Set("service_id", parts[1])
	d.SetId(parts[2])
	return []*schema.ResourceData{d}, nil
}
