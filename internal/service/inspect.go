package service

import (
	"fmt"
	"io"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/types"
	"gopkg.in/yaml.v3"
)

// Helper function for empty field logic for non-crucial fields.
func orNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

// Inspect function for inspecting OpenApi files.
func Inspect(File string, writer io.Writer) error {
	// load file
	// parse spec
	// calculate information
	// print output
	data, err := loader.Load(File)
	if err != nil {
		return err
	}

	var spec api_types.OpenAPISpec
	err = yaml.Unmarshal(data, &spec)

	if err != nil {
		return err
	}

	// logic and error handling for spec.OpenAPI in case file is not a OpenAPI spec.
	var OpenApiVersion = spec.OpenAPI
	if OpenApiVersion == "" {
		return fmt.Errorf("invalid spec: missing required \"openapi\" field")
	}
	fmt.Fprintln(writer, "OpenApiVersion : ", spec.OpenAPI)

	// output for title.
	var title = spec.Info.Title
	fmt.Fprintf(writer, "Title: %s\n", orNA(title))

	// output for version.
	var version = spec.Info.Version
	fmt.Fprintf(writer, "Version: %s\n", orNA(version))

	// output for summary.
	var summary = spec.Operation.Summary
	fmt.Fprintf(writer, "Summary: %s\n", orNA(summary))

	fmt.Println("-------------------------")
	// output endpoints
	fmt.Fprintln(writer, "Endpoints: ", len(spec.Paths))

	// output for operations
	count := 0
	for _, item := range spec.Paths {
		for _, op := range []*api_types.Operation{item.Get, item.Post, item.Put, item.Delete, item.Patch, item.Head, item.Options, item.Trace} {
			if op != nil {
				count++
			}
		}
	}
	fmt.Fprintln(writer, "Operations: ", count)

	// output for schemas
	fmt.Fprintln(writer, "Schemas: ", len(spec.Components.Schemas))

	// output for servers
	fmt.Println("------ Servers -------")
	for _, s := range spec.Servers {
		fmt.Fprintln(writer, "-", s.URL)
	}

	return nil
}
