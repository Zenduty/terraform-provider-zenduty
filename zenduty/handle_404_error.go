package zenduty

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// wrapReadWith404 wraps any ReadContextFunc to handle 404s globally
func wrapReadWith404(readFunc schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		diags := readFunc(ctx, d, m)

		for _, diagItem := range diags {
			if isNotFound(diagItem) {
				log.Printf("[WARN] Removing resource %s because it's gone", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diags
	}
}

// isNotFound detects a 404 from the error or message
func isNotFound(diagItem diag.Diagnostic) bool {
	if diagItem.Severity != diag.Error {
		return false
	}

	re := regexp.MustCompile("404 Not Found")
	return re.MatchString(diagItem.Summary) || re.MatchString(diagItem.Detail)
}
