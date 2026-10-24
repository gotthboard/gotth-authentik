# Definitive review Two findings

Verdict: findings.

1. WHATWG hosts whose last label is a numeric form remained admissible.
2. The Makefile fixture loop could hide an earlier failure.
3. The workflow integration command had ignored prose arguments and did not pin
   Go inside the harness.
4. v1 workflow state/event/affected-file records were incomplete.
5. Golden review provenance did not bind its contract, fixture source, runtime,
   model target, and reviewer identity.

Remediation:

- reject domains ending in decimal, octal-like, or hexadecimal numeric labels;
- enable shell fail-fast before the fixture loop;
- reject harness arguments and enforce Go 1.26.6 internally;
- use in_progress state, append activation history, and scope the event log;
- bind the golden review to exact hashes, runtime digest, Authentik image, date,
  and independent reviewer label.
