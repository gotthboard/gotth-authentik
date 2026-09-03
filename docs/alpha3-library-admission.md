# Alpha.3 library admission

## Scope and authority

This pass admits `pkg/authentik` as the canonical package for GOTTH Board
alpha.3. It does not apply blueprints, own a tenant, handle API credentials, or
accept confidential client secrets.

## Requirement traceability

| Requirement | Design/specification | Code | Verification |
|---|---|---|---|
| `AUT-A3-01` | architecture and README layout | `pkg/authentik/` | canonical outside-package test |
| `AUT-A3-02` | implementation specification | `pkg/authentik/` and `cmd/` | build and import inspection |
| `AUT-A3-03` | PRD/security contract | manifest and renderer | secret-absence tests/import evidence |
| `AUT-A3-04` | verification contract | integration harness | pinned double-import completeness oracle |
| `AUT-A3-05` | workflow admission | tests/evidence | clean clone, graph, two Judge passes |

## Runtime boundary

- Go 1.26.6, Authentik 2026.5.2 at the recorded digest, and PostgreSQL 16 at
  the recorded digest.
- Authentik importer behavior, object linkage, idempotent re-import, isolated
  access binding, and generated-secret preservation are correctness boundaries.
- The completeness oracle is a fresh disposable tenant imported twice with
  exact object/linkage counts and unchanged non-empty confidential secret.
- Shared tenant flows and signing keys are explicit consumer references; no
  live tenant is a valid admission target.

## Performance admission

No renderer, validation, importer, container, or database mechanism changes.
The canonical package is the original implementation and the root owns
governance only. No speedup is claimed; benchmark/Amdahl evidence is N/A for
this structural admission.

## Failure and rollback

Rollback is a revert before the first consumer pin. Disposable integration
state must be destroyed on every exit. No live Authentik object or secret is
created, read, changed, logged, or copied.
