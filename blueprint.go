package authentik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

var enrollmentTemplate = template.Must(template.New("enrollment").Funcs(template.FuncMap{
	"quote": func(value string) string {
		encoded, _ := json.Marshal(value)
		return string(encoded)
	},
}).Parse(`version: 1
metadata:
  name: {{quote .Name}}
  labels:
    blueprints.goauthentik.io/instantiate: "false"
entries:
  - identifiers: {name: {{.AccessGroup}}}
    id: access-group
    model: authentik_core.group
    attrs: {is_superuser: false}

  - identifiers: {slug: {{.Slug}}-enrollment}
    id: enrollment-flow
    model: authentik_flows.flow
    attrs:
      name: {{quote .Name}}
      title: {{quote .RegistrationTitle}}
      designation: enrollment
      authentication: require_unauthenticated

  - identifiers: {name: {{.Slug}}-username}
    id: username
    model: authentik_stages_prompt.prompt
    attrs: {field_key: username, label: Username, type: username, required: true, order: 0}
  - identifiers: {name: {{.Slug}}-password}
    id: password
    model: authentik_stages_prompt.prompt
    attrs: {field_key: password, label: Password, type: password, required: true, order: 1}
  - identifiers: {name: {{.Slug}}-password-repeat}
    id: password-repeat
    model: authentik_stages_prompt.prompt
    attrs: {field_key: password_repeat, label: Password (repeat), type: password, required: true, order: 2}
  - identifiers: {name: {{.Slug}}-name}
    id: display-name
    model: authentik_stages_prompt.prompt
    attrs: {field_key: name, label: Display name, type: text, required: true, order: 3}
  - identifiers: {name: {{.Slug}}-email}
    id: email
    model: authentik_stages_prompt.prompt
    attrs: {field_key: email, label: Email, type: email, required: true, order: 4}

  - identifiers: {name: {{.Slug}}-prompts}
    id: prompts
    model: authentik_stages_prompt.promptstage
    attrs:
      fields: [!KeyOf username, !KeyOf password, !KeyOf password-repeat, !KeyOf display-name, !KeyOf email]
  - identifiers: {name: {{.Slug}}-user-write}
    id: user-write
    model: authentik_stages_user_write.userwritestage
    attrs:
      create_users_as_inactive: true
      user_creation_mode: always_create
      user_type: external
      user_path_template: {{.UserPath}}
      create_users_group: !KeyOf access-group
  - identifiers: {name: {{.Slug}}-email-verification}
    id: email-verification
    model: authentik_stages_email.emailstage
    attrs:
      use_global_settings: true
      template: email/account_confirmation.html
      activate_user_on_success: true
  - identifiers: {name: {{.Slug}}-login}
    id: login
    model: authentik_stages_user_login.userloginstage

  - identifiers: {target: !KeyOf enrollment-flow, stage: !KeyOf prompts, order: 10}
    model: authentik_flows.flowstagebinding
  - identifiers: {target: !KeyOf enrollment-flow, stage: !KeyOf user-write, order: 20}
    model: authentik_flows.flowstagebinding
  - identifiers: {target: !KeyOf enrollment-flow, stage: !KeyOf email-verification, order: 30}
    model: authentik_flows.flowstagebinding
  - identifiers: {target: !KeyOf enrollment-flow, stage: !KeyOf login, order: 100}
    model: authentik_flows.flowstagebinding

  - identifiers:
      target: !Find [authentik_core.application, [slug, {{.Slug}}]]
      group: !KeyOf access-group
      order: 20
    model: authentik_policies.policybinding
    attrs: {enabled: true, negate: false, failure_result: false}
`))

// RenderEnrollmentBlueprint returns deterministic Authentik blueprint YAML for
// verified-email enrollment and a per-application access policy.
func RenderEnrollmentBlueprint(application Application) ([]byte, error) {
	if err := application.Validate(); err != nil {
		return nil, err
	}
	var output bytes.Buffer
	if err := enrollmentTemplate.Execute(&output, application); err != nil {
		return nil, fmt.Errorf("render Authentik enrollment blueprint: %w", err)
	}
	return output.Bytes(), nil
}
