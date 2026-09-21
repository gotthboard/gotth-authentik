package authentik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

const secretDirectory = "/run/secrets/"

// Deployment is the desired Authentik state for an application whose OIDC
// credential and SCIM bearer token are supplied through mounted files.
// Credential values are deliberately absent from this type.
type Deployment struct {
	Application          Application
	OIDCClientSecretFile string
	SCIM                 SCIMBackchannel
}

// SCIMBackchannel describes the one SCIM provider attached to a deployment.
type SCIMBackchannel struct {
	Name      string
	URL       string
	TokenFile string
}

// Validate rejects deployment state that could broaden synchronization or put
// credential material into the desired-state boundary.
func (deployment Deployment) Validate() error {
	if err := deployment.Application.Validate(); err != nil {
		return fmt.Errorf("application: %w", err)
	}
	if deployment.Application.Provider.ClientType != "confidential" {
		return fmt.Errorf("OIDC provider must be confidential")
	}
	if err := validateSecretFile(deployment.OIDCClientSecretFile); err != nil {
		return fmt.Errorf("OIDC client secret file: %w", err)
	}
	if err := deployment.SCIM.Validate(); err != nil {
		return fmt.Errorf("SCIM backchannel: %w", err)
	}
	return nil
}

// Validate rejects ambiguous SCIM provider names, insecure endpoints, and
// credential references outside the fixed runtime secret directory.
func (backchannel SCIMBackchannel) Validate() error {
	if !validText(backchannel.Name, 1, 255) {
		return fmt.Errorf("name is invalid")
	}
	if _, err := validateSecureURL(backchannel.URL); err != nil {
		return fmt.Errorf("URL is invalid: %w", err)
	}
	if err := validateSecretFile(backchannel.TokenFile); err != nil {
		return fmt.Errorf("token file: %w", err)
	}
	return nil
}

func validateSecretFile(value string) error {
	if value == "" || len(value) > 255 || !utf8.ValidString(value) {
		return fmt.Errorf("path is empty, overlong, or invalid UTF-8")
	}
	if strings.IndexFunc(value, unicode.IsControl) >= 0 {
		return fmt.Errorf("path contains a control character")
	}
	if !strings.HasPrefix(value, secretDirectory) || path.Clean(value) != value {
		return fmt.Errorf("path must be canonical and beneath %s", strings.TrimSuffix(secretDirectory, "/"))
	}
	return nil
}

var deploymentOverlayTemplate = template.Must(template.New("deployment-overlay").Funcs(template.FuncMap{
	"quote": func(value string) string {
		encoded, _ := json.Marshal(value)
		return string(encoded)
	},
}).Parse(`
  - identifiers: {name: {{quote .Application.Provider.Name}}}
    model: authentik_providers_oauth2.oauth2provider
    attrs:
      client_secret: !File {{quote .OIDCClientSecretFile}}

  - identifiers: {name: {{quote .SCIM.Name}}}
    id: scim-provider
    model: authentik_providers_scim.scimprovider
    attrs:
      url: {{quote .SCIM.URL}}
      token: !File {{quote .SCIM.TokenFile}}
      auth_mode: token
      verify_certificates: true
      exclude_users_service_account: true
      compatibility_mode: default
      property_mappings:
        - !Find [authentik_providers_scim.scimmapping, [managed, goauthentik.io/providers/scim/user]]
      property_mappings_group:
        - !Find [authentik_providers_scim.scimmapping, [managed, goauthentik.io/providers/scim/group]]
      group_filters: [!KeyOf access-group]

  - identifiers: {slug: {{.Application.Slug}}}
    model: authentik_core.application
    attrs:
      backchannel_providers: [!KeyOf scim-provider]
`))

// RenderDeploymentBlueprint returns the existing enrollment blueprint plus a
// deterministic overlay whose credential values are resolved by Authentik
// through !File references at explicit import time.
func RenderDeploymentBlueprint(deployment Deployment) ([]byte, error) {
	if err := deployment.Validate(); err != nil {
		return nil, err
	}
	base, err := RenderEnrollmentBlueprint(deployment.Application)
	if err != nil {
		return nil, err
	}
	var output bytes.Buffer
	_, _ = output.Write(base)
	if err := deploymentOverlayTemplate.Execute(&output, deployment); err != nil {
		return nil, fmt.Errorf("render Authentik deployment blueprint: %w", err)
	}
	return output.Bytes(), nil
}
