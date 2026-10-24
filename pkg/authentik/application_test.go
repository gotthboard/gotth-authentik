package authentik

import (
	"strings"
	"testing"
)

func validApplication() Application {
	return Application{
		Slug: "example-app", Name: "Example App", RegistrationTitle: "Register for Example App",
		UserPath: "users/example-app", AccessGroup: "example-app-users", LaunchURL: "https://app.example/",
		Provider: Provider{
			Name: "example-app-oidc", ClientType: "confidential", ClientID: "example-app",
			RedirectURIs:      []string{"https://app.example/auth/callback"},
			AuthorizationFlow: "default-provider-authorization-implicit-consent",
			InvalidationFlow:  "default-provider-invalidation-flow", SigningKey: "authentik Self-signed Certificate",
		},
	}
}

func TestApplicationValidationAndIsolation(t *testing.T) {
	valid := validApplication()
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}
	invalid := []Application{
		{},
		{Slug: "Bad Slug", Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: "Bad Group", Provider: valid.Provider},
		{Slug: valid.Slug, Name: "", RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: strings.Repeat("users/a", 24), AccessGroup: valid.AccessGroup, Provider: valid.Provider},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: "groups/no", AccessGroup: valid.AccessGroup},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: "users/Bad", AccessGroup: valid.AccessGroup},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup, LaunchURL: "http://app.example/", Provider: valid.Provider},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup, LaunchURL: "https://APP.example/", Provider: valid.Provider},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup, LaunchURL: "https://app.example/" + strings.Repeat("x", 2049), Provider: valid.Provider},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup, Provider: Provider{}},
	}
	for index, application := range invalid {
		if err := application.Validate(); err == nil {
			t.Errorf("invalid application %d passed", index)
		}
	}
	if err := ValidateIsolation([]Application{valid, valid}); err == nil {
		t.Fatal("duplicate application passed isolation")
	}
	sibling := valid
	sibling.Slug = "sibling"
	if err := ValidateIsolation([]Application{valid, sibling}); err == nil {
		t.Fatal("shared group passed isolation")
	}
	sibling.AccessGroup = "sibling-users"
	sibling.Provider.Name = "sibling-oidc"
	sibling.Provider.ClientID = "sibling"
	if err := ValidateIsolation([]Application{valid, sibling}); err != nil {
		t.Fatalf("isolated siblings rejected: %v", err)
	}
}

func TestProviderValidation(t *testing.T) {
	valid := validApplication().Provider
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}
	invalid := []Provider{
		{},
		{Name: valid.Name, ClientType: "public", ClientID: valid.ClientID, RedirectURIs: valid.RedirectURIs, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: "device", ClientID: valid.ClientID, RedirectURIs: valid.RedirectURIs, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: "", RedirectURIs: valid.RedirectURIs, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"http://app.example/callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://APP.example/callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://app.example:443/callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://app.example.:8443/callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://app.example:/callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://[0:0:0:0:0:0:0:1]:8443/callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://app.example/%63allback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://app.example/a/../callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://127.0.0.01/callback"}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{valid.RedirectURIs[0], valid.RedirectURIs[0]}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: make([]string, 33), AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: []string{"https://app.example/" + strings.Repeat("x", 2049)}, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: valid.RedirectURIs, AuthorizationFlow: "Bad", InvalidationFlow: valid.InvalidationFlow, SigningKey: valid.SigningKey},
		{Name: valid.Name, ClientType: valid.ClientType, ClientID: valid.ClientID, RedirectURIs: valid.RedirectURIs, AuthorizationFlow: valid.AuthorizationFlow, InvalidationFlow: valid.InvalidationFlow, SigningKey: ""},
	}
	for index, provider := range invalid {
		if err := provider.Validate(); err == nil {
			t.Errorf("invalid provider %d passed", index)
		}
	}
	loopback := valid
	loopback.RedirectURIs = []string{"http://127.0.0.1:8080/callback", "http://[::1]:8080/callback"}
	if err := loopback.Validate(); err != nil {
		t.Fatalf("numeric loopback redirect rejected: %v", err)
	}
}

func TestTextValidation(t *testing.T) {
	if validText(string([]byte{0xff}), 1, 10) {
		t.Fatal("invalid UTF-8 passed")
	}
	if validText("line\nbreak", 1, 20) {
		t.Fatal("control character passed")
	}
	if validText("too-long", 1, 3) {
		t.Fatal("overlong text passed")
	}
}

