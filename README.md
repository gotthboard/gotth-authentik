# gotth-authentik

> **Distribution:** GitHub is the public clone, Go import, and future release endpoint.
> Forgejo remains canonical development. Report public bugs through GitHub
> Issues and security vulnerabilities through GitHub private reporting.
> See [the distribution contract](docs/distribution.md).


`gotth-authentik` is a versioned desired-state and validation toolkit for
isolated Authentik OIDC applications. One secret-free manifest generates the
OAuth2 provider, application, verified-email enrollment flow,
inactive-until-verified user creation, dedicated access group, and fail-closed
application policy binding.

Canonical Go package:
`github.com/gotthboard/gotth-authentik/pkg/authentik`. The module root
contains repository governance only; the CLI imports the canonical package.

Repository layout:

- `pkg/authentik/` — public manifest validation and blueprint implementation;
- `cmd/gotth-authentik/` — bounded render/check CLI;
- module root — module metadata and repository governance;
- `integration/`, `examples/`, `docs/`, and `workflow/` — disposable provider
  evidence, provenance, contracts, and canonical workflow state.

It is not a daemon and does not silently apply remote changes. Rendering and
inspection are local; a future remote apply command must remain an explicit
operator action.

Provider client IDs, exact redirect URIs, shared authorization/invalidation
flow slugs, and the signing-key name are explicit. Confidential client secrets
are not accepted or rendered; Authentik generates one when the provider is
first created and preserves it on re-import. The v1 contract supports
confidential clients and Authentik 2026.5.2 only.

```text
gotth-authentik render manifest.json > application.yaml
gotth-authentik check manifest.json application.yaml
```

`check` is a local byte-for-byte desired-state comparison. Applying a blueprint
is always a separate, explicit operator action. The retained integration gate
uses a fresh internal Docker network, PostgreSQL database, Authentik server,
and worker; it imports twice and destroys all disposable state afterward.

The accepted GOTTH Board alpha.2 blueprint remains provenance. Producer-owned
Board and WorldWeave fixtures exercise separate profiles; the generic
two-application oracle creates complete isolated boundaries and has been
admitted by the Authentik 2026.5.2 importer.

## Installation, compatibility, and support

The source is a prepared v1 release candidate pending final admission; it
remains unpublished, with no v1 tag or hosted Release. The intended stable API, manifest, CLI,
runtime, and support surface is in [the v1 contract](docs/v1-contract.md).

Before the first tag, consumers must use an exact reviewed commit rather than a
moving branch. From a pinned checkout, install the CLI with:

```sh
go install ./cmd/gotth-authentik
```

The owner-authored repository contents are licensed under the
[MIT License](LICENSE). Third-party dependencies and assets retain their own
licenses. No long-term support promise is implied. Versioning, release
admission, security reporting, and contribution details are in
[the release policy](docs/RELEASING.md), [security policy](SECURITY.md), and
[contribution guide](CONTRIBUTING.md). Revision-matched and fixture-based
consumer evidence is separated in
[the v1 consumer admission](docs/v1-consumer-admission.md).
