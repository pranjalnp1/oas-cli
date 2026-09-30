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

// exampleValueFor returns a placeholder value for a given OpenAPI schema type.
func exampleValueFor(schemaType string) interface{} {
	switch schemaType {
	case "integer":
		return 0
	case "number":
		return 0
	case "boolean":
		return true
	case "array":
		return []interface{}{}
	default:
		return "string"
	}
}

// exampleBody generates example JSON data from a schema's properties.
func exampleBody(spec *api_types.OpenAPISpec, schema map[string]interface{}) (map[string]interface{}, error) {
	schema, err := resolver.ResolveSchema(spec, schema)
	if err != nil {
		return nil, err
	}

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
		propType, _ := prop["type"].(string)
		body[name] = exampleValueFor(propType)
	}
	return body, nil
}

// Curl loads and parses an OpenAPI spec and assembles an example request for one operation.
func Curl(file string, method string, path string) (*CurlResult, error) {
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

	if len(spec.Servers) == 0 {
		return nil, fmt.Errorf("no servers defined in spec")
	}
	baseURL := spec.Servers[0].URL

	resolvedPath := path
	var queryParams []string
	for _, p := range op.Parameters {
		example := fmt.Sprintf("%v", exampleValueFor(p.Schema.Type))
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
