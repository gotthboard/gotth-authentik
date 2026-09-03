# Architecture

The root library validates one small JSON-compatible application/provider model
and renders one complete Authentik blueprint deterministically. The provider is
an OAuth2/OIDC authorization-code provider with exact redirect URI matching and
the built-in `openid`, `profile`, and `email` mappings. Authentik creates the
confidential client secret because secret material is not manifest data.

The application refers to the provider with `!KeyOf`; its access binding refers
to the application with `!KeyOf`. Import therefore does not depend on a
pre-existing product object. Authorization flow, invalidation flow, and signing
key are explicit manifest references because installations may choose different
objects. The toolkit never owns those shared tenant primitives.

The CLI renders desired state or compares an existing blueprint byte-for-byte.
URL validators enforce exact canonical issuer, enrollment, launch, and redirect
targets. Isolation validation rejects reuse of application slugs, provider
names, client IDs, or access groups between sibling manifests.

The accepted alpha.2 YAML remains evidence, not the generic API. Consumer
manifests are data; no product-specific behavior is compiled into the library.
Remote mutation remains absent. Import verification runs only against a
disposable tenant and proves creation plus idempotent re-import.
