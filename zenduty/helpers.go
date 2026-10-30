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

// localTimeLayout is the YYYY-MM-DD HH:MM wall-clock format resources use.
const localTimeLayout = "2006-01-02 15:04"

// parseAPITime parses the API's ISO 8601 timestamps.
func parseAPITime(value string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised timestamp %q", value)
}

// localTimeToAPI converts a local wall-clock time in zone to UTC RFC 3339.
func localTimeToAPI(value, zone string) (string, error) {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return "", fmt.Errorf("invalid time zone %q: %w", zone, err)
	}
	t, err := time.ParseInLocation(localTimeLayout, value, loc)
	if err != nil {
		return "", fmt.Errorf("expected %q to be a time in the format YYYY-MM-DD HH:MM", value)
	}
	// a time inside a daylight-saving gap would never read back as written
	if t.Format(localTimeLayout) != value {
		return "", fmt.Errorf("%q does not exist in time zone %s (it falls in a daylight-saving gap); use a time outside the gap", value, zone)
	}
	return t.UTC().Format(time.RFC3339), nil
}

// apiTimeToLocal converts an API timestamp to localTimeLayout in zone.
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

// localPeriodToAPI converts a start and end in zone to UTC and requires end to be after start.
func localPeriodToAPI(startKey, start, endKey, end, zone string) (string, string, error) {
	from, err := localTimeToAPI(start, zone)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", startKey, err)
	}
	to, err := localTimeToAPI(end, zone)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", endKey, err)
	}
	if to <= from {
		return "", "", fmt.Errorf("%s must be after %s", endKey, startKey)
	}
	return from, to, nil
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
				Detail:   fmt.Sprintf("expected a time in the format YYYY-MM-DD HH:MM, got %q", s),
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
				Detail:   fmt.Sprintf("expected an IANA time zone name such as UTC or Europe/Berlin, got %q", s),
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
				Detail:   "expected valid JSON; use jsonencode()",
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
				Detail:   "expected a JSON object; use jsonencode({...})",
			}}
		}
		return nil
	}
}

// expandStringSet converts a TypeSet of strings to a slice.
func expandStringSet(v interface{}) []string {
	set, ok := v.(*schema.Set)
	if !ok || set == nil {
		return []string{}
	}
	return expandStringList(set.List())
}

// setJSONAttribute stores API JSON normalised so DiffSuppress compares equal.
func setJSONAttribute(d *schema.ResourceData, key string, raw []byte) error {
	return d.Set(key, normalizeRawJSON(raw))
}

// normalizeRawJSON returns compact JSON, or "" for empty or null.
func normalizeRawJSON(raw []byte) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	if normalized, err := normalizeJSON(string(raw)); err == nil {
		return normalized
	}
	return string(raw)
}

// isSetInConfig reports whether the attribute is present in the configuration.
func isSetInConfig(d *schema.ResourceData, key string) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return false
	}
	attr := raw.GetAttr(key)
	return !attr.IsNull() && attr.IsKnown()
}

func configuredBool(d *schema.ResourceData, key string) *bool {
	if !isSetInConfig(d, key) {
		return nil
	}
	v := d.Get(key).(bool)
	return &v
}

func configuredInt(d *schema.ResourceData, key string) *int {
	if !isSetInConfig(d, key) {
		return nil
	}
	v := d.Get(key).(int)
	return &v
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
