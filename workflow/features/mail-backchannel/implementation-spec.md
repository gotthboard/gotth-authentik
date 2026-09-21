# Implementation specification: deployment overlay

## Public API

Add these bounded public types to `pkg/authentik`:

```go
type Deployment struct {
    Application          Application
    OIDCClientSecretFile string
    SCIM                 SCIMBackchannel
}

type SCIMBackchannel struct {
    Name      string
    URL       string
    TokenFile string
}

func (Deployment) Validate() error
func (SCIMBackchannel) Validate() error
func RenderDeploymentBlueprint(Deployment) ([]byte, error)
```

`Deployment.Validate` requires a valid application, a confidential OIDC
provider, two safe credential references, and a valid SCIM backchannel.
`SCIMBackchannel.Name` uses bounded control-free display text because
Authentik's provider identifier is a name rather than a slug.

## Credential reference validation

A credential reference is valid only when it:

- is valid UTF-8 and at most 255 bytes;
- begins with `/run/secrets/`;
- is absolute and already equal to `filepath.Clean(value)`;
- names a child rather than `/run/secrets` itself; and
- contains no NUL or control character.

The renderer must not call `Open`, `ReadFile`, `Stat`, or any filesystem API on
the reference.

## Rendering

Render the existing enrollment bytes first. Append one deterministic template
whose YAML nodes use JSON string quoting for all caller data. `!File` receives
a quoted path scalar. Managed mapping identifiers are constants:

- `goauthentik.io/providers/scim/user`
- `goauthentik.io/providers/scim/group`

Fixed fields are `auth_mode: token`, `verify_certificates: true`,
`exclude_users_service_account: true`, and `compatibility_mode: default`.
The access group comes only from `Deployment.Application.AccessGroup`.

## Tests

Unit tests shall prove deterministic output, exact required fields and tags,
absence of secret values from the API surface, legacy renderer byte stability,
and rejection of non-confidential OIDC, unsafe paths, insecure/noncanonical
URLs, bad names, and invalid applications. Public-package compilation shall
cover all new exported types and functions.

The disposable integration shall mount two test credential files into both
server and worker containers, import twice, and query Authentik models to prove:

- the OAuth2 provider secret equals the file content;
- the SCIM provider token equals the file content;
- auth mode, certificate verification, compatibility mode, mappings, and
  access-group filter are exact;
- the application has exactly that SCIM backchannel;
- provider and application counts remain one after re-import.

Synthetic integration values must never be presented as live credentials.
