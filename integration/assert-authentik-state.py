from pathlib import Path

from django.urls import reverse

from authentik.blueprints.v1.importer import Importer
from authentik.core.models import Application, Group
from authentik.flows.models import Flow, FlowStageBinding
from authentik.policies.models import PolicyBinding
from authentik.providers.oauth2.models import OAuth2Provider
from authentik.stages.email.models import EmailStage
from authentik.stages.prompt.models import Prompt, PromptStage
from authentik.stages.user_login.models import UserLoginStage
from authentik.stages.user_write.models import UserWriteStage

SPECS = (
    {
        "blueprint": "/test/primary.yaml",
        "slug": "integration-example",
        "name": "Integration Example",
        "title": "Register for Integration Example",
        "path": "users/integration-example",
        "group": "integration-example-users",
        "launch_url": "https://app.example/",
        "provider": "integration-example-oidc",
        "client_id": "integration-example",
        "redirects": ("https://app.example/auth/callback",),
    },
    {
        "blueprint": "/test/sibling.yaml",
        "slug": "integration-sibling",
        "name": "Integration Sibling",
        "title": "Register for Integration Sibling",
        "path": "users/integration-sibling",
        "group": "integration-sibling-users",
        "launch_url": "https://sibling.example/",
        "provider": "integration-sibling-oidc",
        "client_id": "integration-sibling",
        "redirects": (
            "https://sibling.example/auth/callback",
            "https://sibling.example/auth/secondary-callback",
        ),
    },
)

AUTHORIZATION_FLOW = "default-provider-authorization-implicit-consent"
INVALIDATION_FLOW = "default-provider-invalidation-flow"
SIGNING_KEY = "authentik Self-signed Certificate"
SCOPES = {
    "goauthentik.io/providers/oauth2/scope-openid",
    "goauthentik.io/providers/oauth2/scope-profile",
    "goauthentik.io/providers/oauth2/scope-email",
}
PROMPTS = (
    ("username", "Username", "username", True, 0),
    ("password", "Password", "password", True, 1),
    ("password_repeat", "Password (repeat)", "password", True, 2),
    ("name", "Display name", "text", True, 3),
    ("email", "Email", "email", True, 4),
)


def import_blueprint(spec):
    raw = Path(spec["blueprint"]).read_text()
    importer = Importer.from_string(raw)
    valid, logs = importer.validate(raise_validation_errors=True)
    assert valid, logs
    assert OAuth2Provider.objects.filter(name=spec["provider"]).count() == 0
    assert Application.objects.filter(slug=spec["slug"]).count() == 0
    assert Group.objects.filter(name=spec["group"]).count() == 0
    assert Flow.objects.filter(slug=f"{spec['slug']}-enrollment").count() == 0
    assert importer.apply()


def reimport_blueprint(spec):
    raw = Path(spec["blueprint"]).read_text()
    importer = Importer.from_string(raw)
    valid, logs = importer.validate(raise_validation_errors=True)
    assert valid, logs
    assert importer.apply()


