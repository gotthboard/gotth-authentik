# Coverage map

| Surface | Evidence |
|---|---|
| Application/provider manifest validation and sibling isolation | `application_test.go` |
| Exact issuer and enrollment URLs | `application_test.go` |
| Complete secret-free provider/application/enrollment blueprint | `blueprint_test.go` |
| Bounded strict render and desired-state comparison CLI | `cmd/gotth-authentik/main_test.go` |
| Authentik 2026.5.2 creation, linkage, re-import, and secret preservation | `integration/verify-authentik-2026.5.2.sh` |
| Accepted alpha.2 configuration | `examples/gotth-bb/accepted-alpha2-blueprint.yaml` |

Implementation tests above now live under `pkg/authentik/`; the root
`public_api_test.go` and CLI import that canonical package.
