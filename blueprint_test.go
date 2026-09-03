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
		"authentik_providers_oauth2.oauth2provider", "authentik_core.application",
		"client_type: confidential", `client_id: "example-app"`,
		`url: "https://app.example/auth/callback"`, "matching_mode: strict",
		"grant_types: [authorization_code, refresh_token]", "issuer_mode: per_provider",
		"provider: !KeyOf oidc-provider", "target: !KeyOf application",
		"create_users_as_inactive: true", "activate_user_on_success: true",
		"create_users_group: !KeyOf access-group", "authentik_policies.policybinding",
		"failure_result: false", "[slug, default-provider-authorization-implicit-consent]",
	} {
		if !strings.Contains(string(first), required) {
			t.Errorf("blueprint missing %q", required)
		}
	}
	if strings.Contains(string(first), "client_secret") {
		t.Fatal("secret field rendered into blueprint")
	}
	withoutLaunch := validApplication()
	withoutLaunch.LaunchURL = ""
	output, err := RenderEnrollmentBlueprint(withoutLaunch)
	if err != nil || strings.Contains(string(output), "meta_launch_url") {
		t.Fatalf("render without launch URL = (%q, %v)", output, err)
	}
	if output, err := RenderEnrollmentBlueprint(Application{}); err == nil || output != nil {
		t.Fatalf("invalid render = (%q, %v)", output, err)
	}
}
