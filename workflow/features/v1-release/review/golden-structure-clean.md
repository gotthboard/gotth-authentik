# Independent golden semantic review

Verdict: CLEAN.

Reviewed artifact:

- pkg/authentik/testdata/example-blueprint.yaml
- SHA-256: 8b60022d06315f3e8348654c15c588a451bd67f37e8cbc3cddd44b7a125b1444
- runtime-input digest:
  99f199e21249790c5cc16ab7df61b68d5851c8d39a60a561afec056239736a06
- v1 contract SHA-256:
  2a871fb5f14294b84d44e47961499d628fb439235e4e965074fabbd779bdd6ab
- validApplication source/test SHA-256:
  ff1aff9f0247ff72d36b3e1a387caf6a4ec3fb736e7801b591d88c1c8dff7b7c
- Authentik model target: 2026.5.2 image digest
  sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581
- independent reviewer identity: OpenClaw child run labeled
  `golden semantic review`, read-only, 2026-10-24

The review did not trust the renderer or byte-comparison test. It checked the
golden entry by entry against the v1 contract, validApplication fixture, and the
Authentik 2026.5.2 blueprint model:

- exactly 18 entries: provider, application, group, flow, five prompts, four
  stages, four ordered bindings, and one policy binding;
- confidential provider fields, strict callback, grants, ID-token claims,
  per-provider issuer, shared flows, signing key, and managed scopes;
- application/provider linkage and launch URL;
- non-superuser group and unauthenticated enrollment flow;
- exact prompt fields/orders and prompt-stage membership;
- inactive external user creation in the dedicated group;
- verified-email activation and login stage;
- binding orders 10, 20, 30, and 100;
- enabled, non-negated, fail-closed group policy;
- no client_secret or other secret material.

No omission or semantic mismatch was found.
