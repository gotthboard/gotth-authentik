.PHONY: verify

verify:
	version="$$(go env GOVERSION)"; test "$${version%%-*}" = "go1.26.6"
	test -z "$$(gofmt -l $$(git ls-files --cached --others --exclude-standard -- '*.go'))"
	go vet -mod=readonly ./...
	go test -mod=readonly -race -cover ./...
	set -e; blueprint="$$(mktemp)"; trap 'rm -f "$$blueprint"' EXIT; for manifest in examples/gotth-bb/manifest.json examples/worldweave/manifest.json; do go run ./cmd/gotth-authentik render "$$manifest" >"$$blueprint"; go run ./cmd/gotth-authentik check "$$manifest" "$$blueprint"; done
