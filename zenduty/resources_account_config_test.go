package zenduty

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testUUID     = "1cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	testUUID2    = "2cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	testUUID3    = "3cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	testUsername = "1cd4b3f5-5f4c-4d1e-9c2a-0"
)

type importerFunc func(*schema.ResourceData, interface{}) ([]*schema.ResourceData, error)

// Every new importer must accept its documented id format and reject others.
func TestImportersParseDocumentedIDs(t *testing.T) {
	cases := []struct {
		name     string
		resource *schema.Resource
		importer importerFunc
		good     string
		bad      []string
		want     map[string]string // attributes set from the id
		wantID   string
	}{
		{"heartbeat", resourceHeartbeat(), resourceHeartbeatImporter, testUUID, []string{"x", testUUID + "/" + testUUID}, nil, testUUID},
		{"webhook_configuration", resourceWebhookConfiguration(), resourceWebhookConfigurationImporter, testUUID, []string{"x"}, nil, testUUID},
		{"account_stakeholder_template", resourceAccountStakeholderTemplate(), resourceAccountStakeholderTemplateImporter, testUUID, []string{"x"}, nil, testUUID},
		{"incident_custom_field", resourceIncidentCustomField(), resourceIncidentCustomFieldImporter, testUUID, []string{"x"}, nil, testUUID},
		{"incident_type", resourceIncidentType(), resourceIncidentTypeImporter, testUUID, []string{"x"}, nil, testUUID},
		{"workflow", resourceWorkflow(), resourceWorkflowImporter, testUUID, []string{"x"}, nil, testUUID},
		{"maintenance_template", resourceMaintenanceTemplate(), resourceMaintenanceTemplateImporter, testUUID + "/" + testUUID2,
			[]string{testUUID, "x/" + testUUID2, testUUID + "/x", testUUID + "/" + testUUID2 + "/" + testUUID3},
			map[string]string{"team_id": testUUID}, testUUID2},
		{"service_maintenance_window", resourceServiceMaintenanceWindow(), resourceServiceMaintenanceWindowImporter, testUUID + "/" + testUUID2 + "/" + testUUID3,
			[]string{testUUID + "/" + testUUID2, "x/" + testUUID2 + "/" + testUUID3, testUUID + "/x/" + testUUID3, testUUID + "/" + testUUID2 + "/x"},
			map[string]string{"team_id": testUUID, "service_id": testUUID2}, testUUID3},
		{"service_context", resourceServiceContext(), resourceServiceContextImporter, testUUID + "/" + testUUID2 + "/" + testUUID3,
			[]string{testUUID + "/" + testUUID2, testUUID + "/" + testUUID2 + "/x"},
			map[string]string{"team_id": testUUID, "service_id": testUUID2}, testUUID3},
		{"user_forwarding_rule", resourceUserForwardingRule(), resourceUserForwardingRuleImporter, testUsername + "/" + testUUID,
			[]string{testUUID, testUUID + "/" + testUUID, testUsername + "/x", testUsername},
			map[string]string{"username": testUsername}, testUUID},
		{"assign_account_role", resourceAssignAccountRole(), resourceAssignAccountRoleImporter, testUsername,
			[]string{testUUID, "x", testUsername + "/" + testUUID},
			map[string]string{"username": testUsername}, testUsername},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, id := range c.bad {
				d := c.resource.Data(nil)
				d.SetId(id)
				if _, err := c.importer(d, nil); err == nil {
					t.Errorf("import %q: expected an error", id)
				}
			}
			d := c.resource.Data(nil)
			d.SetId(c.good)
			if _, err := c.importer(d, nil); err != nil {
				t.Fatalf("import of %q failed: %v", c.good, err)
			}
			if d.Id() != c.wantID {
				t.Errorf("id = %q, want %q", d.Id(), c.wantID)
			}
			for key, want := range c.want {
				if got := d.Get(key).(string); got != want {
					t.Errorf("%s = %q, want %q", key, got, want)
				}
			}
		})
	}
}

