# Admission review Beta findings

Verdict: findings.

1. WHATWG hexadecimal IPv4 aliases remained admissible.
2. Re-import identity preservation covered only provider, application, group,
   and flow, not prompts, stages, bindings, and policy.
3. The scratch directory was allocated before cleanup traps existed.
4. Docker image ID comparison depended on image-store implementation.

Remediation:

- reject decimal/octal/hexadecimal numeric-host candidates unless net.ParseIP
  accepts their canonical form;
- capture and compare every owned prompt, stage, flow binding, and policy ID;
- install guarded EXIT/signal traps before mktemp;
- verify exact registry references through RepoDigests rather than .Id.
