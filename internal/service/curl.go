package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/internal/parser"
	"github.com/yourusername/oas-cli/internal/resolver"
	"github.com/yourusername/oas-cli/types"
)

const maxExampleDepth = 10

// exampleValueFor returns a placeholder value for a given OpenAPI primitive schema type.
func exampleValueFor(schemaType string) interface{} {
	switch schemaType {
	case "integer":
		return 0
	case "number":
		return 0
	case "boolean":
		return true
	default:
		return "string"
	}
}

// exampleValue generates an example value for an arbitrary schema, following
// $refs and recursing into object properties / array items as needed.
func exampleValue(spec *api_types.OpenAPISpec, schema map[string]interface{}, depth int) (interface{}, error) {
	if depth > maxExampleDepth {
		return nil, fmt.Errorf("schema nesting too deep or circular")
	}

	schema, err := resolver.ResolveSchema(spec, schema)
	if err != nil {
		return nil, err
	}

	if enum, ok := schema["enum"].([]interface{}); ok && len(enum) > 0 {
		return enum[0], nil
	}

	schemaType, _ := schema["type"].(string)

	switch schemaType {
	case "object":
		return exampleObject(spec, schema, depth)
	case "array":
		items, ok := schema["items"].(map[string]interface{})
		if !ok {
			return []interface{}{}, nil
		}
		item, err := exampleValue(spec, items, depth+1)
		if err != nil {
			return nil, err
		}
		return []interface{}{item}, nil
	case "":
		// No explicit type: if it has properties, treat it as an implicit object.
		if _, hasProps := schema["properties"]; hasProps {
			return exampleObject(spec, schema, depth)
		}
		return exampleValueFor(schemaType), nil
	default:
		return exampleValueFor(schemaType), nil
	}
}

// exampleObject generates an example JSON object from a schema's properties.
func exampleObject(spec *api_types.OpenAPISpec, schema map[string]interface{}, depth int) (map[string]interface{}, error) {
	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}

	body := make(map[string]interface{})
	for name, raw := range properties {
		prop, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		value, err := exampleValue(spec, prop, depth+1)
		if err != nil {
			return nil, err
		}
		body[name] = value
	}
	return body, nil
}

// exampleBody generates example JSON data from a request body's schema.
func exampleBody(spec *api_types.OpenAPISpec, schema map[string]interface{}) (map[string]interface{}, error) {
	value, err := exampleValue(spec, schema, 0)
	if err != nil {
		return nil, err
	}
	body, ok := value.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}
	return body, nil
}

// Curl loads and parses an OpenAPI spec and assembles an example request for one operation.
// If baseURLOverride is non-empty, it is used instead of the spec's first server URL —
// this is needed when the spec declares a relative server URL (e.g. "/api/v3"), since
// a relative server is only resolvable against wherever the document was originally
// served from, information a local spec file does not carry.
func Curl(file string, method string, path string, baseURLOverride string) (*CurlResult, error) {
	data, err := loader.Load(file)
	if err != nil {
		return nil, err
	}

	spec, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	item, ok := spec.Paths[path]
	if !ok {
		return nil, fmt.Errorf("path not found: %s", path)
	}

	op := operationFor(item, method)
	if op == nil {
		return nil, fmt.Errorf("method %s not defined for path %s", strings.ToUpper(method), path)
	}

	baseURL := baseURLOverride
	if baseURL == "" {
		if len(spec.Servers) == 0 {
			return nil, fmt.Errorf("no servers defined in spec")
		}
		baseURL = spec.Servers[0].URL
		if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
			return nil, fmt.Errorf("server URL %q is relative; pass --base-url to specify the host", baseURL)
		}
	}

	resolvedPath := path
	var queryParams []string
	for _, p := range op.Parameters {
		var example string
		if len(p.Schema.Enum) > 0 {
			example = p.Schema.Enum[0]
		} else {
			example = fmt.Sprintf("%v", exampleValueFor(p.Schema.Type))
		}
		switch p.In {
		case "path":
			resolvedPath = strings.ReplaceAll(resolvedPath, "{"+p.Name+"}", example)
		case "query":
			queryParams = append(queryParams, p.Name+"="+example)
		}
	}

	url := baseURL + resolvedPath
	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	result := &CurlResult{
		Method:  strings.ToUpper(method),
		URL:     url,
		Headers: map[string]string{},
	}

	if op.RequestBody == nil {
		return result, nil
	}

	content, ok := op.RequestBody.Content["application/json"]
	if !ok {
		return result, nil
	}

	body, err := exampleBody(spec, content.Schema)
	if err != nil {
		return nil, err
	}
	bodyJSON, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return nil, err
	}

	result.Headers["Content-Type"] = "application/json"
	result.Body = string(bodyJSON)

	return result, nil
}
