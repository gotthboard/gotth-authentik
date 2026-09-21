# Mail backchannel admission

This feature adds the smallest reusable deployment boundary needed by a mail
consumer: a confidential OIDC provider whose client secret is loaded by
Authentik from a file, and one SCIM provider attached to the same application
as a backchannel.

The feature is local desired-state generation. It neither reads credentials
nor applies state to an Authentik tenant. The canonical contracts are:

- [PRD](prd.md)
- [architecture](architecture.md)
- [implementation specification](implementation-spec.md)

Admission records:

- [verification evidence](evidence/verification.md)
- [cold review](evidence/review.md)
