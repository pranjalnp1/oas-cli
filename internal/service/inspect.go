package service

import (
	"fmt"
	"io"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/types"
	"gopkg.in/yaml.v3"
)

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
	fmt.Fprintln(writer, "API:", spec.Info.Title)

	for path := range spec.Paths {
		fmt.Fprintln(writer, " ", path)
	}

	fmt.Fprintln(writer, "OpenApiVersion : ", spec.OpenAPI)
	fmt.Fprintln(writer, "Title : ", spec.Info.Title)
	fmt.Fprintln(writer, "Version : ", spec.Info.Version)

	var summary = spec.Operation.Summary
	if summary == "" {
		fmt.Fprintln(writer, "Summary : N/A")
	} else {
		fmt.Fprintln(writer, "Summary : ", summary)
	}

	return nil

}
