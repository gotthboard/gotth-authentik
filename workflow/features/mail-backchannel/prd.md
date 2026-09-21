# Product requirements: file-backed OIDC and SCIM backchannel

## Problem

The admitted renderer creates an isolated OIDC application but deliberately
leaves its confidential client secret generated inside Authentik. A deployed
service needs a repeatable, auditable way to supply that credential and a SCIM
bearer token without placing either value in a manifest, environment variable,
command line, generated blueprint, log, or evidence artifact. The application
also needs Authentik's SCIM provider attached as a backchannel so provisioning
is restricted by the existing application access group.

## Required behavior

`gotth-authentik` shall render a deterministic deployment blueprint from:

1. one already-valid `Application`;
2. an absolute file reference beneath `/run/secrets/` for the confidential
   OIDC client secret; and
3. one bounded SCIM provider name, canonical endpoint URL, and absolute token
   file reference beneath `/run/secrets/`.

The blueprint shall use Authentik 2026.5's documented `!File` tag. It shall
patch the OIDC provider identified by the application's provider name, create
the SCIM provider with token authentication and certificate verification,
select Authentik's managed user and group mappings, filter synchronization to
the application's access group, and attach that provider to the application's
`backchannel_providers` relation.

## Security and userspace requirements

- The public API accepts file paths, never secret bytes.
- The renderer never opens or stats credential files.
- Paths are absolute, clean, non-root paths strictly below `/run/secrets/`.
- OIDC and SCIM credentials use distinct file references and rotation domains.
- Production SCIM URLs require HTTPS. Numeric loopback HTTP is accepted only
  to support a disposable integration environment.
- The existing `RenderEnrollmentBlueprint` API and its generated-secret
  semantics remain byte-for-byte unchanged.
- Repeated imports converge without duplicate providers or applications and
  without changing supplied credential values.
- Rendering remains a local Class 2 artifact operation. Applying the artifact
  to any tenant remains a separate, explicit Class 3 operation.

## Exclusions

This feature does not provide arbitrary Authentik object rendering, alternate
SCIM authentication modes, custom property mappings, multiple groups, role
mapping, live tenant mutation, secret rotation orchestration, or release
publication.

## Admission

Admission requires unit and negative tests, unchanged legacy output, full Go
format/vet/race/coverage gates, and a disposable Authentik 2026.5.2 double
import proving the exact database state. No claim of live readiness is allowed
without a public hostname, real credentials, and a separately approved apply.
