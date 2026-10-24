# v1 release-candidate verification

Date: 2026-10-24

## Exact source boundary

- canonical base and worktree HEAD before the local release commit:
  77c811c18fe3c107f6f3e97ce3f9a3850a685300
- branch: workflow/feature/v1-release
- worktree:
  /home/linus/.openclaw/workspace/.artifacts/worktrees/gotth-authentik-v1-release
- runtime-input tree SHA-256:
  99f199e21249790c5cc16ab7df61b68d5851c8d39a60a561afec056239736a06

The runtime-input hash covers go.mod, .go-version, Makefile, and every file
under cmd/, pkg/, integration/, and examples/ in sorted path/content order.
Documentation and this evidence record are intentionally outside that hash.
The final local commit containing this record is reported in the handoff; no
remote ref, tag, hosted Release, or deployment was created.

Exact hash recipe, executed from the candidate root:

    LC_ALL=C; export LC_ALL
    { printf '%s\0' go.mod .go-version Makefile; find cmd pkg integration examples -type f -print0; } |
      sort -z | xargs -0 sha256sum | sha256sum

Retained independent inputs:

- pkg/authentik/testdata/example-blueprint.yaml:
  8b60022d06315f3e8348654c15c588a451bd67f37e8cbc3cddd44b7a125b1444
- integration/assert-authentik-state.py:
  0a7abc6a37f91915e4fe51ccf54142c8f11c1f5ef6ed40dc59a59cc05c3f055d
- integration/fixture-manifest.json:
  85d73987da84cb6d8de99d1ca15d579faa568e7df9d960f758a9df53f8ac3481
- integration/fixture-sibling-manifest.json:
  ffaa5d66a65202583ec9f2bb06160c44ee4a6c9df417101dcc55e8aa7c403b12
- integration/verify-authentik-2026.5.2.sh:
  880bfee0abd522909ca83d0a40a48edb753842ad431ef58395a2e72a899b633c

## Go and CLI gates

Executed on development, not in the Gateway cgroup:

- go version: go1.26.6 linux/amd64
- go mod verify: PASS, all modules verified
- make verify: PASS
- gofmt inventory: PASS
- go vet -mod=readonly ./...: PASS
- go test -mod=readonly -race -cover ./...: PASS
- cmd coverage: 92.7% statements
- pkg/authentik coverage: 98.8% statements
- Board and WorldWeave render/check fixtures: PASS
- staged-tree diff check in the disposable copy: PASS

Result marker: REMOTE_GO_GATES_OK. The reported runtime-input hash matched the
runtime gate below.

Focused local checks also passed with GOTOOLCHAIN=go1.26.6 and GOMAXPROCS=2.
The final candidate CLI checked the unchanged Mail manifest and blueprint
byte-for-byte: MAIL_CANDIDATE_CHECK_OK.

## Disposable Authentik runtime gate

Executed only on development with a fresh internal Docker network, PostgreSQL
database, Authentik server, and worker. The harness enforces these exact image
IDs and accepts no image override:

- Authentik 2026.5.2:
  sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581
- PostgreSQL 16:
  sha256:57c72fd2a128e416c7fcc499958864df5301e940bca0a56f58fddf30ffc07777

The oracle validated and applied two confidential-client blueprints, then
validated and applied both again. Database assertions covered exact provider
fields and shared references, strict redirects, grants, scopes, issuer route,
application/provider linkage, launch URLs, dedicated groups, enrollment-flow
attributes, five prompt definitions, prompt-stage membership, user-write/email/
login stages, four ordered flow-stage bindings, fail-closed policy bindings,
exact per-application counts, distinct object identities, distinct non-empty
client secrets, and unchanged identities/secrets after re-import.

Final result:

- AUTHENTIK_IMPORT_OK
- DISPOSABLE_CLEANUP_OK
- independent host container/network/volume absence check: HOST_RESOURCE_ABSENCE_OK
- runtime-input SHA-256:
  99f199e21249790c5cc16ab7df61b68d5851c8d39a60a561afec056239736a06

No live tenant, external database, API credential, consumer repository,
deployment, or service was read or changed by this runtime gate.

## Retained failed attempts

Failures were not converted into passes:

1. The first development unit attempt used the host default Go 1.27.3 and was
   rejected by the exact Go 1.26.6 Makefile gate before runtime startup.
2. The next disposable run passed Go gates but the container could not traverse
   the mode-0700 temporary directory. It failed before importing either fixture.
   A separate bounded check proved all named containers and networks were gone.
3. An earlier complete oracle passed before the reviews strengthened exact image
   and cleanup enforcement. It is retained as preliminary evidence only. The
   final result above is the post-remediation gate.
4. A later exact-source run reached the importer before all three managed OAuth
   scope mappings existed and failed validation. Failure cleanup removed the
   named containers, network, and volume. The readiness oracle now waits for
   the exact three managed mappings before import; only the subsequent run is
   final evidence.

## Consumer admission

Revision-matched Mail evidence, its exact hashes and limitations are recorded in
docs/v1-consumer-admission.md. The inspected Mail evidence commit
c783440ebd46da8c851d75b3c384568c9777e729 is an ancestor of clean inspected
revision fdb367d63f69849bebb87c6260f74e2a5e54b3d8, and the relevant blobs are
unchanged. Board remains provenance plus a producer fixture. WorldWeave has only
an explicitly non-deployment producer fixture. No consumer was mutated.

## Reviews

Successive independent cold reviews found and drove remediation of URL
canonicality, size, identity, readiness, image, cleanup, workflow, golden, and
evidence defects. All findings and fixes are retained under review/. Two final
fresh independent passes reviewed runtime-input digest
99f199e21249790c5cc16ab7df61b68d5851c8d39a60a561afec056239736a06
and returned CLEAN. Their records are final-clean-review-1.md and
final-clean-review-2.md.
