package service

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/oas-cli/internal/loader"
	"github.com/yourusername/oas-cli/types"
	"gopkg.in/yaml.v3"
)

// refName extracts the trailing component name from a local $ref string,
// e.g. "#/components/schemas/Pets" -> "Pets".
func refName(ref string) string {
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}

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

func Show(File string, Method string, Path string, writer io.Writer) error {
	data, err := loader.Load(File)
	if err != nil {
		return err
	}

	var spec api_types.OpenAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return err
	}

	item, ok := spec.Paths[Path]
	if !ok {
		return fmt.Errorf("path not found: %s", Path)
	}

	op := operationFor(item, Method)
	if op == nil {
		return fmt.Errorf("method %s not defined for path %s", strings.ToUpper(Method), Path)
	}

	fmt.Fprintf(writer, "%s %s\n\n", strings.ToUpper(Method), Path)

	fmt.Fprintln(writer, "Summary:")
	fmt.Fprintln(writer, orNA(op.Summary))
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "Parameters:")
	if len(op.Parameters) == 0 {
		fmt.Fprintln(writer, "(none)")
	} else {
		for _, p := range op.Parameters {
			required := "optional"
			if p.Required {
				required = "required"
			}
			fmt.Fprintf(writer, "- %s (%s, %s, %s)\n", p.Name, p.In, required, orNA(p.Schema.Type))
		}
	}
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "Responses:")
	for code, resp := range op.Responses {
		schemaName := ""
		for _, content := range resp.Content {
			if ref, ok := content.Schema["$ref"].(string); ok {
				schemaName = refName(ref)
			}
		}
		if schemaName == "" {
			fmt.Fprintf(writer, "%s (no content)\n", code)
		} else {
			fmt.Fprintf(writer, "%s %s\n", code, schemaName)
		}
	}
	return nil
}
