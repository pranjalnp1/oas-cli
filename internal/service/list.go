package service

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/types"
	"gopkg.in/yaml.v3"
)

func List(File string, writer io.Writer) error {
	data, err := loader.Load(File)
	if err != nil {
		return err
	}

	var spec api_types.OpenAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return err
	}

	paths := make([]string, 0, len(spec.Paths))
	for path := range spec.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		item := spec.Paths[path]
		methods := []struct {
			name string
			op   *api_types.Operation
		}{
			{"GET", item.Get}, {"POST", item.Post}, {"PUT", item.Put}, {"DELETE", item.Delete},
		}
		for _, m := range methods {
			if m.op != nil {
				fmt.Fprintf(writer, "%-7s %s\n", m.name, path)
			}
		}
	}
	return nil
}
