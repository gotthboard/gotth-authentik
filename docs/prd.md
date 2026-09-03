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
