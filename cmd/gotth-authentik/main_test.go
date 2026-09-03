package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

const validManifest = `{
  "slug":"example-app",
  "name":"Example application",
  "registration_title":"Register",
  "user_path":"users/example-app",
  "access_group":"example-app-users",
  "launch_url":"https://app.example/",
  "provider":{
    "name":"example-app-oidc",
    "client_type":"confidential",
    "client_id":"example-app",
    "redirect_uris":["https://app.example/auth/callback"],
    "authorization_flow":"default-provider-authorization-implicit-consent",
    "invalidation_flow":"default-provider-invalidation-flow",
    "signing_key":"authentik Self-signed Certificate"
  }
}`

func TestRunRenderAndCheck(t *testing.T) {
	directory := t.TempDir()
	manifest := filepath.Join(directory, "manifest.json")
	if err := os.WriteFile(manifest, []byte(validManifest), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := run([]string{"render", manifest}, &output); err != nil || output.Len() == 0 {
		t.Fatalf("run(render) = (%d bytes, %v)", output.Len(), err)
	}
	blueprint := filepath.Join(directory, "blueprint.yaml")
	if err := os.WriteFile(blueprint, output.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"check", manifest, blueprint}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run(check) returned error: %v", err)
	}
	if err := run([]string{"render", manifest}, failingWriter{}); err == nil {
		t.Fatal("render output error was hidden")
	}
	if err := os.WriteFile(blueprint, []byte("changed"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"check", manifest, blueprint}, &bytes.Buffer{}); err == nil {
		t.Fatal("changed blueprint passed check")
	}
}

func TestRunRejectsBadInputs(t *testing.T) {
	directory := t.TempDir()
	valid := filepath.Join(directory, "valid.json")
	if err := os.WriteFile(valid, []byte(validManifest), 0o600); err != nil {
		t.Fatal(err)
	}
	badJSON := filepath.Join(directory, "bad.json")
	if err := os.WriteFile(badJSON, []byte(validManifest+` {}`), 0o600); err != nil {
		t.Fatal(err)
	}
	unknown := filepath.Join(directory, "unknown.json")
	if err := os.WriteFile(unknown, []byte(`{"unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	oversize := filepath.Join(directory, "oversize.json")
	if err := os.WriteFile(oversize, []byte(strings.Repeat("x", maxManifestBytes+1)), 0o600); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(directory, "missing")
	invalidApplication := filepath.Join(directory, "invalid-application.json")
	if err := os.WriteFile(invalidApplication, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, arguments := range [][]string{
		nil,
		{"unknown", valid},
		{"render", valid, valid},
		{"check", valid},
		{"render", missing},
		{"render", badJSON},
		{"render", unknown},
		{"render", oversize},
		{"render", invalidApplication},
		{"render", directory},
		{"check", valid, missing},
	} {
		if err := run(arguments, &bytes.Buffer{}); err == nil {
			t.Errorf("run(%v) returned no error", arguments)
		}
	}
}
