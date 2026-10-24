# v1 consumer admission

The v1 candidate reuses revision-matched real-consumer evidence where it exists
and adds isolated producer-owned fixtures only for missing proof. No GOTTH Mail,
GOTTH Board, or WorldWeave repository or deployment was changed.

## GOTTH Mail: real consumer evidence

Reusable evidence is bound to toolkit revision
77c811c18fe3c107f6f3e97ce3f9a3850a685300 and GOTTH Mail evidence revision
c783440ebd46da8c851d75b3c384568c9777e729. That revision is an ancestor of the
clean inspected Mail revision fdb367d63f69849bebb87c6260f74e2a5e54b3d8, and
the relevant manifest and blueprint blobs are unchanged.

- manifest: deploy/authentik/gotth-mail/manifest.json
- manifest SHA-256: a59e50d6f7c54418ccf68db5e6e5795c1f416c53c43521fa9fefac5529ddfc7b
- blueprint: deploy/authentik/gotth-mail/blueprint.yaml
- blueprint SHA-256: 5bcbe9141f7a59eb7416faa0de0ced9ae4f95ffdfb4b15408cf576eed369a9ad
- evidence: workflow/features/v2.identity-provisioning/live-authentik-persistence/evidence/2026-09-13-live-authentik-profile.md

The evidence records use of the exact toolkit base revision, validation and
application to Authentik, idempotent second import, unchanged generated client
secret, exact discovery/JWKS/authorization/enrollment checks, and isolation from
a retained historical provider. The v1 candidate also rechecked the unchanged
Mail manifest and blueprint byte-for-byte with the candidate CLI. The evidence
does not prove a completed browser callback or SCIM lifecycle, and this project
does not claim either.

## GOTTH Board: retained provenance plus producer fixture

examples/gotth-bb/accepted-alpha2-blueprint.yaml remains the accepted historical
consumer blueprint. examples/gotth-bb/manifest.json is a producer-owned v1
fixture that renders the complete provider/application boundary and is checked
by make verify. The historical blueprint predates that complete generic shape,
so it is provenance rather than revision-compatible proof of every v1 object.
No Board repository was changed or relabeled as a current package consumer.

## WorldWeave: isolated missing-proof fixture

No revision-matched WorldWeave package consumption or deployment evidence was
available. examples/worldweave/manifest.json is therefore explicitly a
producer-owned, non-deployment fixture. It proves that a second branded profile
with multiple exact redirects validates and renders without adding
product-specific code. It is not evidence that WorldWeave is deployed or pinned
to this module.

## Admission boundary

The reusable Mail evidence covers real consumer use at the exact base revision.
The candidate changes only validation, input bounds, verification, and contract
documentation; renderer output is unchanged and the Mail candidate check
passes. Board and WorldWeave fixtures cover producer-side profile portability.
The disposable two-application Authentik 2026.5.2 oracle independently proves
complete creation, isolation, nonduplication, linkage, and generated-secret
preservation. Consumer repository pins and deployments remain consumer-owned.
