# Distribution Contract

## Endpoints

- Canonical development source:
  <https://git.dannyhunn.com/gotthboard/gotth-authentik>
- Public clone, Go import, and future releases:
  <https://github.com/gotthboard/gotth-authentik>
- Public bug tracker:
  <https://github.com/gotthboard/gotth-authentik/issues>
- Private vulnerability reports:
  <https://github.com/gotthboard/gotth-authentik/security/advisories/new>

Forgejo pushes one way to GitHub. GitHub does not feed commits or tags back to
Forgejo. A ref is distributed only when the exact object ID is visible at both
endpoints.

## Maturity and compatibility

Current status: admitted local v1 release candidate; no v1 tag or hosted Release.

## Current source use

No post-migration version has been tagged. Until publication, source users must
select an exact reviewed commit; the moving `main` branch is not a compatibility
promise:

```sh
go get github.com/gotthboard/gotth-authentik/pkg/authentik@<reviewed-commit>
go install github.com/gotthboard/gotth-authentik/cmd/gotth-authentik@<reviewed-commit>
```

```go
import authentik "github.com/gotthboard/gotth-authentik/pkg/authentik"
```

A Go command records a pseudo-version in the consumer's `go.mod`; review that
exact revision. The intended v1 stability and support boundary is in
`docs/v1-contract.md` and begins as a published contract only when a v1 tag
exists.

The repository pins Go 1.26.6 where a Go module exists. Supported protocol,
runtime, database, and tool versions remain the ones stated in the README and
project verification documents; this distribution change does not widen those
contracts.

## License

The maintainer selected the MIT license for the owner-authored contents of
this repository. The standard grant is in `LICENSE`; third-party dependencies
and assets retain their own licenses. This decision does not publish a release
or create a support or compatibility promise.

## Migration traceability

| Requirement | Repository implementation | Verification |
| --- | --- | --- |
| DIST-001 | Existing history, tags, worktrees, and mirror direction remain unchanged | pinned ref and worktree inventory |
| DIST-002 | Module directive, exact self-imports, fixtures, and examples use the GitHub identity | stale-prefix search, tidy, vet, test, and clean public import |
| DIST-003/004 | README, contribution, security, changelog, and release contracts describe public use and support | documentation audit |
| DIST-006 | MIT license is present and scoped to owner-authored contents | license and dependency inventory |
| DIST-008 | Forgejo remains source and GitHub remains the one-way mirror target | push-mirror configuration and exact ref comparison |
