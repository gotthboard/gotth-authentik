#!/bin/sh
set -eu

authentik_image=${AUTHENTIK_TEST_IMAGE:-ghcr.io/goauthentik/server@sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581}
postgres_image=${AUTHENTIK_TEST_POSTGRES_IMAGE:-postgres@sha256:57c72fd2a128e416c7fcc499958864df5301e940bca0a56f58fddf30ffc07777}
docker_bin=${DOCKER_BIN:-docker}
case_id="gotth-authentik-it-$$"
network_name="$case_id-net"
postgres_name="$case_id-postgres"
server_name="$case_id-server"
worker_name="$case_id-worker"
scratch=$(mktemp -d)

docker_cmd() {
	if [ "${DOCKER_SUDO:-0}" = 1 ]; then
		sudo -n "$docker_bin" "$@"
	else
		"$docker_bin" "$@"
	fi
}

cleanup() {
	docker_cmd rm -f "$worker_name" "$server_name" "$postgres_name" >/dev/null 2>&1 || true
	docker_cmd network rm "$network_name" >/dev/null 2>&1 || true
	rm -rf "$scratch"
}
trap cleanup EXIT HUP INT TERM

fail_with_logs() {
	docker_cmd logs --tail 160 "$server_name" >&2 2>/dev/null || true
	docker_cmd logs --tail 160 "$worker_name" >&2 2>/dev/null || true
	exit 1
}

wait_for() {
	container=$1
	command=$2
	limit=$3
	i=0
	until docker_cmd exec "$container" ak shell -c "$command" >/dev/null 2>&1; do
		i=$((i + 1))
		if [ "$i" -ge "$limit" ]; then
			fail_with_logs
		fi
		sleep 1
	done
}

docker_cmd image inspect "$authentik_image" >/dev/null
docker_cmd image inspect "$postgres_image" >/dev/null
go run ./cmd/gotth-authentik render integration/fixture-manifest.json >"$scratch/blueprint.yaml"

docker_cmd network create --internal "$network_name" >/dev/null
docker_cmd run -d --name "$postgres_name" --network "$network_name" \
	-e POSTGRES_DB=authentik -e POSTGRES_USER=authentik \
	-e POSTGRES_PASSWORD=disposable-password "$postgres_image" >/dev/null

i=0
until docker_cmd exec "$postgres_name" pg_isready -U authentik -d authentik >/dev/null 2>&1; do
	i=$((i + 1))
	if [ "$i" -ge 60 ]; then
		fail_with_logs
	fi
	sleep 1
done

docker_cmd run -d --name "$server_name" --network "$network_name" \
	-v "$scratch/blueprint.yaml:/test/blueprint.yaml:ro" \
	-e AUTHENTIK_SECRET_KEY=disposable-only-secret-key-at-least-fifty-characters-1234567890 \
	-e POSTGRES_PASSWORD=disposable-password \
	-e AUTHENTIK_POSTGRESQL__HOST="$postgres_name" \
	"$authentik_image" server >/dev/null

wait_for "$server_name" 'from django.db import connection; connection.ensure_connection()' 240

docker_cmd run -d --name "$worker_name" --network "$network_name" \
	-e AUTHENTIK_SECRET_KEY=disposable-only-secret-key-at-least-fifty-characters-1234567890 \
	-e POSTGRES_PASSWORD=disposable-password \
	-e AUTHENTIK_POSTGRESQL__HOST="$postgres_name" \
	"$authentik_image" worker >/dev/null

wait_for "$server_name" 'from authentik.flows.models import Flow; assert Flow.objects.filter(slug="default-provider-authorization-implicit-consent").exists(); from authentik.crypto.models import CertificateKeyPair; assert CertificateKeyPair.objects.filter(name="authentik Self-signed Certificate").exists()' 180

if ! docker_cmd exec "$server_name" ak shell -c '
from pathlib import Path
from authentik.blueprints.v1.importer import Importer
from authentik.core.models import Application, Group
from authentik.policies.models import PolicyBinding
from authentik.providers.oauth2.models import OAuth2Provider

raw = Path("/test/blueprint.yaml").read_text()
first = Importer.from_string(raw)
valid, logs = first.validate(raise_validation_errors=True)
assert valid, logs
assert first.apply()

provider = OAuth2Provider.objects.get(name="integration-example-oidc")
assert provider.client_secret
secret = provider.client_secret
assert provider.client_type == "confidential"
assert provider.client_id == "integration-example"
assert len(provider.redirect_uris) == 1
assert provider.redirect_uris[0].url == "https://app.example/auth/callback"

application = Application.objects.get(slug="integration-example")
assert application.provider_id == provider.pk
group = Group.objects.get(name="integration-example-users")
binding = PolicyBinding.objects.get(target=application, group=group)
assert binding.enabled and not binding.negate and not binding.failure_result

second = Importer.from_string(raw)
valid, logs = second.validate(raise_validation_errors=True)
assert valid, logs
assert second.apply()
provider.refresh_from_db()
assert provider.client_secret == secret
assert OAuth2Provider.objects.filter(name="integration-example-oidc").count() == 1
assert Application.objects.filter(slug="integration-example").count() == 1
print("AUTHENTIK_IMPORT_OK")
' >"$scratch/import.log" 2>&1; then
	cat "$scratch/import.log" >&2
	exit 1
fi
grep -Fx AUTHENTIK_IMPORT_OK "$scratch/import.log"
