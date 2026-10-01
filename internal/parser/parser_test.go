package parser

import (
	"os"
	"testing"
)

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", path, err)
	}
	return data
}

func TestParse_ValidYAML(t *testing.T) {
	data := mustReadFile(t, "../../testdata/valid/test.yml")

	spec, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	if spec.OpenAPI != "3.0.0" {
		t.Errorf("OpenAPI = %q, want %q", spec.OpenAPI, "3.0.0")
	}
	if spec.Info.Title != "Swagger Petstore" {
		t.Errorf("Info.Title = %q, want %q", spec.Info.Title, "Swagger Petstore")
	}
	if len(spec.Paths) != 2 {
		t.Errorf("len(Paths) = %d, want 2", len(spec.Paths))
	}
}

func TestParse_ValidJSON(t *testing.T) {
	data := mustReadFile(t, "../../testdata/validLiveAPI/petStore3.json")

	spec, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	if spec.OpenAPI == "" {
		t.Error("OpenAPI version is empty, expected it to be populated from JSON input")
	}
	if spec.Info.Title == "" {
		t.Error("Info.Title is empty, expected it to be populated from JSON input")
	}
	if len(spec.Paths) == 0 {
		t.Error("Paths is empty, expected real paths to be parsed from JSON input")
	}
}

func TestParse_MalformedYAML(t *testing.T) {
	data := mustReadFile(t, "../../testdata/invalid/malformed.yml")

	_, err := Parse(data)
	if err == nil {
		t.Fatal("Parse returned nil error for malformed YAML, want an error")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	spec, err := Parse([]byte{})
	if err != nil {
		t.Fatalf("Parse returned unexpected error for empty input: %v", err)
	}
	if spec.OpenAPI != "" {
		t.Errorf("OpenAPI = %q, want empty for empty input", spec.OpenAPI)
	}
}
