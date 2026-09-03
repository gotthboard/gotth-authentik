# Verification evidence

## Exact state

- Structural implementation: `4a71e45278a710b58e31ff237a86b4d8a8925f6b`.
- Corrected review candidate: `cb87f0ebe27d4f2bc6574bd409fcf4d28b23021a`.
- Base/distribution prerequisite: `0279e0f127f9bb79f95d5bf41e3df1e5066a431e`.
- Canonical package: `github.com/gotthboard/gotth-authentik/pkg/authentik`.

## Coding-setup admission

- Root byte/inode preflight: 5% bytes, 1% inodes; below both stop thresholds.
- Context broker 0.1.0: clean revision, cache miss, untruncated bounded packet;
  cache path `/home/linus/.cache/openclaw-code-context/2ae8e717de2dbf67/9e68f6dea599ab3d/0fe39734a6e7d12b8b67f85a620401f55d1787cabe55e85e22d49bc2560a92b9.json`.
- Production library units were not changed; files are 100% identical renames.
  The CLI changed only its import. Prospective complexity comments are N/A.
- Performance admission: N/A. Validation, rendering, importer/container and
  database mechanisms are unchanged; no speedup is claimed.
- Runtime contract: Go 1.26.6, Authentik 2026.5.2 and PostgreSQL 16 pinned by
  digest, disposable isolated network/database, double import, exact linkage,
  no duplicate objects, unchanged generated secret, cleanup on every exit.
- `gopls` was unavailable and was not installed; compiler, vet, tests, pinned
  importer integration, and outside-package compilation are authoritative.

## Verification

- `go mod verify && make verify`: PASS; library coverage 99.1%, CLI 91.9%,
  aggregate 97.2%.
- Fifty consecutive `go test -mod=readonly -race ./...` runs: PASS.
- Deterministic render/check: PASS.
- Pinned Authentik 2026.5.2/PostgreSQL 16 double import: `AUTHENTIK_IMPORT_OK`;
  containers, network, database, files, and checkout state removed.
- Module root contains zero Go files; CLI/canonical import identity: PASS.
- Two independent cold Judge passes on one exact committed state: CLEAN.
- No live tenant, credential, secret, database, tag, or deployment changed.

## Graph evidence

Graphify 0.9.32, code-only, implementation revision
`4a71e45278a710b58e31ff237a86b4d8a8925f6b`:

- path: `/home/linus/.cache/openclaw-code-index/gotth-authentik/4a71e45278a710b58e31ff237a86b4d8a8925f6b/graphify/graphify-out/graph.json`
- SHA-256: `f2802d9ef4fadfdb18dfa76a3570145c9ac1df6a1753b20617f44358bb5b6d3c`
- 63 nodes, 99 edges, 10 communities; zero self-loops, duplicates,
  same-endpoint collisions, or dangling endpoints.
- Limitation: fixture JSON produced no graph nodes. The actual generated YAML
  was verified by the pinned Authentik importer rather than guessed from graph
  absence.

Graph output remained navigation evidence, not an admission oracle.
