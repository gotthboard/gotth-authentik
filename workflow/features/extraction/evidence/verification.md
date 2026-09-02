# Extraction verification

Verified on 2026-09-02 with Go `go1.26.6-X:nodwarf5`:

- `go vet -mod=readonly ./...`
- `go test -mod=readonly -race -cover ./...`
- aggregate statement coverage: 89.2%
- deterministic GOTTH Board manifest rendering succeeds
- exact issuer/enrollment URL rejection and sibling access-group isolation pass
- accepted alpha.2 Authentik 2026.5.2 blueprint retained unchanged as provenance

No Authentik tenant was modified. The newly generated generic blueprint has not
yet been imported into a disposable Authentik tenant; that is a recorded,
blocking release gate before GOTTH Board or another consumer pins this module.

Graphify 0.9.32 code-only audit: 27 nodes, 37 directed post-build edges, no
self-loops, exact duplicate edges, or same-endpoint relation groups.
