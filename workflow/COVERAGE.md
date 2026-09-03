# Coverage

The canonical behavior map is `workflow/artifacts/global-coverage-map.md`.
`make verify` covers `pkg/authentik`, the bounded CLI, deterministic rendering,
and desired-state checking. Current statement coverage is 99.1% for the public
library and 91.9% for the CLI; disposable Authentik import is a separate
runtime-boundary gate.
