package zenduty

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func isJSONString(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// IsValidUsername matches Zenduty usernames, which are the first 25 characters
// of a UUID (8-4-4-4-1) — not a full UUID.
func IsValidUsername(username string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]$")
	return r.MatchString(username)
}

func ValidateUserName() schema.SchemaValidateDiagFunc {
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
		if !IsValidUsername(id) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid username",
				Detail:   fmt.Sprintf("expected %s to be a valid Zenduty username (the 25-character prefix of a UUID)", path),
			})
		}

		return diags
	}
}

func emptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// isRetryableError reports whether an SDK error is worth retrying: 5xx, 429,
// or transport-level failures (timeouts, connection resets). A transport
// error carries no HTTP status because it never reached the API, and is safe
// to retry.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	if code := client.StatusCode(err); code != 0 {
		return code == http.StatusTooManyRequests || code >= 500
	}
	return true
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
	emailAddress, err := mail.ParseAddress(e)
	return err == nil && emailAddress.Address == e
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

// ValidateNonZeroLength rejects only zero-length strings. Unlike
// ValidateRequired it accepts whitespace-only values, which the Zenduty API
// permits for some fields (e.g. a user's last name) — rejecting them would
// make real remote state unrepresentable in config.
func ValidateNonZeroLength() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		s, ok := v.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid",
				Detail:   "expected type of string",
			})
			return diags
		}
		if len(s) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "This field is required",
				Detail:   "This field cannot be empty",
			})
		}
		return diags
	}
}

func generateUUID() string {
	id := uuid.New()
	uuidString := id.String()
	return uuidString
}

func normalizeJSON(jsonString string) (string, error) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		return "", err
	}

	normalizedBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	return string(normalizedBytes), nil
}

// suppressEquivalentJSONDiffs suppresses diffs between JSON strings that
// decode to the same value. Reads store the API's compact, key-sorted JSON, so
// without this any human-formatted config diffs on every plan.
func suppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	normalizedOld, errOld := normalizeJSON(old)
	normalizedNew, errNew := normalizeJSON(new)
	if errOld != nil || errNew != nil {
		return old == new
	}
	return normalizedOld == normalizedNew
}

// suppressEquivalentJSONObjectDiffs is suppressEquivalentJSONDiffs for fields
// whose backend default is an empty object, where an unset value and "{}" mean
// the same thing.
func suppressEquivalentJSONObjectDiffs(k, old, new string, d *schema.ResourceData) bool {
	if old == "" {
		old = "{}"
	}
	if new == "" {
		new = "{}"
	}
	return suppressEquivalentJSONDiffs(k, old, new, d)
}

// localTimeLayout is how resources accept and show wall-clock times: a
// "YYYY-MM-DD HH:MM" string interpreted in the resource's time zone.
const localTimeLayout = "2006-01-02 15:04"

// parseAPITime parses the timestamps the API returns: ISO 8601 in UTC with a
// "Z" suffix or an explicit offset, with or without fractional seconds.
func parseAPITime(value string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised timestamp %q", value)
}

// localTimeToAPI converts a localTimeLayout string in zone to the UTC
// RFC 3339 form the API expects.
func localTimeToAPI(value, zone string) (string, error) {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return "", fmt.Errorf("invalid time zone %q: %w", zone, err)
	}
	t, err := time.ParseInLocation(localTimeLayout, value, loc)
	if err != nil {
		return "", fmt.Errorf("expected %q to be a time in the format YYYY-MM-DD HH:MM", value)
	}
	// a wall-clock time inside a daylight-saving gap is normalised forward by
	// ParseInLocation and would never read back as written
	if t.Format(localTimeLayout) != value {
		return "", fmt.Errorf("%q does not exist in time zone %s (it falls in a daylight-saving gap); use a time outside the gap", value, zone)
	}
	return t.UTC().Format(time.RFC3339), nil
}

// apiTimeToLocal converts an API timestamp to localTimeLayout in zone. An
// empty timestamp (a null on the API) stays empty.
func apiTimeToLocal(value, zone string) (string, error) {
	if value == "" {
		return "", nil
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return "", fmt.Errorf("invalid time zone %q: %w", zone, err)
	}
	t, err := parseAPITime(value)
	if err != nil {
		return "", err
	}
	return t.In(loc).Format(localTimeLayout), nil
}

// ValidateLocalTime checks the YYYY-MM-DD HH:MM format at plan time.
func ValidateLocalTime() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		s, ok := v.(string)
		if !ok {
			return diag.Diagnostics{{Severity: diag.Error, Summary: "Invalid", Detail: "expected type of string"}}
		}
		if _, err := time.Parse(localTimeLayout, s); err != nil {
			return diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Invalid time",
				Detail:   fmt.Sprintf("expected %s to be a time in the format YYYY-MM-DD HH:MM, got %q", path, s),
			}}
		}
		return nil
	}
}

// ValidateTimeZone checks that the value is an IANA time zone name.
func ValidateTimeZone() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		s, ok := v.(string)
		if !ok {
			return diag.Diagnostics{{Severity: diag.Error, Summary: "Invalid", Detail: "expected type of string"}}
		}
		if _, err := time.LoadLocation(s); err != nil || s == "" {
			return diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Invalid time zone",
				Detail:   fmt.Sprintf("expected %s to be an IANA time zone name such as UTC or Europe/Berlin, got %q", path, s),
			}}
		}
		return nil
	}
}

// ValidateJSONString rejects values that are not valid JSON.
func ValidateJSONString() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		s, ok := v.(string)
		if !ok {
			return diag.Diagnostics{{Severity: diag.Error, Summary: "Invalid", Detail: "expected type of string"}}
		}
		if !isJSONString(s) {
			return diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Invalid JSON",
				Detail:   fmt.Sprintf("expected %s to be valid JSON; use jsonencode()", path),
			}}
		}
		return nil
	}
}

// ValidateJSONObjectString rejects values that are not a JSON object.
func ValidateJSONObjectString() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		s, ok := v.(string)
		if !ok {
			return diag.Diagnostics{{Severity: diag.Error, Summary: "Invalid", Detail: "expected type of string"}}
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(s), &obj); err != nil || obj == nil {
			return diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Invalid JSON object",
				Detail:   fmt.Sprintf("expected %s to be a JSON object; use jsonencode({...})", path),
			}}
		}
		return nil
	}
}

// expandStringSet turns a TypeSet of strings into a sorted-by-set-order slice.
func expandStringSet(v interface{}) []string {
	set, ok := v.(*schema.Set)
	if !ok || set == nil {
		return []string{}
	}
	return expandStringList(set.List())
}

// setJSONAttribute stores raw JSON from the API in its normalised form so it
// compares equal to any equivalent config value under suppressEquivalentJSONDiffs.
func setJSONAttribute(d *schema.ResourceData, key string, raw []byte) error {
	if len(raw) == 0 || string(raw) == "null" {
		return d.Set(key, "")
	}
	normalized, err := normalizeJSON(string(raw))
	if err != nil {
		return d.Set(key, string(raw))
	}
	return d.Set(key, normalized)
}

// jsonStringField returns the string value of key in a JSON object, or "".
func jsonStringField(raw []byte, key string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ""
	}
	s, _ := obj[key].(string)
	return s
}
