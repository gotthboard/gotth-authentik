package authentik

import (
	"bytes"
	"strings"
	"testing"
)

func validDeployment() Deployment {
	return Deployment{
		Application:          validApplication(),
		OIDCClientSecretFile: "/run/secrets/example-oidc-client-secret",
		SCIM: SCIMBackchannel{
			Name:      "Example application SCIM",
			URL:       "https://app.example/api/scim/v2/",
			TokenFile: "/run/secrets/example-scim-token",
		},
	}
}

func TestRenderDeploymentBlueprint(t *testing.T) {
	deployment := validDeployment()
	first, err := RenderDeploymentBlueprint(deployment)
	if err != nil {
		t.Fatalf("RenderDeploymentBlueprint() returned error: %v", err)
	}
	second, err := RenderDeploymentBlueprint(deployment)
	if err != nil || !bytes.Equal(first, second) {
		t.Fatal("deployment blueprint rendering is not deterministic")
	}
	base, err := RenderEnrollmentBlueprint(deployment.Application)
	if err != nil || !bytes.HasPrefix(first, base) {
		t.Fatal("deployment blueprint does not preserve the enrollment blueprint")
	}
	for _, required := range []string{
		`client_secret: !File "/run/secrets/example-oidc-client-secret"`,
		"model: authentik_providers_scim.scimprovider",
		`url: "https://app.example/api/scim/v2/"`,
		`token: !File "/run/secrets/example-scim-token"`,
		"auth_mode: token", "verify_certificates: true",
		"exclude_users_service_account: true", "compatibility_mode: default",
		"goauthentik.io/providers/scim/user", "goauthentik.io/providers/scim/group",
		"group_filters: [!KeyOf access-group]",
		"backchannel_providers: [!KeyOf scim-provider]",
	} {
		if !strings.Contains(string(first), required) {
			t.Errorf("deployment blueprint missing %q", required)
		}
	}
	if strings.Count(string(first), "model: authentik_providers_scim.scimprovider") != 1 {
		t.Fatal("deployment blueprint does not contain exactly one SCIM provider")
	}
}

func TestDeploymentValidationRejectsUnsafeState(t *testing.T) {
	valid := validDeployment()
	tests := []struct {
		name   string
		mutate func(*Deployment)
	}{
		{"invalid application", func(value *Deployment) { value.Application = Application{} }},
		{"public OIDC provider", func(value *Deployment) { value.Application.Provider.ClientType = "public" }},
		{"empty OIDC path", func(value *Deployment) { value.OIDCClientSecretFile = "" }},
		{"relative OIDC path", func(value *Deployment) { value.OIDCClientSecretFile = "run/secrets/client" }},
		{"parent OIDC path", func(value *Deployment) { value.OIDCClientSecretFile = "/run/secrets/../client" }},
		{"unclean OIDC path", func(value *Deployment) { value.OIDCClientSecretFile = "/run/secrets/a/../client" }},
		{"control OIDC path", func(value *Deployment) { value.OIDCClientSecretFile = "/run/secrets/client\nname" }},
		{"overlong OIDC path", func(value *Deployment) { value.OIDCClientSecretFile = "/run/secrets/" + strings.Repeat("a", 244) }},
		{"empty SCIM name", func(value *Deployment) { value.SCIM.Name = "" }},
		{"control SCIM name", func(value *Deployment) { value.SCIM.Name = "bad\nname" }},
		{"insecure SCIM URL", func(value *Deployment) { value.SCIM.URL = "http://app.example/scim/" }},
		{"SCIM URL credentials", func(value *Deployment) { value.SCIM.URL = "https://user@app.example/scim/" }},
		{"SCIM URL query", func(value *Deployment) { value.SCIM.URL = "https://app.example/scim/?x=1" }},
		{"unsafe SCIM token path", func(value *Deployment) { value.SCIM.TokenFile = "/tmp/token" }},
		{"shared credential path", func(value *Deployment) { value.SCIM.TokenFile = value.OIDCClientSecretFile }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := valid
			test.mutate(&candidate)
			if err := candidate.Validate(); err == nil {
				t.Fatal("invalid deployment passed validation")
			}
			if output, err := RenderDeploymentBlueprint(candidate); err == nil || output != nil {
				t.Fatalf("invalid render = (%q, %v)", output, err)
			}
		})
	}
}

func TestDeploymentAcceptsNumericLoopbackForDisposableIntegration(t *testing.T) {
	deployment := validDeployment()
	deployment.SCIM.URL = "http://127.0.0.1:8080/scim/v2/"
	if err := deployment.Validate(); err != nil {
		t.Fatalf("numeric loopback deployment rejected: %v", err)
	}
	deployment.SCIM.URL = "http://localhost:8080/scim/v2/"
	if err := deployment.Validate(); err == nil {
		t.Fatal("hostname loopback deployment passed")
	}
}
