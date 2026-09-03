package authentik_test

import (
	"testing"

	authentik "github.com/gotthboard/gotth-authentik/pkg/authentik"
)

func TestCanonicalPublicPackageIsUsable(t *testing.T) {
	t.Parallel()

	var _ = authentik.Application{}
	var _ = authentik.Provider{}
	if err := authentik.ValidateIssuerURL("https://identity.example/application/o/example/"); err != nil {
		t.Fatal(err)
	}
}