// The account settings singleton accepts any id and normalises it.
func TestAccountSettingImporterNormalisesID(t *testing.T) {
	d := resourceAccountSetting().Data(nil)
	d.SetId("whatever")
	if _, err := resourceAccountSettingImporter(d, nil); err != nil {
		t.Fatal(err)
	}
	if d.Id() != accountSettingID {
		t.Errorf("id = %q, want %q", d.Id(), accountSettingID)
	}
}

func planError(t *testing.T, r *schema.Resource, cfg map[string]interface{}) error {
	t.Helper()
	_, err := r.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(cfg), nil)
	return err
}

func expectPlanError(t *testing.T, err error, want string) {
	t.Helper()
	if want == "" {
		if err != nil {
			t.Fatalf("unexpected plan error: %v", err)
		}
		return
	}
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("plan error = %v, want it to contain %q", err, want)
	}
}

// The descriptor rules the API enforces must fail at plan time.
func TestIncidentCustomFieldPlanTimeValidation(t *testing.T) {
	r := resourceIncidentCustomField()
	cases := []struct {
		name       string
		fieldType  int
		properties string
		wantErr    string
	}{
		{"string ok", 0, `{"type":"string","title":"Region"}`, ""},
		{"select ok", 3, `{"type":"string","title":"Env","enum":["prod","dev"]}`, ""},
		{"multiselect ok", 4, `{"type":"array","title":"Tags","items":{"enum":["a"]}}`, ""},
		{"missing title", 0, `{"type":"string"}`, "non-empty title"},
		{"type mismatch", 2, `{"type":"string","title":"Count"}`, `requires properties.type "integer"`},
		{"select without enum", 3, `{"type":"string","title":"Env"}`, "non-empty enum"},
		{"multiselect without items", 4, `{"type":"array","title":"Tags"}`, "items.enum"},
		{"not an object", 0, `"text"`, "JSON object"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := planError(t, r, map[string]interface{}{"field_type": c.fieldType, "properties": c.properties})
			expectPlanError(t, err, c.wantErr)
		})
	}
}

func TestUserForwardingRulePlanTimeValidation(t *testing.T) {
	r := resourceUserForwardingRule()
	base := func() map[string]interface{} {
		return map[string]interface{}{
			"username":  testUsername,
			"to_user":   testUsername,
			"from_date": "2027-01-10 09:00",
			"to_date":   "2027-01-12 09:00",
			"time_zone": "Europe/Berlin",
		}
	}
	expectPlanError(t, planError(t, r, base()), "")

	bad := base()
	bad["to_date"] = "2027-01-10 08:00"
	expectPlanError(t, planError(t, r, bad), "to_date must be after from_date")

	badFormat := base()
	badFormat["from_date"] = "10/01/2027 09:00"
	expectPlanError(t, planError(t, r, badFormat), "YYYY-MM-DD HH:MM")

	badZone := base()
	badZone["time_zone"] = "Mars/Olympus"
	expectPlanError(t, planError(t, r, badZone), "time zone")
}

func TestWorkflowPlanTimeValidation(t *testing.T) {
	r := resourceWorkflow()
	base := func(schemaJSON string) map[string]interface{} {
		return map[string]interface{}{"name": "w", "description": "d", "json_schema": schemaJSON}
	}
	expectPlanError(t, planError(t, r, base(`{"nodes":[{"id":"t","type":"trigger"},{"id":"a","type":"action"}],"edges":[]}`)), "")
	expectPlanError(t, planError(t, r, base(`{"nodes":[{"id":"a","type":"action"}],"edges":[]}`)), "exactly one node of type \"trigger\", found 0")
	expectPlanError(t, planError(t, r, base(`{"edges":[]}`)), "needs a nodes list")
	expectPlanError(t, planError(t, r, base(`[]`)), "JSON object")
}

