package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/internal/parser"
	"github.com/yourusername/oas-cli/internal/resolver"
	"github.com/yourusername/oas-cli/types"
)

// operationFor returns the *Operation matching method on the given PathItem, if any.
func operationFor(item api_types.PathItem, method string) *api_types.Operation {
	switch strings.ToUpper(method) {
	case "GET":
		return item.Get
	case "POST":
		return item.Post
	case "PUT":
		return item.Put
	case "DELETE":
		return item.Delete
	case "PATCH":
		return item.Patch
	case "HEAD":
		return item.Head
	case "OPTIONS":
		return item.Options
	case "TRACE":
		return item.Trace
	default:
		return nil
	}
}

// Show loads and parses an OpenAPI spec and computes detail data for one operation.
func Show(file string, method string, path string) (*ShowResult, error) {
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

	parameters := make([]ParameterView, 0, len(op.Parameters))
	for _, p := range op.Parameters {
		parameters = append(parameters, ParameterView{
			Name:     p.Name,
			In:       p.In,
			Required: p.Required,
			Type:     p.Schema.Type,
		})
	}

	codes := make([]string, 0, len(op.Responses))
	for code := range op.Responses {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	responses := make([]ResponseView, 0, len(op.Responses))
	for _, code := range codes {
		resp := op.Responses[code]
		schemaName := ""
		for _, content := range resp.Content {
			if ref, ok := content.Schema["$ref"].(string); ok {
				name, err := resolver.ResolveName(spec, ref)
				if err != nil {
					return nil, err
				}
				schemaName = name
			}
		}
		responses = append(responses, ResponseView{Code: code, SchemaName: schemaName})
	}

	return &ShowResult{
		Method:     strings.ToUpper(method),
		Path:       path,
		Summary:    op.Summary,
		Parameters: parameters,
		Responses:  responses,
	}, nil
}
