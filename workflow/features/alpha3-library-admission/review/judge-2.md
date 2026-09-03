# Judge pass 2 — clean

Reviewed revision: `cb87f0ebe27d4f2bc6574bd409fcf4d28b23021a`.

All library and white-box test files are exact 100% renames into
`pkg/authentik`. The CLI changes only its import to the canonical package. The
root has no Go package, the public package is unique, and blueprint generation
remains distinct from operator-owned application/import execution.

Verdict: CLEAN.
