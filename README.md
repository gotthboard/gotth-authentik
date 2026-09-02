# gotth-authentik

`gotth-authentik` is a versioned desired-state and validation toolkit for
isolated Authentik application enrollment. It generates verified-email
enrollment flows, inactive-until-verified users, dedicated access groups, and
fail-closed application policy bindings.

It is not a daemon and does not silently apply remote changes. Rendering and
inspection are local; a future remote apply command must remain an explicit
operator action.

The accepted GOTTH Board alpha.2 blueprint is retained as provenance beside a
generic manifest that reproduces the same contract.
