package zenduty

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestValidateAlertGroupingFields(t *testing.T) {
	cases := []struct {
		name    string
		static  []string
		custom  []string
		wantErr string
	}{
		{"static only", []string{"incident_title"}, nil, ""},
		{"custom only", nil, []string{"message"}, ""},
		{"both", []string{"incident_title", "integration"}, []string{"message"}, ""},
		{"none", nil, nil, "at least one of static_fields or custom_fields"},
		{"empty lists", []string{}, []string{}, "at least one of static_fields or custom_fields"},
		{"duplicate within static", []string{"incident_title", "incident_title"}, nil, "static_fields lists \"incident_title\" more than once"},
		{"duplicate within custom", nil, []string{"host", "host"}, "custom_fields lists \"host\" more than once"},
		{"duplicate across lists", []string{"integration"}, []string{"integration"}, "listed in both static_fields and custom_fields"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateAlertGroupingFields(c.static, c.custom)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("error = %v, want it to contain %q", err, c.wantErr)
			}
		})
	}
}

func TestAlertGroupingPolicyImporterRejectsBadIDs(t *testing.T) {
	r := resourceAlertGroupingPolicy()
	team := "1cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	for _, id := range []string{
		"only-one-part",
		team + "/" + team,
		"not-a-uuid/" + team + "/" + team,
		team + "/not-a-uuid/" + team,
		team + "/" + team + "/not-a-uuid",
	} {
		d := r.Data(nil)
		d.SetId(id)
		if _, err := resourceAlertGroupingPolicyImporter(d, nil); err == nil {
			t.Errorf("import %q: expected an error", id)
		}
	}
	d := r.Data(nil)
	d.SetId(team + "/" + team + "/" + team)
	if _, err := resourceAlertGroupingPolicyImporter(d, nil); err != nil {
		t.Fatalf("import of a valid id failed: %v", err)
	}
	if d.Get("team_id") != team || d.Get("service_id") != team || d.Id() != team {
		t.Errorf("import did not split the id into team_id/service_id/id")
	}
}

// The field rules must fail at plan time, through CustomizeDiff, not only at apply.
func TestAlertGroupingPolicyPlanTimeValidation(t *testing.T) {
	r := resourceAlertGroupingPolicy()
	id := "1cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	base := func() map[string]interface{} {
		return map[string]interface{}{"team_id": id, "service_id": id, "time_window": 5}
	}
	cases := []struct {
		name    string
		mutate  func(map[string]interface{})
		wantErr string
	}{
		{"valid", func(c map[string]interface{}) { c["static_fields"] = []interface{}{"integration"} }, ""},
		{"no fields", func(c map[string]interface{}) {}, "at least one of static_fields or custom_fields"},
		{"across lists", func(c map[string]interface{}) {
			c["static_fields"] = []interface{}{"integration"}
			c["custom_fields"] = []interface{}{"integration"}
		}, "listed in both"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := base()
			c.mutate(cfg)
			_, err := r.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(cfg), nil)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected plan error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("plan error = %v, want it to contain %q", err, c.wantErr)
			}
		})
	}
}
