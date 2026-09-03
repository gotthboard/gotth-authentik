# Reusable v0 admission verification

Verified on 2026-09-02/03 with contract commit `75808e2` and implementation
commit `428ff3d7a3be57abf67933c705a8343189835e92`.

## Local gates

- `make verify`
- `go vet -mod=readonly ./...`
- `go test -mod=readonly -race -cover ./...`
- 50 consecutive `go test -mod=readonly -race ./...` runs
- clean-clone `make verify`, including the tracked CLI source and tests
- public library coverage: 99.1% statements
- CLI package coverage: 91.9% statements
- aggregate coverage: 97.2% statements

The only uncovered library branch is the static `text/template` execution
failure after complete input validation. The remaining CLI gap is process-level
`main`/`os.Exit`; `run`, strict manifest decoding, the 1 MiB cap, trailing-data
rejection, output failures, and exact desired-state comparison are fully
covered. The compiled executable is exercised by `make verify`.

## Authentik compatibility gate

`integration/verify-authentik-2026.5.2.sh` passed in a clean archive on the
development validation host with:

- Authentik 2026.5.2
  `sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581`
- PostgreSQL 16
  `sha256:57c72fd2a128e416c7fcc499958864df5301e940bca0a56f58fddf30ffc07777`

The harness used a uniquely named internal Docker network and fresh server,
worker, and database containers. The generated blueprint validated and applied,
created exactly one provider and application, linked them, created the dedicated
group and fail-closed binding, and applied a second time without duplicating
objects or rotating Authentik's generated confidential-client secret. The
result was `AUTHENTIK_IMPORT_OK`. All containers, the network, the database,
temporary source archive, and rendered blueprint were removed afterward.

No live Authentik tenant, live database, credentials, GOTTH Board repository,
or deployed service was read or changed by the integration test.

## Code graph

Graphify 0.9.32 code-only extraction for implementation commit `428ff3d`
produced 60 nodes and 95 directed edges. Graph SHA-256:
`cc38ae91da37151f6abadbb408779d669de309b2c9f7cdf1454343eaeab21e54`.
There were no self-loops, exact duplicate relations, or same-endpoint relation
collisions. Ten documentation files and five unsupported repository metadata
files were skipped; the JSON fixture produced no semantic nodes. The graph was
written outside the repository, and extraction did not change the worktree.
