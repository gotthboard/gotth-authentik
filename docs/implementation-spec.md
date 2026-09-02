# Implementation specification

- One manifest per application.
- Slugs, access groups, and user-path segments are bounded lowercase slugs.
- Display strings are valid bounded UTF-8 without control characters.
- New users are inactive until the email stage succeeds.
- The creation stage adds users only to that application's access group.
- The application policy binding is enabled, non-negated, and fail-closed.
- Blueprint rendering is deterministic and secret-free.
- `render` writes desired YAML; `check` compares it with an existing file.
- Remote apply is intentionally absent from the extraction baseline.
