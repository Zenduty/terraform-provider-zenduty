package zenduty

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Editing team_id on an existing service must fail at plan time with guidance,
// never fall through to the replacement that ForceNew would otherwise plan.
func TestServiceTeamChangeIsRejected(t *testing.T) {
	r := resourceServices()
	teamA := "1cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	teamB := "2cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	esp := "3cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"
	svc := "4cd4b3f5-5f4c-4d1e-9c2a-0a1b2c3d4e5f"

	config := func(team string) *terraform.ResourceConfig {
		return terraform.NewResourceConfigRaw(map[string]interface{}{
			"name": "payments", "team_id": team, "escalation_policy": esp,
		})
	}
	existing := &terraform.InstanceState{
		ID:         svc,
		Attributes: map[string]string{"id": svc, "name": "payments", "team_id": teamA, "escalation_policy": esp},
	}

	_, err := r.Diff(context.Background(), existing, config(teamB), nil)
	if err == nil {
		t.Fatal("changing team_id on an existing service: expected a plan error")
	}
	for _, want := range []string{"cannot be changed", teamA, teamB, teamB + "/" + svc, "terraform state rm", "terraform import"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should mention %q", err.Error(), want)
		}
	}

	if _, err := r.Diff(context.Background(), existing, config(teamA), nil); err != nil {
		t.Fatalf("unchanged team_id must plan cleanly: %v", err)
	}
	if _, err := r.Diff(context.Background(), nil, config(teamB), nil); err != nil {
		t.Fatalf("a new service must plan cleanly: %v", err)
	}
}
