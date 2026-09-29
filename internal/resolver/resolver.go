package resolver

import (
	"fmt"
	"strings"

	"github.com/yourusername/oas-cli/types"
)

const maxDepth = 10

// Resolve follows a local $ref string (e.g. "#/components/schemas/Pet")
// and returns the schema object it points to.
func Resolve(spec *api_types.OpenAPISpec, ref string) (map[string]interface{}, error) {
	name, err := parseRef(ref)
	if err != nil {
		return nil, err
	}

	schema, ok := spec.Components.Schemas[name].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unresolved $ref: %s", ref)
	}
	return schema, nil
}

// ResolveName follows a local $ref string and returns just the referenced
// component's name (e.g. "#/components/schemas/Pet" -> "Pet"), verifying
// that the reference actually exists.
func ResolveName(spec *api_types.OpenAPISpec, ref string) (string, error) {
	name, err := parseRef(ref)
	if err != nil {
		return "", err
	}
	if _, ok := spec.Components.Schemas[name]; !ok {
		return "", fmt.Errorf("unresolved $ref: %s", ref)
	}
	return name, nil
}

// ResolveSchema resolves a schema map that may itself be a $ref, following
// nested $refs until a concrete schema is reached. Returns an error if the
// reference chain exceeds maxDepth, which guards against circular references.
func ResolveSchema(spec *api_types.OpenAPISpec, schema map[string]interface{}) (map[string]interface{}, error) {
	return resolveWithDepth(spec, schema, 0)
}

func resolveWithDepth(spec *api_types.OpenAPISpec, schema map[string]interface{}, depth int) (map[string]interface{}, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("possible circular $ref detected")
	}

	ref, isRef := schema["$ref"].(string)
	if !isRef {
		return schema, nil
	}

	resolved, err := Resolve(spec, ref)
	if err != nil {
		return nil, err
	}
	return resolveWithDepth(spec, resolved, depth+1)
}

// parseRef validates that ref is a supported local component reference
// and returns the trailing component name.
func parseRef(ref string) (string, error) {
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return "", fmt.Errorf("unsupported or invalid $ref syntax: %s", ref)
	}
	name := strings.TrimPrefix(ref, prefix)
	if name == "" {
		return "", fmt.Errorf("invalid $ref syntax: %s", ref)
	}
	return name, nil
}
