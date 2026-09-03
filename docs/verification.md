# Verification

`make verify` checks format, vet, race tests, coverage, generation of the GOTTH
Board example, and byte-for-byte desired-state comparison.

`integration/verify-authentik-2026.5.2.sh` requires preloaded pinned Authentik
and PostgreSQL images plus a usable Docker CLI. Set `DOCKER_SUDO=1` only on a
host whose operator requires passwordless sudo for Docker. The harness creates
an internal network with no external route, a fresh PostgreSQL database, and
fresh Authentik server/worker containers. It imports the generic fixture twice,
checks provider/application linkage and the fail-closed group binding, proves
that the generated confidential secret is non-empty and unchanged, and removes
all state on exit.

The gate passed against Authentik 2026.5.2 image digest
`sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581`
and PostgreSQL 16 image digest
`sha256:57c72fd2a128e416c7fcc499958864df5301e940bca0a56f58fddf30ffc07777`.
No live Authentik tenant is a valid test target.
