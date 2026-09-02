# Product requirements

Give multiple products a consistent, reviewable Authentik enrollment and access
boundary. Every application receives its own slug, access group, enrollment
flow, verified-email activation sequence, and fail-closed policy binding.

The toolkit must render deterministically, reject unsafe identifiers and URLs,
detect sibling applications sharing an access group, and support exact desired
state comparison before an operator applies anything.

Non-goals: running Authentik, storing API credentials, owning tenants, managing
users at runtime, or applying remote changes implicitly.
