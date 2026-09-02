package authentik

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderEnrollmentBlueprintContainsFailClosedContract(t *testing.T) {
	first, err := RenderEnrollmentBlueprint(validApplication())
	if err != nil {
		t.Fatalf("RenderEnrollmentBlueprint() returned error: %v", err)
	}
	second, err := RenderEnrollmentBlueprint(validApplication())
	if err != nil || !bytes.Equal(first, second) {
		t.Fatal("blueprint rendering is not deterministic")
	}
	for _, required := range []string{
		"create_users_as_inactive: true", "activate_user_on_success: true",
		"create_users_group: !KeyOf access-group", "authentik_policies.policybinding",
		"failure_result: false", "[slug, example-app]",
	} {
		if !strings.Contains(string(first), required) {
			t.Errorf("blueprint missing %q", required)
		}
	}
	if output, err := RenderEnrollmentBlueprint(Application{}); err == nil || output != nil {
		t.Fatalf("invalid render = (%q, %v)", output, err)
	}
}
