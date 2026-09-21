//go:build ignore

package main

import (
	"fmt"
	"os"

	authentik "github.com/gotthboard/gotth-authentik/pkg/authentik"
)

func main() {
	blueprint, err := authentik.RenderDeploymentBlueprint(authentik.Deployment{
		Application: authentik.Application{
			Slug: "integration-example", Name: "Integration Example", RegistrationTitle: "Register for Integration Example",
			UserPath: "users/integration-example", AccessGroup: "integration-example-users", LaunchURL: "https://app.example/",
			Provider: authentik.Provider{
				Name: "integration-example-oidc", ClientType: "confidential", ClientID: "integration-example",
				RedirectURIs:      []string{"https://app.example/auth/callback"},
				AuthorizationFlow: "default-provider-authorization-implicit-consent",
				InvalidationFlow:  "default-provider-invalidation-flow", SigningKey: "authentik Self-signed Certificate",
			},
		},
		OIDCClientSecretFile: "/run/secrets/integration-oidc-client-secret",
		SCIM: authentik.SCIMBackchannel{
			Name:      "Integration Example SCIM",
			URL:       "https://app.example/api/scim/v2/",
			TokenFile: "/run/secrets/integration-scim-token",
		},
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(blueprint); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
