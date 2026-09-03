# Implementation specification

- Canonical package:
  `github.com/gotthboard/gotth-authentik/pkg/authentik`.
- The module root contains no Go package; the CLI depends only on the canonical
  public package.

- One manifest per application.
- Slugs, access groups, and user-path segments are bounded lowercase slugs.
- Display strings are valid bounded UTF-8 without control characters.
- Provider name and client ID are explicit, bounded, non-secret identifiers.
- Provider client type is exactly `confidential` or `public`.
- Redirect URIs are unique, canonical, and use HTTPS, except numeric loopback
  HTTP for disposable local consumers.
- Authorization/invalidation flow slugs and the signing-key name are explicit.
- OAuth2 provider desired state uses strict redirect matching, authorization
  code plus refresh-token grants, per-provider issuer mode, and the built-in
  `openid`, `profile`, and `email` mappings.
- The client secret is never accepted, rendered, logged, or compared. Authentik
  generates it on first confidential-provider creation.
- The Authentik application is created in the same blueprint and refers to the
  provider with `!KeyOf`.
- New users are inactive until the email stage succeeds.
- The creation stage adds users only to that application's access group.
- The application policy binding is enabled, non-negated, and fail-closed.
- Blueprint rendering is deterministic and secret-free.
- `render` writes desired YAML; `check` compares it with an existing file.
- Remote apply is intentionally absent from the extraction baseline.
- Disposable import verification creates a fresh tenant, imports twice, and
  proves the confidential secret is non-empty and unchanged across re-import.
