#!/bin/sh
set -eu

if [ "$#" -ne 0 ]; then
	printf '%s\n' 'usage: verify-authentik-2026.5.2.sh' >&2
	exit 2
fi
GOTOOLCHAIN=go1.26.6
export GOTOOLCHAIN
version=$(go env GOVERSION)
[ "${version%%-*}" = go1.26.6 ] || {
	printf 'required Go 1.26.6, got %s\n' "$version" >&2
	exit 1
}

authentik_digest=sha256:3ddf09bbf69ded6a9634ecd753a01608d477f811e99bb5ffe9fc2ef7ad1c6581
postgres_digest=sha256:57c72fd2a128e416c7fcc499958864df5301e940bca0a56f58fddf30ffc07777
authentik_image=ghcr.io/goauthentik/server:2026.5.2@$authentik_digest
postgres_image=postgres:16@$postgres_digest
docker_bin=${DOCKER_BIN:-docker}
case_id="gotth-authentik-it-$$"
network_name="$case_id-net"
volume_name="$case_id-postgres-data"
postgres_name="$case_id-postgres"
server_name="$case_id-server"
worker_name="$case_id-worker"
scratch=

docker_cmd() {
	if [ "${DOCKER_SUDO:-0}" = 1 ]; then
		sudo -n "$docker_bin" "$@"
	else
		"$docker_bin" "$@"
	fi
}

cleanup_best_effort() {
	docker_cmd rm -f "$worker_name" "$server_name" "$postgres_name" >/dev/null 2>&1 || true
	docker_cmd network rm "$network_name" >/dev/null 2>&1 || true
	docker_cmd volume rm "$volume_name" >/dev/null 2>&1 || true
	if [ -n "$scratch" ]; then
		rm -rf "$scratch"
	fi
}

cleanup_checked() {
	docker_cmd rm -f "$worker_name" "$server_name" "$postgres_name" >/dev/null
	docker_cmd network rm "$network_name" >/dev/null
	docker_cmd volume rm "$volume_name" >/dev/null
	rm -rf "$scratch"
	for container in "$worker_name" "$server_name" "$postgres_name"; do
		if docker_cmd inspect "$container" >/dev/null 2>&1; then
			return 1
		fi
	done
	if docker_cmd network inspect "$network_name" >/dev/null 2>&1; then
		return 1
	fi
	if docker_cmd volume inspect "$volume_name" >/dev/null 2>&1; then
		return 1
	fi
	[ ! -e "$scratch" ]
}
trap cleanup_best_effort EXIT
trap 'exit 1' HUP INT TERM
scratch=$(mktemp -d)

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

docker_cmd image inspect --format '{{range .RepoDigests}}{{println .}}{{end}}' "$authentik_image" |
	grep -Fx "ghcr.io/goauthentik/server@$authentik_digest" >/dev/null
docker_cmd image inspect --format '{{range .RepoDigests}}{{println .}}{{end}}' "$postgres_image" |
	grep -Fx "postgres@$postgres_digest" >/dev/null
go run ./cmd/gotth-authentik render integration/fixture-manifest.json >"$scratch/primary.yaml"
go run ./cmd/gotth-authentik render integration/fixture-sibling-manifest.json >"$scratch/sibling.yaml"
go run ./cmd/gotth-authentik check integration/fixture-manifest.json "$scratch/primary.yaml"
go run ./cmd/gotth-authentik check integration/fixture-sibling-manifest.json "$scratch/sibling.yaml"
chmod 0755 "$scratch"
chmod 0644 "$scratch/primary.yaml" "$scratch/sibling.yaml"

docker_cmd network create --internal "$network_name" >/dev/null
docker_cmd volume create "$volume_name" >/dev/null
docker_cmd run -d --name "$postgres_name" --network "$network_name" \
	-v "$volume_name:/var/lib/postgresql/data" \
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
	-v "$scratch:/test:ro" \
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

wait_for "$server_name" 'from authentik.flows.models import Flow; required_flows={"default-provider-authorization-implicit-consent", "default-provider-invalidation-flow"}; assert set(Flow.objects.filter(slug__in=required_flows).values_list("slug", flat=True)) == required_flows; from authentik.crypto.models import CertificateKeyPair; assert CertificateKeyPair.objects.filter(name="authentik Self-signed Certificate").exists(); from authentik.providers.oauth2.models import ScopeMapping; required={"goauthentik.io/providers/oauth2/scope-openid", "goauthentik.io/providers/oauth2/scope-profile", "goauthentik.io/providers/oauth2/scope-email"}; assert set(ScopeMapping.objects.filter(managed__in=required).values_list("managed", flat=True)) == required' 180

if ! docker_cmd exec -i "$server_name" ak shell \
	<integration/assert-authentik-state.py >"$scratch/import.log" 2>&1; then
	cat "$scratch/import.log" >&2
	exit 1
fi
grep -Fx AUTHENTIK_IMPORT_OK "$scratch/import.log"
cleanup_checked
trap - EXIT HUP INT TERM
printf '%s\n' DISPOSABLE_CLEANUP_OK
