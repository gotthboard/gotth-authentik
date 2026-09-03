# Judge pass 1 — rejected and repaired

The first cold review rejected a module-root compatibility facade. No released
tag or consumer pin established that import path, so it duplicated the public
desired-state API without preserving userspace.

Repair: remove the facade, retain exactly one public library package at
`pkg/authentik`, make the CLI import it directly, and keep the module root for
governance. Provider and secret boundaries remain unchanged.
