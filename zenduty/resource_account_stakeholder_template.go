package zenduty

import (
	"context"
	"fmt"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccountStakeholderTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAccountStakeholderTemplate,
		ReadContext:   resourceReadAccountStakeholderTemplate,
		UpdateContext: resourceUpdateAccountStakeholderTemplate,
		DeleteContext: resourceDeleteAccountStakeholderTemplate,
		Importer: &schema.ResourceImporter{
			State: resourceAccountStakeholderTemplateImporter,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
				Description:      "Title of the template, unique within the account.",
			},
			"subject": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateNonZeroLength(),
				Description:      "Subject line of the stakeholder update.",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Body of the stakeholder update.",
			},
		},
	}
}

func expandAccountStakeholderTemplate(d *schema.ResourceData) *client.AccountStakeholderTemplate {
	return &client.AccountStakeholderTemplate{
		Title:    d.Get("title").(string),
		Subject:  d.Get("subject").(string),
		Template: d.Get("template").(string),
	}
}

func resourceCreateAccountStakeholderTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	created, err := apiclient.AccountStakeholderTemplate.CreateAccountStakeholderTemplate(expandAccountStakeholderTemplate(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.UniqueID)
	return resourceReadAccountStakeholderTemplate(ctx, d, m)
}

func resourceUpdateAccountStakeholderTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if _, err := apiclient.AccountStakeholderTemplate.UpdateAccountStakeholderTemplate(d.Id(), expandAccountStakeholderTemplate(d)); err != nil {
		return diag.FromErr(err)
	}
	return resourceReadAccountStakeholderTemplate(ctx, d, m)
}

func resourceReadAccountStakeholderTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	template, err := apiclient.AccountStakeholderTemplate.GetAccountStakeholderTemplate(d.Id())
	if err != nil {
		return handleReadError(d, err)
	}
	d.Set("title", template.Title)
	d.Set("subject", template.Subject)
	d.Set("template", template.Template)
	return nil
}

func resourceDeleteAccountStakeholderTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiclient, _ := m.(*Config).Client()
	if err := apiclient.AccountStakeholderTemplate.DeleteAccountStakeholderTemplate(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAccountStakeholderTemplateImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if !IsValidUUID(d.Id()) {
		return nil, fmt.Errorf("unexpected format of id (%q), expected <stakeholder_template_id>", d.Id())
	}
	return []*schema.ResourceData{d}, nil
}
