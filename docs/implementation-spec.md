# Implementation specification

- Canonical package:
  `github.com/gotthboard/gotth-authentik/pkg/authentik`.
- The module root contains no Go package; the CLI depends only on the canonical
  public package.

- One manifest per application.
- Slugs, access groups, and user-path segments are bounded lowercase slugs.
- Display strings are valid bounded UTF-8 without control characters.
- Provider name and client ID are explicit, bounded, non-secret identifiers.
- Provider client type is exactly `confidential`. Public clients are outside the
  v1 contract because no pinned importer evidence admits them.
- Redirect URIs are unique, canonical, and use HTTPS, except numeric loopback
  HTTP for disposable local consumers.
- Launch and redirect URLs are bounded to 2048 bytes. Canonical authorities use
  lowercase DNS labels or canonical IP text, omit default/empty/zero-padded
  ports, and reject encoded or dot-segment paths.
- Authorization/invalidation flow slugs and the signing-key name are explicit.
- OAuth2 provider desired state uses strict redirect matching, authorization
  code plus refresh-token grants, per-provider issuer mode, and the built-in
  `openid`, `profile`, and `email` mappings.
- Production manifest, renderer, and CLI code never accepts, renders, logs, or
  compares the client secret. Authentik generates it on first
  confidential-provider creation. The disposable oracle reads and compares
  generated secrets only in memory and never prints them.
- The Authentik application is created in the same blueprint and refers to the
  provider with `!KeyOf`.
- New users are inactive until the email stage succeeds.
- The creation stage adds users only to that application's access group.
- The application policy binding is enabled, non-negated, and fail-closed.
- Blueprint rendering is deterministic and secret-free.
- `render` writes desired YAML; `check` compares it with an existing file.
  Both CLI inputs are bounded to 1 MiB.
- Issuer validation requires exactly `/application/o/<application-slug>/`.
- Remote apply is intentionally absent from the extraction baseline.
- Disposable import verification creates a fresh tenant, imports twice, and
  proves the confidential secret is non-empty and unchanged across re-import.
