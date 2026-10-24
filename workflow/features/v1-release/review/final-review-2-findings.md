# Final review 2 findings

Verdict: findings.

1. High: encoded and dot-segment path aliases and noncanonical authorities could
   bypass raw redirect deduplication.
2. High: an unbounded launch URL could produce a blueprint larger than the
   check command's 1 MiB input limit.
3. Medium: the golden lacked hash-bound independent semantic review provenance.
4. Medium: the runtime-input digest recipe did not state the required bytewise
   locale.
5. Medium: verification prose overclaimed every provider field rather than the
   v1 provider contract fields.

Remediation:

- reject encoded, repeated-separator, backslash, and dot-segment paths plus
  noncanonical authority aliases;
- bound launch URLs to 2048 bytes and add a maximum round-trip render/check test;
- retain an independent hash-bound semantic golden review;
- record the exact LC_ALL=C digest recipe;
- enumerate the provider fields actually asserted by the runtime oracle.
