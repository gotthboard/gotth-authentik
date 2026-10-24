# Version 1 contract

This document freezes the admitted source contract intended for the first v1
release. No tag, hosted Release, or deployment is created by this work.

## Supported runtime

- Go toolchain: Go 1.26.6.
- Authentik: exactly 2026.5.2. Newer and older Authentik releases are outside
  the v1 compatibility claim until separately admitted.
- Provider type: confidential OAuth2/OIDC clients only.
- Execution model: local validation, rendering, and comparison. There is no
  daemon, Authentik API client, credential input, or remote apply command.

The toolkit accepts no client secret. Authentik generates the confidential
client secret on initial import, and an unchanged blueprint re-import preserves
it.

## Go API

The public package is

github.com/gotthboard/gotth-authentik/pkg/authentik

The v1 surface is:

- Application and Provider, including their documented JSON field names;
- Application.Validate and Provider.Validate;
- ValidateIsolation;
- ValidateIssuerURL;
- ValidateEnrollmentURL;
- RenderEnrollmentBlueprint.

ValidateIssuerURL accepts only a canonical secure URL whose path is exactly
/application/o/<application-slug>/. ValidateEnrollmentURL additionally binds an
exact origin and enrollment-flow slug. Launch and redirect URLs accept numeric
loopback HTTP for local consumers; issuer and enrollment URLs require HTTPS.
URL authorities use lowercase DNS labels or canonical IP text, canonical
non-default ports, and unencoded paths without dot segments.

## Manifest v1

The v1 JSON manifest is exactly the JSON representation of Application:

- slug, name, registration_title, user_path, access_group;
- optional launch_url;
- provider.name, provider.client_type, provider.client_id;
- provider.redirect_uris;
- provider.authorization_flow, provider.invalidation_flow, and
  provider.signing_key.

client_type must be confidential. Launch URLs and each redirect URI are bounded
to 2048 bytes. The CLI rejects unknown fields, trailing JSON,
invalid values, files larger than 1 MiB, and manifests with anything other than
one JSON value. There is no secret field and no schema extension mechanism in
v1.

## Blueprint contract

Rendering is deterministic and secret-free. It creates exactly one provider,
application, access group, enrollment flow, five prompts, prompt/user-write/
email/login stages, four ordered flow-stage bindings, and one fail-closed group
policy binding. Shared authorization and invalidation flows, the signing key,
and managed openid/profile/email scope mappings are referenced rather than
owned.

The provider uses strict redirect matching, authorization_code and refresh_token
grants, claims in the ID token, and per-provider issuer mode. The enrollment
flow creates external users inactive in the application-specific group and
activates them only after the email stage succeeds.

## CLI contract

The complete v1 command surface is:

    gotth-authentik render <manifest.json>
    gotth-authentik check <manifest.json> <blueprint.yaml>

render writes the blueprint to standard output. check performs a byte-for-byte
comparison with locally rendered desired state and writes no success output.
Both inputs are bounded to 1 MiB. Invalid usage, invalid input, read/write
failure, or mismatch returns a non-zero process status. No command applies a
blueprint or contacts Authentik.

## Compatibility and support

Within v1, backwards-incompatible changes to the public Go API, JSON field
meaning, rendered blueprint structure, or CLI command/exit contract require a
new major version. Additive changes still require tests, documentation, and a
fresh Authentik 2026.5.2 importer run because importer behavior is part of the
contract.

Before a v1 tag exists, consumers must pin an exact commit and receive no tagged
compatibility or security-support promise. After publication, the latest v1
patch is the supported v1 line; fixes are not promised for superseded v1 patch
releases. Authentik versions other than 2026.5.2, public clients, remote apply,
and live-tenant operation remain unsupported.
