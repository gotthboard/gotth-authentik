# Product requirements

Give unrelated products a consistent, reviewable Authentik OIDC, enrollment,
and access boundary. One manifest must create the OAuth2 provider, application,
enrollment flow, verified-email activation sequence, dedicated access group,
and fail-closed policy binding required by one consumer.

The toolkit must render deterministically, contain no client secret, use exact
redirect URIs, reject unsafe identifiers and URLs, detect sibling applications
sharing an application slug, provider name, client ID, or access group, and
support exact desired-state comparison before an operator applies anything.

The generated blueprint must import cleanly and idempotently into the pinned
Authentik 2026.5 release. Re-import must preserve an Authentik-generated
confidential-client secret.

Non-goals: running Authentik, storing API credentials, owning tenants, managing
users at runtime, or applying remote changes implicitly.

## Alpha.3 admission requirements

- `AUT-A3-01`: New consumers import the documented `pkg/authentik` package.
- `AUT-A3-02`: The CLI imports the sole public package and the module root owns
  no provider, tenant, or Go implementation mechanism.
- `AUT-A3-03`: Layout work cannot accept, render, log, compare, or persist a
  confidential client secret.
- `AUT-A3-04`: Disposable pinned Authentik double-import remains the runtime
  completeness oracle.
- `AUT-A3-05`: Clean-clone, race, canonical-consumer, graph, and two clean Judge
  passes gate alpha.3 admission.