def assert_spec(spec):
    slug = spec["slug"]
    provider = OAuth2Provider.objects.get(name=spec["provider"])
    assert provider.client_type == "confidential"
    assert provider.client_id == spec["client_id"]
    assert provider.client_secret
    assert provider.authorization_flow.slug == AUTHORIZATION_FLOW
    assert provider.invalidation_flow.slug == INVALIDATION_FLOW
    assert provider.signing_key.name == SIGNING_KEY
    assert set(provider.grant_types) == {"authorization_code", "refresh_token"}
    assert provider.issuer_mode == "per_provider"
    assert provider.include_claims_in_id_token is True
    assert {mapping.managed for mapping in provider.property_mappings.all()} == SCOPES
    redirects = tuple(
        (str(item.matching_mode), item.url, str(item.redirect_uri_type))
        for item in provider.redirect_uris
    )
    assert redirects == tuple(
        ("strict", redirect, "authorization") for redirect in spec["redirects"]
    )

    application = Application.objects.get(slug=slug)
    assert application.name == spec["name"]
    assert application.provider_id == provider.pk
    assert application.meta_launch_url == spec["launch_url"]
    assert reverse(
        "authentik_providers_oauth2:provider-root",
        kwargs={"application_slug": slug},
    ) == f"/application/o/{slug}/"

    group = Group.objects.get(name=spec["group"])
    assert group.is_superuser is False
    flow = Flow.objects.get(slug=f"{slug}-enrollment")
    assert flow.name == spec["name"]
    assert flow.title == spec["title"]
    assert str(flow.designation) == "enrollment"
    assert str(flow.authentication) == "require_unauthenticated"

    prompt_objects = []
    for suffix, expected in zip(
        ("username", "password", "password-repeat", "name", "email"),
        PROMPTS,
        strict=True,
    ):
        prompt = Prompt.objects.get(name=f"{slug}-{suffix}")
        actual = (
            prompt.field_key,
            prompt.label,
            str(prompt.type),
            prompt.required,
            prompt.order,
        )
        assert actual == expected
        prompt_objects.append(prompt)

    prompt_stage = PromptStage.objects.get(name=f"{slug}-prompts")
    assert set(prompt_stage.fields.values_list("pk", flat=True)) == {
        prompt.pk for prompt in prompt_objects
    }
    user_write = UserWriteStage.objects.get(name=f"{slug}-user-write")
    assert user_write.create_users_as_inactive is True
    assert str(user_write.user_creation_mode) == "always_create"
    assert str(user_write.user_type) == "external"
    assert user_write.user_path_template == spec["path"]
    assert user_write.create_users_group_id == group.pk
    email = EmailStage.objects.get(name=f"{slug}-email-verification")
    assert email.use_global_settings is True
    assert email.template == "email/account_confirmation.html"
    assert email.activate_user_on_success is True
    login = UserLoginStage.objects.get(name=f"{slug}-login")

    binding_objects = list(FlowStageBinding.objects.filter(target=flow))
    bindings = {binding.order: binding.stage_id for binding in binding_objects}
    assert bindings == {
        10: prompt_stage.pk,
        20: user_write.pk,
        30: email.pk,
        100: login.pk,
    }
    policy = PolicyBinding.objects.get(target=application)
    assert policy.group_id == group.pk
    assert policy.order == 20
    assert policy.enabled is True
    assert policy.negate is False
    assert policy.failure_result is False
    return {
        "provider": provider.pk,
        "application": application.pk,
        "group": group.pk,
        "flow": flow.pk,
        "prompts": tuple(prompt.pk for prompt in prompt_objects),
        "stages": (prompt_stage.pk, user_write.pk, email.pk, login.pk),
        "bindings": tuple(sorted(binding.pk for binding in binding_objects)),
        "policy": policy.pk,
        "secret": provider.client_secret,
    }


def assert_exact_counts():
    slugs = [spec["slug"] for spec in SPECS]
    names = [spec["provider"] for spec in SPECS]
    groups = [spec["group"] for spec in SPECS]
    flow_slugs = [f"{slug}-enrollment" for slug in slugs]
    assert OAuth2Provider.objects.filter(name__in=names).count() == len(SPECS)
    assert Application.objects.filter(slug__in=slugs).count() == len(SPECS)
    assert Group.objects.filter(name__in=groups).count() == len(SPECS)
    assert Flow.objects.filter(slug__in=flow_slugs).count() == len(SPECS)
    for spec in SPECS:
        slug = spec["slug"]
        assert Prompt.objects.filter(name__startswith=f"{slug}-").count() == 5
        assert PromptStage.objects.filter(name=f"{slug}-prompts").count() == 1
        assert UserWriteStage.objects.filter(name=f"{slug}-user-write").count() == 1
        assert EmailStage.objects.filter(name=f"{slug}-email-verification").count() == 1
        assert UserLoginStage.objects.filter(name=f"{slug}-login").count() == 1
        flow = Flow.objects.get(slug=f"{slug}-enrollment")
        assert FlowStageBinding.objects.filter(target=flow).count() == 4
        app = Application.objects.get(slug=slug)
        assert PolicyBinding.objects.filter(target=app).count() == 1


for spec in SPECS:
    import_blueprint(spec)
first = {spec["slug"]: assert_spec(spec) for spec in SPECS}
assert_exact_counts()
assert len({state["provider"] for state in first.values()}) == len(SPECS)
assert len({state["application"] for state in first.values()}) == len(SPECS)
assert len({state["group"] for state in first.values()}) == len(SPECS)
assert len({state["flow"] for state in first.values()}) == len(SPECS)
assert len({state["secret"] for state in first.values()}) == len(SPECS)

for spec in SPECS:
    reimport_blueprint(spec)
second = {spec["slug"]: assert_spec(spec) for spec in SPECS}
assert_exact_counts()
for spec in SPECS:
    slug = spec["slug"]
    for identity in (
        "provider",
        "application",
        "group",
        "flow",
        "prompts",
        "stages",
        "bindings",
        "policy",
        "secret",
    ):
        assert second[slug][identity] == first[slug][identity]
assert len({state["secret"] for state in second.values()}) == len(SPECS)

print("AUTHENTIK_IMPORT_OK")
