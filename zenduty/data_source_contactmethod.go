package zenduty

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceUserContacts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationRulesRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 5),
			},

			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNotificationRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiclient, _ := m.(*Config).Client()
	username := d.Get("user_id").(string)
	if emptyString(username) {
		return diag.Errorf("username is required")
	}

	contactType := d.Get("contact_type").(int)
	value := d.Get("value").(string)

	var diags diag.Diagnostics

	contactmethod, err := apiclient.ContactMethod.GetContactMethods(username)
	if err != nil {
		return diag.Errorf("Error getting contactmethod: %s", err)
	}
	if len(contactmethod) == 0 {
		return diag.Errorf("No contactmethod found for contact_type %s", strconv.Itoa(contactType))
	}
	if value != "" {
		for _, contactmethod := range contactmethod {
			if contactmethod.Value == value && contactmethod.ContactType == contactType {
				d.SetId(contactmethod.UniqueID)
				d.Set("value", contactmethod.Value)
				d.Set("name", contactmethod.Name)
				return diags
			}
		}
		return diag.Errorf("No contactmethod found for contact_type %s and value %s", strconv.Itoa(contactType), value)
	}
	for _, contactmethod := range contactmethod {
		if contactmethod.ContactType == contactType {
			d.SetId(contactmethod.UniqueID)
			d.Set("value", contactmethod.Value)
			d.Set("name", contactmethod.Name)
			return diags
		}
	}
	return diag.Errorf("No contactmethod found for contact_type %s", strconv.Itoa(contactType))

}
