.PHONY: verify

verify:
	version="$$(go env GOVERSION)"; test "$${version%%-*}" = "go1.26.6"
	test -z "$$(gofmt -l $$(git ls-files --cached --others --exclude-standard -- '*.go'))"
	go vet -mod=readonly ./...
	go test -mod=readonly -race -cover ./...
	blueprint="$$(mktemp)"; trap 'rm -f "$$blueprint"' EXIT; go run ./cmd/gotth-authentik render examples/gotth-bb/manifest.json >"$$blueprint"; go run ./cmd/gotth-authentik check examples/gotth-bb/manifest.json "$$blueprint"