func TestServiceMaintenanceWindowExpandConvertsToUTC(t *testing.T) {
	d := resourceServiceMaintenanceWindow().Data(nil)
	d.Set("name", "patching")
	d.Set("timezone", "Europe/Berlin")
	d.Set("start_time", "2027-01-15 10:00")
	d.Set("end_time", "2027-01-15 12:00")
	d.Set("repeat_until", "2027-02-15 12:00")
	window, err := expandServiceMaintenanceWindow(d)
	if err != nil {
		t.Fatal(err)
	}
	if window.StartTime != "2027-01-15T09:00:00Z" || window.EndTime != "2027-01-15T11:00:00Z" || window.RepeatUntil == nil || *window.RepeatUntil != "2027-02-15T11:00:00Z" {
		t.Errorf("got start=%s end=%s until=%v", window.StartTime, window.EndTime, window.RepeatUntil)
	}
	d.Set("repeat_until", "")
	if window, err = expandServiceMaintenanceWindow(d); err != nil || window.RepeatUntil != nil {
		t.Errorf("an unset repeat_until must be sent as null, got %v, %v", window.RepeatUntil, err)
	}
	if window.RepeatInterval != 0 {
		t.Errorf("repeat_interval must be sent explicitly as 0, got %d", window.RepeatInterval)
	}
	d.Set("end_time", "2027-01-15 09:00")
	if _, err := expandServiceMaintenanceWindow(d); err == nil || !strings.Contains(err.Error(), "end_time must be after start_time") {
		t.Errorf("expected ordering error, got %v", err)
	}
}

func TestTimeHelpersRoundTrip(t *testing.T) {
	api, err := localTimeToAPI("2027-07-01 10:30", "Europe/Berlin")
	if err != nil || api != "2027-07-01T08:30:00Z" {
		t.Fatalf("got %q, %v", api, err)
	}
	for _, variant := range []string{"2027-07-01T08:30:00Z", "2027-07-01T08:30:00.123456Z", "2027-07-01T10:30:00+02:00"} {
		local, err := apiTimeToLocal(variant, "Europe/Berlin")
		if err != nil || local != "2027-07-01 10:30" {
			t.Errorf("%s: got %q, %v", variant, local, err)
		}
	}
	if local, err := apiTimeToLocal("", "UTC"); err != nil || local != "" {
		t.Errorf("empty timestamp must stay empty, got %q, %v", local, err)
	}
	// 2027-03-28 02:30 does not exist in Berlin (clocks jump 02:00 -> 03:00)
	if _, err := localTimeToAPI("2027-03-28 02:30", "Europe/Berlin"); err == nil || !strings.Contains(err.Error(), "daylight-saving gap") {
		t.Errorf("a DST-gap time must be rejected, got %v", err)
	}
	if _, err := localTimeToAPI("2027-03-28 03:30", "Europe/Berlin"); err != nil {
		t.Errorf("a time after the gap must pass, got %v", err)
	}
}

