# gotth-authentik

`gotth-authentik` is a versioned desired-state and validation toolkit for
isolated Authentik OIDC applications. One secret-free manifest generates the
OAuth2 provider, application, verified-email enrollment flow,
inactive-until-verified user creation, dedicated access group, and fail-closed
application policy binding.

It is not a daemon and does not silently apply remote changes. Rendering and
inspection are local; a future remote apply command must remain an explicit
operator action.

Provider client IDs, exact redirect URIs, shared authorization/invalidation
flow slugs, and the signing-key name are explicit. Confidential client secrets
are not accepted or rendered; Authentik generates one when the provider is
first created and preserves it on re-import.

```text
gotth-authentik render manifest.json > application.yaml
gotth-authentik check manifest.json application.yaml
```

`check` is a local byte-for-byte desired-state comparison. Applying a blueprint
is always a separate, explicit operator action. The retained integration gate
uses a fresh internal Docker network, PostgreSQL database, Authentik server,
and worker; it imports twice and destroys all disposable state afterward.

The accepted GOTTH Board alpha.2 blueprint remains provenance. The generic
example creates a complete application boundary and has been admitted by the
Authentik 2026.5.2 importer.
