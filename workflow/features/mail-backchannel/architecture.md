# Architecture: deployment overlay

## Boundary

The existing enrollment blueprint remains the owner of the OIDC provider,
application, enrollment flow, access group, and policy binding. The new API
renders that exact blueprint followed by a narrow overlay. Authentik's importer
resolves repeated model identifiers as updates, so the overlay can supply the
OIDC `client_secret` without duplicating ownership or copying the base
renderer.

The overlay contains exactly three entries:

1. an OAuth2 provider update keyed by the existing provider name, with
   `client_secret: !File <path>`;
2. one SCIM provider keyed by its name, with `auth_mode: token`,
   `token: !File <path>`, certificate verification, the two Authentik-managed
   SCIM mappings, and `group_filters` containing only the existing access
   group; and
3. an application update keyed by the existing slug, attaching the SCIM
   provider through `backchannel_providers`.

## Data and trust flow

The caller supplies non-secret identifiers, URLs, and secret-file references.
Validation rejects ambiguous input. Rendering serializes only those validated
references. Authentik opens the files when its importer processes `!File`.
Consequently secret bytes do not cross the Go process and cannot leak through
its heap, output, errors, logs, tests, or evidence.

The container or host supervisor owns file creation, permissions, rotation,
and mount policy. The library owns none of that state. It also owns no network
client and has no apply method.

## Failure behavior

Invalid applications, non-confidential OIDC providers, malformed names,
noncanonical or insecure endpoints, and unsafe file references fail before any
output is returned. Missing or unreadable files fail later inside Authentik at
explicit import time. This separation is intentional: checking file existence
inside the renderer would couple artifact generation to deployment state and
would pull credentials into the wrong trust boundary.

## Compatibility

`RenderEnrollmentBlueprint` and `Application` remain untouched. The additive
types and function are pre-1.0 API. The generated overlay is pinned to the
Authen­tik 2026.5 importer contract and must be revalidated before claiming
compatibility with another Authentik release.
