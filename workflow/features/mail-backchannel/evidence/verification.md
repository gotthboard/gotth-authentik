# Verification evidence

Feature source commit: `f47749261d91433b677af5b8efcabf8f6c0861f8`

Environment: `development`, Go 1.26.6, pinned Authentik 2026.5.2 image
`ghcr.io/goauthentik/server@sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581`.

## Results

- `make verify`: PASS; format, vet, race, package tests, CLI render/check.
- `go test -mod=readonly -race -run 'Deployment|SCIM' -count=100 ./pkg/authentik`: PASS.
- `go test -mod=readonly -shuffle=on -count=20 ./...`: PASS.
- `go test -mod=readonly -coverprofile=... ./...`: PASS; public library
  97.9%, CLI 91.9%, repository total 96.6% statement coverage.
- `sh -n integration/verify-authentik-2026.5.2.sh`: PASS.
- `shellcheck integration/verify-authentik-2026.5.2.sh`: PASS when available
  on the verification host.
- clean `--no-local` clone at the exact source commit followed by
  `make verify`: PASS.
- `DOCKER_SUDO=1 ./integration/verify-authentik-2026.5.2.sh`: PASS with
  `AUTHENTIK_IMPORT_OK`.

## Disposable importer proof

The integration gate created an internal-only Docker network, PostgreSQL,
Authentik server, and worker, and removed all of them on exit. It proved:

- a missing OIDC credential file fails with an Authentik `EntryInvalidError`
  naming `client_secret` and leaves no partial OAuth provider;
- the rendered blueprint contains neither synthetic credential value;
- Authentik's `!File` loader resolves the mounted OIDC client secret and SCIM
  bearer token exactly;
- SCIM token authentication, certificate verification, default compatibility,
  service-account exclusion, managed user/group mappings, one access-group
  filter, and one application backchannel are exact;
- a second import preserves both credentials and all relations without
  duplicate provider or application rows.

The first disposable attempt failed because synthetic host files were mode
0600 and owned by host UID 1002 while Authentik runs as UID 1000. The fixture
was corrected to Docker-secret-style read-only 0444 mounts. This was a test
fixture ownership error, not an application or blueprint relaxation.

No live tenant, credential, DNS record, release, or deployment was changed.
