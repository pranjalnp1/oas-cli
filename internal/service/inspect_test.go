package service

import "testing"

func TestInspect_ValidSpec(t *testing.T) {
	result, err := Inspect("../../testdata/valid/test.yml")
	if err != nil {
		t.Fatalf("Inspect returned unexpected error: %v", err)
	}

	if result.OpenAPIVersion != "3.0.0" {
		t.Errorf("OpenAPIVersion = %q, want %q", result.OpenAPIVersion, "3.0.0")
	}
	if result.Title != "Swagger Petstore" {
		t.Errorf("Title = %q, want %q", result.Title, "Swagger Petstore")
	}
	if result.Endpoints != 2 {
		t.Errorf("Endpoints = %d, want 2", result.Endpoints)
	}
	if result.Operations != 3 {
		t.Errorf("Operations = %d, want 3", result.Operations)
	}
	if result.Schemas != 3 {
		t.Errorf("Schemas = %d, want 3", result.Schemas)
	}
	if len(result.Servers) != 1 || result.Servers[0] != "http://petstore.swagger.io/v1" {
		t.Errorf("Servers = %v, want [\"http://petstore.swagger.io/v1\"]", result.Servers)
	}
}

func TestInspect_MissingOpenAPIField(t *testing.T) {
	_, err := Inspect("../../testdata/invalid/missing-openapi.yml")
	if err == nil {
		t.Fatal("Inspect returned nil error for a spec missing the openapi field, want an error")
	}
}

func TestInspect_MalformedSpec(t *testing.T) {
	_, err := Inspect("../../testdata/invalid/malformed.yml")
	if err == nil {
		t.Fatal("Inspect returned nil error for malformed YAML, want an error")
	}
}

func TestInspect_FileNotFound(t *testing.T) {
	_, err := Inspect("../../testdata/valid/does-not-exist.yml")
	if err == nil {
		t.Fatal("Inspect returned nil error for a nonexistent file, want an error")
	}
}
