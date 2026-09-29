package service

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/internal/parser"
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

// resolveSchema follows a single local $ref (e.g. "#/components/schemas/Pet")
// into spec.Components.Schemas, returning the referenced schema object.
func resolveSchema(spec *api_types.OpenAPISpec, schema map[string]interface{}) map[string]interface{} {
	ref, ok := schema["$ref"].(string)
	if !ok {
		return schema
	}
	name := refName(ref)
	resolved, ok := spec.Components.Schemas[name].(map[string]interface{})
	if !ok {
		return schema
	}
	return resolved
}

// exampleBody generates example JSON data from a schema's properties.
func exampleBody(spec *api_types.OpenAPISpec, schema map[string]interface{}) map[string]interface{} {
	schema = resolveSchema(spec, schema)

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
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
	return body
}

func Curl(filename string, method string, path string, writer io.Writer) error {
	data, err := loader.Load(filename)
	if err != nil {
		return err
	}
	spec, err := parser.Parse(data)
	if err != nil {
		return err
	}

	item, ok := spec.Paths[path]
	if !ok {
		return fmt.Errorf("path not found: %s", path)
	}

	op := operationFor(item, method)
	if op == nil {
		return fmt.Errorf("method %s not defined for path %s", strings.ToUpper(method), path)
	}

	if len(spec.Servers) == 0 {
		return fmt.Errorf("no servers defined in spec")
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

	fmt.Fprintln(writer, "curl \\")
	fmt.Fprintf(writer, "  -X %s \\\n", strings.ToUpper(method))

	if op.RequestBody == nil {
		fmt.Fprintf(writer, "  \"%s\"\n", url)
		return nil
	}

	content, ok := op.RequestBody.Content["application/json"]
	if !ok {
		fmt.Fprintf(writer, "  \"%s\"\n", url)
		return nil
	}

	body := exampleBody(spec, content.Schema)
	bodyJSON, err := json.MarshalIndent(body, "    ", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "  \"%s\" \\\n", url)
	fmt.Fprintln(writer, "  -H \"Content-Type: application/json\" \\")
	fmt.Fprintf(writer, "  -d '%s'\n", bodyJSON)

	return nil
}
