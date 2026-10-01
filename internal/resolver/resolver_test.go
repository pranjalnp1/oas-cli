package resolver

import (
	"testing"

	"github.com/yourusername/oas-cli/types"
)

func testSpec() *api_types.OpenAPISpec {
	return &api_types.OpenAPISpec{
		Components: api_types.ComponentsObject{
			Schemas: map[string]interface{}{
				"Pet": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{"type": "string"},
					},
				},
				"PetRef": map[string]interface{}{
					"$ref": "#/components/schemas/Pet",
				},
				"CircularA": map[string]interface{}{
					"$ref": "#/components/schemas/CircularB",
				},
				"CircularB": map[string]interface{}{
					"$ref": "#/components/schemas/CircularA",
				},
			},
		},
	}
}

func TestResolve_ValidRef(t *testing.T) {
	spec := testSpec()

	schema, err := Resolve(spec, "#/components/schemas/Pet")
	if err != nil {
		t.Fatalf("Resolve returned unexpected error: %v", err)
	}
	if schema["type"] != "object" {
		t.Errorf("resolved schema type = %v, want %q", schema["type"], "object")
	}
}

func TestResolve_MissingRef(t *testing.T) {
	spec := testSpec()

	_, err := Resolve(spec, "#/components/schemas/DoesNotExist")
	if err == nil {
		t.Fatal("Resolve returned nil error for a missing schema, want an error")
	}
}

func TestResolve_InvalidSyntax(t *testing.T) {
	spec := testSpec()

	cases := []string{
		"not-a-ref-at-all",
		"#/components/schemas/",
		"#/definitions/Pet", // Swagger 2.0 style, unsupported
	}
	for _, ref := range cases {
		if _, err := Resolve(spec, ref); err == nil {
			t.Errorf("Resolve(%q) returned nil error, want an error for invalid ref syntax", ref)
		}
	}
}

func TestResolveName(t *testing.T) {
	spec := testSpec()

	name, err := ResolveName(spec, "#/components/schemas/Pet")
	if err != nil {
		t.Fatalf("ResolveName returned unexpected error: %v", err)
	}
	if name != "Pet" {
		t.Errorf("ResolveName = %q, want %q", name, "Pet")
	}
}

func TestResolveSchema_NestedRef(t *testing.T) {
	spec := testSpec()

	resolved, err := ResolveSchema(spec, map[string]interface{}{"$ref": "#/components/schemas/PetRef"})
	if err != nil {
		t.Fatalf("ResolveSchema returned unexpected error: %v", err)
	}
	if resolved["type"] != "object" {
		t.Errorf("resolved schema type = %v, want %q (expected PetRef to resolve through to Pet)", resolved["type"], "object")
	}
}

func TestResolveSchema_NotARef(t *testing.T) {
	spec := testSpec()

	plain := map[string]interface{}{"type": "string"}
	resolved, err := ResolveSchema(spec, plain)
	if err != nil {
		t.Fatalf("ResolveSchema returned unexpected error: %v", err)
	}
	if resolved["type"] != "string" {
		t.Errorf("resolved schema type = %v, want %q (expected a non-ref schema to pass through unchanged)", resolved["type"], "string")
	}
}

func TestResolveSchema_CircularRef(t *testing.T) {
	spec := testSpec()

	_, err := ResolveSchema(spec, map[string]interface{}{"$ref": "#/components/schemas/CircularA"})
	if err == nil {
		t.Fatal("ResolveSchema returned nil error for a circular $ref chain, want an error")
	}
}
