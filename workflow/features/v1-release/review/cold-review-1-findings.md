# Cold review 1

Verdict: findings.

1. High: loopback HTTP was accepted by the shared URL validator for launch,
   issuer, and enrollment URLs although the v1 contract allowed it only for
   redirect URIs.
2. Medium: enrollment validation did not reject encoded RawPath values.
3. Medium: the runtime oracle did not prove confidential secrets differed
   between sibling providers.

Remediation:

- require strict HTTPS for issuer and enrollment URLs while preserving numeric
  loopback launch and redirect URLs used by the revision-matched Mail consumer;
- reject encoded enrollment/origin paths and add regression cases;
- assert distinct non-empty secrets before and after re-import.
