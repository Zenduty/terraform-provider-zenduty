package zenduty

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func isJSONString(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func checkList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func validateDate(date string) bool {

	return true
}
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func emptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func ValidateUUID() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		id, ok := v.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid",
				Detail:   "expected type of string",
			})
		}
		if !IsValidUUID(id) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid ID",
				Detail:   fmt.Sprintf("expected %s to be a valid UUID", path),
			})
		}

		return diags
	}
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func ValidateEmail() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		id, ok := v.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid",
				Detail:   "expected type of string",
			})
		}
		if !isEmailValid(id) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid Email",
				Detail:   fmt.Sprintf("expected %s to be a valid Email", path),
			})
		}

		return diags
	}

}

func ValidateRequired() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		id, ok := v.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid",
				Detail:   "expected type of string",
			})
		}
		if emptyString(id) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "This field is required",
				Detail:   "This field cannot be empty",
			})
		}

		return diags
	}

}

func genrateUUID() string {
	id := uuid.New()
	uuidString := id.String()
	return uuidString
}
