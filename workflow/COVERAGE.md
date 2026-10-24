# Coverage

The canonical behavior map is `workflow/artifacts/global-coverage-map.md`.
`make verify` covers `pkg/authentik`, the bounded CLI, deterministic rendering,
and desired-state checking. Current statement coverage is 98.8% for the public
library and 92.7% for the CLI; disposable Authentik import is a separate
runtime-boundary gate.
