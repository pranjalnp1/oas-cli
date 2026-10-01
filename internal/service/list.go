package service

import (
	"sort"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/internal/parser"
	"github.com/yourusername/oas-cli/types"
)

// List loads and parses an OpenAPI spec and returns every operation it defines.
func List(file string) (*ListResult, error) {
	data, err := loader.Load(file)
	if err != nil {
		return nil, err
	}

	spec, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(spec.Paths))
	for path := range spec.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var operations []OperationSummary
	for _, path := range paths {
		item := spec.Paths[path]
		methods := []struct {
			name string
			op   *api_types.Operation
		}{
			{"GET", item.Get}, {"POST", item.Post}, {"PUT", item.Put}, {"DELETE", item.Delete},
			{"PATCH", item.Patch}, {"HEAD", item.Head}, {"OPTIONS", item.Options}, {"TRACE", item.Trace},
		}
		for _, m := range methods {
			if m.op != nil {
				operations = append(operations, OperationSummary{Method: m.name, Path: path})
			}
		}
	}

	return &ListResult{Operations: operations}, nil
}
