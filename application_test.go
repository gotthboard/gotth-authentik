package authentik

import "testing"

func validApplication() Application {
	return Application{Slug: "example-app", Name: "Example App enrollment", RegistrationTitle: "Register for Example App", UserPath: "users/example-app", AccessGroup: "example-app-users"}
}

func TestApplicationValidationAndIsolation(t *testing.T) {
	valid := validApplication()
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}
	invalid := []Application{
		{},
		{Slug: "Bad Slug", Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup},
		{Slug: valid.Slug, Name: "", RegistrationTitle: valid.RegistrationTitle, UserPath: valid.UserPath, AccessGroup: valid.AccessGroup},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: "groups/no", AccessGroup: valid.AccessGroup},
		{Slug: valid.Slug, Name: valid.Name, RegistrationTitle: valid.RegistrationTitle, UserPath: "users/Bad", AccessGroup: valid.AccessGroup},
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
	if err := ValidateIsolation([]Application{valid, sibling}); err != nil {
		t.Fatalf("isolated siblings rejected: %v", err)
	}
}

func TestURLValidation(t *testing.T) {
	if err := ValidateIssuerURL("https://auth.example/application/o/example/"); err != nil {
		t.Fatalf("issuer rejected: %v", err)
	}
	for _, raw := range []string{"/relative", "ftp://auth.example/provider/", "https://user@auth.example/provider/", "https://auth.example/provider/?x=1", "https://auth.example/provider"} {
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
		{"https://auth.example/if/flow/example/", "https://auth.example/base/", "example"},
		{"https://auth.example/if/flow/example/", "https://auth.example/", "Bad"},
	} {
		if err := ValidateEnrollmentURL(values[0], values[1], values[2]); err == nil {
			t.Errorf("enrollment values %#v passed", values)
		}
	}
}
