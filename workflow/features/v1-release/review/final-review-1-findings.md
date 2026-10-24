# Final review 1 findings

Verdict: finding.

High: canonical URL validation still admitted malformed DNS labels,
noncanonical numeric hosts, leading-zero/default-like ports, and out-of-range
ports.

Remediation: authority validation now requires canonical decimal ports in the
1-65535 range, canonical IP text, or bounded lowercase DNS labels with valid
label syntax. Numeric-looking non-IP hosts are rejected, with public-surface
regression cases.
