package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gotthboard/gotth-authentik/pkg/authentik"
)

const maxManifestBytes = 1 << 20

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(arguments []string, output io.Writer) error {
	if len(arguments) < 2 || len(arguments) > 3 ||
		(arguments[0] == "render" && len(arguments) != 2) ||
		(arguments[0] == "check" && len(arguments) != 3) ||
		(arguments[0] != "render" && arguments[0] != "check") {
		return fmt.Errorf("usage: gotth-authentik render <manifest.json> | check <manifest.json> <blueprint.yaml>")
	}
	application, err := readManifest(arguments[1])
	if err != nil {
		return err
	}
	blueprint, err := authentik.RenderEnrollmentBlueprint(application)
	if err != nil {
		return err
	}
	if arguments[0] == "render" {
		_, err = output.Write(blueprint)
		return err
	}
	current, err := os.ReadFile(arguments[2])
	if err != nil {
		return fmt.Errorf("read blueprint: %w", err)
	}
	if !bytes.Equal(current, blueprint) {
		return fmt.Errorf("blueprint differs from generated desired state")
	}
	return nil
}

func readManifest(path string) (authentik.Application, error) {
	file, err := os.Open(path)
	if err != nil {
		return authentik.Application{}, fmt.Errorf("read manifest: %w", err)
	}
	defer file.Close()
	raw, err := io.ReadAll(io.LimitReader(file, maxManifestBytes+1))
	if err != nil {
		return authentik.Application{}, fmt.Errorf("read manifest: %w", err)
	}
	if len(raw) > maxManifestBytes {
		return authentik.Application{}, fmt.Errorf("manifest exceeds 1 MiB")
	}
	var application authentik.Application
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&application); err != nil {
		return authentik.Application{}, fmt.Errorf("decode manifest: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return authentik.Application{}, fmt.Errorf("decode manifest: trailing data")
	}
	return application, nil
}
