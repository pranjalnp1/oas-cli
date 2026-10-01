package service

import (
	"fmt"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/internal/parser"
	"github.com/yourusername/oas-cli/types"
)

// Inspect loads and parses an OpenAPI spec and computes its summary data.
func Inspect(file string) (*InspectResult, error) {
	data, err := loader.Load(file)
	if err != nil {
		return nil, err
	}

	spec, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	if spec.OpenAPI == "" {
		return nil, fmt.Errorf("invalid spec: missing required \"openapi\" field")
	}

	operationCount := 0
	for _, item := range spec.Paths {
		for _, op := range []*api_types.Operation{item.Get, item.Post, item.Put, item.Delete, item.Patch, item.Head, item.Options, item.Trace} {
			if op != nil {
				operationCount++
			}
		}
	}

	servers := make([]string, 0, len(spec.Servers))
	for _, s := range spec.Servers {
		servers = append(servers, s.URL)
	}

	return &InspectResult{
		OpenAPIVersion: spec.OpenAPI,
		Title:          spec.Info.Title,
		Version:        spec.Info.Version,
		Description:    spec.Info.Description,
		Endpoints:      len(spec.Paths),
		Operations:     operationCount,
		Schemas:        len(spec.Components.Schemas),
		Servers:        servers,
	}, nil
}
