package zenduty

import (
	"context"
	"strings"
	"testing"

	"github.com/Zenduty/zenduty-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testEspTeamID = "1cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"

func TestValidateEspAssignmentSettings(t *testing.T) {
	cases := []struct {
		name       string
		strategy   int
		notifyOnly bool
		wantErr    string
	}{
		{"any", client.AssigneeStrategyAny, false, ""},
		{"round robin", client.AssigneeStrategyRoundRobin, false, ""},
		{"round robin, notify assignee only", client.AssigneeStrategyRoundRobin, true, ""},
		{"notify assignee only without round robin", client.AssigneeStrategyAny, true, "notify_round_robin_assignee_only can only be true"},
		{"zero strategy", 0, false, "assignee_strategy must be"},
		{"unknown strategy", 3, false, "assignee_strategy must be"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateEspAssignmentSettings(c.strategy, c.notifyOnly)
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

// Invalid combinations must fail at plan time.
func TestEspPlanTimeValidation(t *testing.T) {
	r := resourceEsp()
	base := func() map[string]interface{} {
		return map[string]interface{}{"name": "esp", "team_id": testEspTeamID}
	}
	cases := []struct {
		name    string
		mutate  func(map[string]interface{})
		wantErr string
	}{
		{"defaults", func(map[string]interface{}) {}, ""},
		{"round robin", func(c map[string]interface{}) { c["assignee_strategy"] = 2 }, ""},
		{"round robin, notify assignee only", func(c map[string]interface{}) {
			c["assignee_strategy"] = 2
			c["notify_round_robin_assignee_only"] = true
		}, ""},
		{"notify assignee only without round robin", func(c map[string]interface{}) {
			c["notify_round_robin_assignee_only"] = true
		}, "notify_round_robin_assignee_only can only be true"},
		{"unknown strategy", func(c map[string]interface{}) { c["assignee_strategy"] = 3 }, "assignee_strategy must be"},
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

// The schema must reject unknown strategies.
func TestEspSchemaRejectsUnknownStrategy(t *testing.T) {
	r := resourceEsp()
	cfg := map[string]interface{}{"name": "esp", "team_id": testEspTeamID, "assignee_strategy": 3}
	diags := r.Validate(terraform.NewResourceConfigRaw(cfg))
	if !diags.HasError() {
		t.Fatal("expected a validation error for assignee_strategy = 3")
	}
}

// assignment_settings must always be present in the payload.
func TestCreateEspAlwaysSendsAssignmentSettings(t *testing.T) {
	r := resourceEsp()
	cases := []struct {
		name           string
		raw            map[string]interface{}
		wantStrategy   int
		wantNotifyOnly bool
	}{
		{"defaults", map[string]interface{}{}, client.AssigneeStrategyAny, false},
		{"round robin", map[string]interface{}{"assignee_strategy": 2}, client.AssigneeStrategyRoundRobin, false},
		{"round robin, notify assignee only", map[string]interface{}{"assignee_strategy": 2, "notify_round_robin_assignee_only": true}, client.AssigneeStrategyRoundRobin, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw := map[string]interface{}{"name": "esp", "team_id": testEspTeamID}
			for k, v := range c.raw {
				raw[k] = v
			}
			d := schema.TestResourceDataRaw(t, r.Schema, raw)
			esp, diags := CreateEsp(context.Background(), d, nil)
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if esp.AssignmentSettings == nil {
				t.Fatal("assignment_settings missing from the payload")
			}
			if esp.AssignmentSettings.AssigneeStrategy != c.wantStrategy || esp.AssignmentSettings.CallRRAssignee != c.wantNotifyOnly {
				t.Errorf("assignment_settings = %+v, want strategy %d notify_round_robin_assignee_only %v",
					*esp.AssignmentSettings, c.wantStrategy, c.wantNotifyOnly)
			}
		})
	}

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name": "esp", "team_id": testEspTeamID, "notify_round_robin_assignee_only": true,
	})
	if _, diags := CreateEsp(context.Background(), d, nil); !diags.HasError() {
		t.Fatal("expected notify_round_robin_assignee_only without round-robin to be rejected")
	}
}

func TestFlattenEspAssignmentSettings(t *testing.T) {
	cases := []struct {
		name           string
		settings       *client.EscalationPolicyAssignmentSettings
		wantStrategy   int
		wantNotifyOnly bool
	}{
		{"missing", nil, client.AssigneeStrategyAny, false},
		{"zero value", &client.EscalationPolicyAssignmentSettings{}, client.AssigneeStrategyAny, false},
		{"round robin", &client.EscalationPolicyAssignmentSettings{AssigneeStrategy: 2, CallRRAssignee: true}, client.AssigneeStrategyRoundRobin, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			strategy, notifyOnly := flattenEspAssignmentSettings(c.settings)
			if strategy != c.wantStrategy || notifyOnly != c.wantNotifyOnly {
				t.Errorf("got (%d, %v), want (%d, %v)", strategy, notifyOnly, c.wantStrategy, c.wantNotifyOnly)
			}
		})
	}
}

func TestFlattenEspItemIncludesPolicyFlags(t *testing.T) {
	item := flattenEspItem(&client.EscalationPolicy{
		UniqueID:     "esp-1",
		Name:         "esp",
		RepeatPolicy: 3,
		MoveToNext:   true,
		GlobalEp:     true,
		Rules:        []client.Rules{{Delay: 0, Position: 1, Targets: []client.Targets{{TargetType: 2, TargetID: "user", Position: 0}}}},
		AssignmentSettings: &client.EscalationPolicyAssignmentSettings{
			AssigneeStrategy: client.AssigneeStrategyRoundRobin,
			CallRRAssignee:   true,
		},
	})
	want := map[string]interface{}{
		"repeat_policy":                    3,
		"move_to_next":                     true,
		"global_ep":                        true,
		"assignee_strategy":                client.AssigneeStrategyRoundRobin,
		"notify_round_robin_assignee_only": true,
	}
	for k, v := range want {
		if item[k] != v {
			t.Errorf("%s = %v, want %v", k, item[k], v)
		}
	}
	if len(item["rules"].([]map[string]interface{})) != 1 {
		t.Errorf("rules were not flattened")
	}
}
