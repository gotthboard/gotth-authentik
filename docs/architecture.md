# Architecture

The root library validates one small JSON-compatible application model and
renders an Authentik blueprint deterministically. The CLI renders desired state
or compares an existing blueprint byte-for-byte. URL validators enforce exact
canonical issuer and enrollment targets. Isolation validation treats group
reuse between sibling applications as an error.

The accepted alpha.2 YAML remains evidence, not the generic API. Consumer
manifests are data; no product-specific behavior is compiled into the library.