func TestIsolationRejectsProviderAndClientReuse(t *testing.T) {
	first := validApplication()
	second := validApplication()
	second.Slug = "second"
	second.AccessGroup = "second-users"
	second.Provider.ClientID = "second"
	if err := ValidateIsolation([]Application{first, second}); err == nil {
		t.Fatal("shared provider name passed isolation")
	}
	second.Provider.Name = "second-oidc"
	second.Provider.ClientID = first.Provider.ClientID
	if err := ValidateIsolation([]Application{first, second}); err == nil {
		t.Fatal("shared client ID passed isolation")
	}
	if err := ValidateIsolation([]Application{{}}); err == nil {
		t.Fatal("invalid application passed isolation")
	}
}

func TestURLValidation(t *testing.T) {
	if err := ValidateIssuerURL("https://auth.example/application/o/example/"); err != nil {
		t.Fatalf("issuer rejected: %v", err)
	}
	for _, raw := range []string{
		"/relative",
		"ftp://auth.example/application/o/example/",
		"http://auth.example/application/o/example/",
		"http://127.0.0.1/application/o/example/",
		"https://user@auth.example/application/o/example/",
		"https://auth.example/application/o/example/?x=1",
		"https://auth.example/provider/",
		"https://auth.example/application/o/wrong/extra/",
		"https://auth.example/application/o/Bad/",
		"https://auth.example/application/o/example%2Dapp/",
		"https://AUTH.example/application/o/example/",
		"https://auth.example:443/application/o/example/",
		"https://auth.example.:8443/application/o/example/",
		"https://auth.example:/application/o/example/",
		"https://auth.example:0443/application/o/example/",
		"https://auth.example:65536/application/o/example/",
		"https://127.0.0.01/application/o/example/",
		"https://127.1/application/o/example/",
		"https://2130706433/application/o/example/",
		"https://0x7f000001/application/o/example/",
		"https://0x7f.1/application/o/example/",
		"https://127.0x0.1/application/o/example/",
		"https://127.0x.1/application/o/example/",
		"https://0x.1/application/o/example/",
		"https://0x.0x/application/o/example/",
		"https://example.1/application/o/example/",
		"https://auth..example/application/o/example/",
		"https://-auth.example/application/o/example/",
		"https://auth_example/application/o/example/",
		"https://auth.example/application/o/example/%2E%2E/",
		"https://auth.example/application/o/example",
	} {
		if err := ValidateIssuerURL(raw); err == nil {
			t.Errorf("issuer %q passed", raw)
		}
	}
	if err := ValidateEnrollmentURL("https://auth.example/if/flow/example-enrollment/", "https://auth.example/", "example-enrollment"); err != nil {
		t.Fatalf("enrollment URL rejected: %v", err)
	}
	for _, values := range [][3]string{
		{"https://other.example/if/flow/example/", "https://auth.example/", "example"},
		{"https://auth.example/if/flow/wrong/", "https://auth.example/", "example"},
		{"https://auth.example/if/flow/example/?x=1", "https://auth.example/", "example"},
		{"https://auth.example/if/flow/example%2Denrollment/", "https://auth.example/", "example-enrollment"},
		{"http://127.0.0.1/if/flow/example/", "https://auth.example/", "example"},
		{"https://auth.example/if/flow/example/", "http://127.0.0.1/", "example"},
		{"https://auth.example/if/flow/example/", "https://auth.example/%2F", "example"},
		{"https://AUTH.example/if/flow/example/", "https://auth.example/", "example"},
		{"https://auth.example/if/flow/example/", "https://AUTH.example/", "example"},
		{"https://auth.example:0443/if/flow/example/", "https://auth.example:0443/", "example"},
		{"https://127.0.0.01/if/flow/example/", "https://127.0.0.01/", "example"},
		{"https://auth.example/if/flow/example/", "https://auth.example/base/", "example"},
		{"https://auth.example/if/flow/example/", "https://auth.example/", "Bad"},
	} {
		if err := ValidateEnrollmentURL(values[0], values[1], values[2]); err == nil {
			t.Errorf("enrollment values %#v passed", values)
		}
	}
}