func TestIncidentTypeFlattenAndExpand(t *testing.T) {
	incidentType := &client.IncidentType{
		Name: "Outage",
		FormToggles: []map[string]interface{}{
			{"impacted_services": map[string]interface{}{"enabled": true, "ui:widget": "x"}},
			{"responders": false},
		},
		CustomFields: client.IncidentTypeCustomFieldsSchema{
			Type: "object",
			Properties: map[string]client.IncidentTypeMappedField{
				"custom_field_9": {Position: 2, Default: json.RawMessage(`["a","b"]`)},
				"custom_field_7": {Position: 1, Default: json.RawMessage(`"eu"`)},
				"custom_field_0": {Position: 0}, // definition deleted: not in keyToID
			},
			Required: []string{"custom_field_7"},
		},
	}
	keyToID := map[string]string{"custom_field_7": testUUID, "custom_field_9": testUUID2}

	toggles := flattenFormToggles(incidentType.FormToggles)
	if toggles["impacted_services"] != true || toggles["responders"] != false {
		t.Errorf("toggles = %v", toggles)
	}
	blocks := flattenIncidentTypeCustomFields(incidentType, keyToID)
	if len(blocks) != 2 {
		t.Fatalf("blocks = %v", blocks)
	}
	first := blocks[0].(map[string]interface{})
	second := blocks[1].(map[string]interface{})
	if first["custom_field"] != testUUID || first["required"] != true || first["default_value"] != `"eu"` {
		t.Errorf("first = %v", first)
	}
	if second["custom_field"] != testUUID2 || second["required"] != false || second["default_value"] != `["a","b"]` {
		t.Errorf("second = %v", second)
	}

	d := resourceIncidentType().Data(nil)
	d.Set("name", "Outage")
	d.Set("form_toggles", map[string]interface{}{"responders": true})
	d.Set("custom_fields", []interface{}{
		map[string]interface{}{"custom_field": testUUID, "required": true, "default_value": `"eu"`},
		map[string]interface{}{"custom_field": testUUID2},
	})
	input := expandIncidentTypeInput(d)
	body, _ := json.Marshal(input)
	s := string(body)
	for _, want := range []string{`"form_toggles":[{"responders":true}]`, `"unique_id":"` + testUUID + `","default_value":"eu","required":true,"position":0`, `"unique_id":"` + testUUID2 + `","required":false,"position":1`} {
		if !strings.Contains(s, want) {
			t.Errorf("body %s lacks %s", s, want)
		}
	}
}

func TestWorkflowExpandAndFlattenTeams(t *testing.T) {
	d := resourceWorkflow().Data(nil)
	d.Set("name", "w")
	d.Set("description", "d")
	d.Set("json_schema", `{"nodes":[],"edges":[]}`)
	d.Set("teams", []interface{}{
		map[string]interface{}{"team_id": testUUID, "service_ids": []interface{}{testUUID2}},
		map[string]interface{}{"team_id": testUUID3},
	})
	input := expandWorkflowInput(d)
	if len(input.Teams) != 2 {
		t.Fatalf("teams = %+v", input.Teams)
	}
	body, _ := json.Marshal(input)
	if !strings.Contains(string(body), `"unique_id":"`+testUUID+`","services":["`+testUUID2+`"]`) || !strings.Contains(string(body), `"unique_id":"`+testUUID3+`","services":[]`) {
		t.Errorf("body = %s", body)
	}

	flat := flattenWorkflowTeams([]client.WorkflowTeam{{UniqueID: testUUID, Services: []client.WorkflowTeamService{{UniqueID: testUUID2}}}})
	if err := d.Set("teams", flat); err != nil {
		t.Fatalf("set teams: %v", err)
	}
	if d.Get("teams").(*schema.Set).Len() != 1 {
		t.Errorf("teams = %v", d.Get("teams"))
	}
}

// Webhook secrets are masked by the API, so the read must not overwrite them.
func TestWebhookConfigurationExpandUsesSet(t *testing.T) {
	d := resourceWebhookConfiguration().Data(nil)
	d.Set("url", "https://example.com/hook")
	d.Set("event", []interface{}{"incident_create", "oncall_change"})
	d.Set("api_key", "secret")
	config := expandWebhookConfiguration(d)
	if len(config.Event) != 2 || config.APIKey != "secret" {
		t.Errorf("config = %+v", config)
	}
}

// A 403 from the custom-role endpoint means the user can no longer hold a
// custom role; the assignment must leave state with a warning, not fail plans.
func TestAssignAccountRoleReadDropsForbiddenAssignment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/account/users/"+testUsername+"/customroles/" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"detail": "Unauthorized"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	d := resourceAssignAccountRole().Data(nil)
	d.SetId(testUsername)
	d.Set("username", testUsername)
	d.Set("account_role", testUUID)
	diags := resourceReadRole(context.Background(), d, &Config{Token: "t", BaseURL: server.URL})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if len(diags) != 1 || diags[0].Severity != 0+1 /* diag.Warning */ {
		t.Fatalf("expected one warning, got %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("resource must leave state, id still %q", d.Id())
	}
}
