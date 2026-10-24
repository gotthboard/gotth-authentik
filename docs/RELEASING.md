# Release Policy

`gotth-authentik` versions independently under Semantic Versioning. A version number
in another GOTTH repository creates no release obligation here.

If this repository contains a Go module, release tags use the Go-compatible
`vMAJOR.MINOR.PATCH` form, including any SemVer prerelease suffix. Historical
tags are immutable even when they predate this convention.

No release may be tagged until required verification is complete. The
repository uses the MIT license. Reserved placeholders are not tagged at all.

For an admitted release:

1. Select a version from this repository's actual compatibility change.
2. Verify an exact clean Forgejo commit with the repository's required gates.
3. Update `docs/CHANGELOG.md` with the exact version, commit, compatibility
   effect, verification, and known limitations.
4. Create and push a new annotated tag on canonical Forgejo. Never move or
   replace an existing tag.
5. Verify the exact Forgejo branch, annotated tag object, peeled commit,
   stable hosted Release, and downloaded asset bytes. Leave the existing
   one-way mirror enabled.
6. Hand the exact Forgejo release identity, notes, assets, checksums, support
   boundary, and verification record to Bryce for GitHub hosted-release
   publication and independent public-distribution verification.

Forgejo remains the development source of truth. A verified stable Forgejo
Release completes canonical publication. GitHub remains the public clone and Go
import/distribution endpoint, but GitHub hosted-release publication is Bryce's
separate handoff and does not block the authorized Forgejo release. Mirror or
GitHub failure never authorizes reversing mirror direction, moving a tag, or
publishing different bytes.
