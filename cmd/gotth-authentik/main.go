package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gotthboard/gotth-authentik/pkg/authentik"
)

const (
	maxManifestBytes  = 1 << 20
	maxBlueprintBytes = 1 << 20
)

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
	current, err := readBoundedFile(arguments[2], maxBlueprintBytes, "blueprint")
	if err != nil {
		return err
	}
	if !bytes.Equal(current, blueprint) {
		return fmt.Errorf("blueprint differs from generated desired state")
	}
	return nil
}

func readManifest(path string) (authentik.Application, error) {
	raw, err := readBoundedFile(path, maxManifestBytes, "manifest")
	if err != nil {
		return authentik.Application{}, err
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

func readBoundedFile(path string, maximum int64, kind string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", kind, err)
	}
	defer file.Close()
	raw, err := io.ReadAll(io.LimitReader(file, maximum+1))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", kind, err)
	}
	if int64(len(raw)) > maximum {
		return nil, fmt.Errorf("%s exceeds 1 MiB", kind)
	}
	return raw, nil
}
