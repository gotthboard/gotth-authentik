# Verification

`make verify` checks format, vet, race tests, coverage, generation of the GOTTH
Board example, and byte-for-byte desired-state comparison.

`integration/verify-authentik-2026.5.2.sh` requires preloaded pinned Authentik
and PostgreSQL images plus a usable Docker CLI. Set `DOCKER_SUDO=1` only on a
host whose operator requires passwordless sudo for Docker. The harness creates
an internal network with no external route, a fresh PostgreSQL database, and
fresh Authentik server/worker containers. It imports two isolated fixtures,
then re-imports both. The independent state oracle checks every v1 provider
contract field (type, ID, generated secret presence, shared flows and signing
key, grants, issuer mode, ID-token claim mode, managed scopes, and redirects),
application/provider linkage, issuer route, group and fail-closed policy,
enrollment flow attributes, all prompts and stages, prompt membership, ordered
flow-stage bindings, exact nonduplication, cross-application isolation, and
distinct non-empty unchanged generated secrets. It removes all named container,
network, volume, database, blueprint, and scratch state before successful exit.

The renderer is also compared byte-for-byte with the independently reviewed
`pkg/authentik/testdata/example-blueprint.yaml` golden structure. CLI tests
cover the 1 MiB manifest and blueprint bounds. `make verify` renders and checks
both producer-owned Board and WorldWeave fixtures.
The hash-bound semantic review is retained at
`workflow/features/v1-release/review/golden-structure-clean.md`.

The gate passed against Authentik 2026.5.2 image digest
`sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581`
and PostgreSQL 16 image digest
`sha256:57c72fd2a128e416c7fcc499958864df5301e940bca0a56f58fddf30ffc07777`.
No live Authentik tenant is a valid test target.
