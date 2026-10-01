package formatter

import (
	"fmt"
	"io"

	"github.com/yourusername/oas-cli/internal/service"
)

// orNA returns s, or "N/A" if s is empty. Used for fields that are optional
// in the spec but should always render something for the user.
func orNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

// FormatInspect renders an InspectResult as the CLI's inspect summary block.
func FormatInspect(w io.Writer, r *service.InspectResult) {
	fmt.Fprintf(w, "OpenApiVersion : %s\n", r.OpenAPIVersion)
	fmt.Fprintf(w, "Title: %s\n", orNA(r.Title))
	fmt.Fprintf(w, "Version: %s\n", orNA(r.Version))
	fmt.Fprintf(w, "Description: %s\n", orNA(r.Description))
	fmt.Fprintln(w, "-------------------------")
	fmt.Fprintf(w, "Endpoints:  %d\n", r.Endpoints)
	fmt.Fprintf(w, "Operations:  %d\n", r.Operations)
	fmt.Fprintf(w, "Schemas:  %d\n", r.Schemas)
	fmt.Fprintln(w, "------ Servers -------")
	for _, url := range r.Servers {
		fmt.Fprintln(w, "-", url)
	}
}

// FormatList renders a ListResult as one "METHOD  path" line per operation.
func FormatList(w io.Writer, r *service.ListResult) {
	for _, op := range r.Operations {
		fmt.Fprintf(w, "%-7s %s\n", op.Method, op.Path)
	}
}

// FormatShow renders a ShowResult's summary, parameters, and responses.
func FormatShow(w io.Writer, r *service.ShowResult) {
	fmt.Fprintf(w, "%s %s\n\n", r.Method, r.Path)

	fmt.Fprintln(w, "Summary:")
	fmt.Fprintln(w, orNA(r.Summary))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Parameters:")
	if len(r.Parameters) == 0 {
		fmt.Fprintln(w, "(none)")
	} else {
		for _, p := range r.Parameters {
			required := "optional"
			if p.Required {
				required = "required"
			}
			fmt.Fprintf(w, "- %s (%s, %s, %s)\n", p.Name, p.In, required, orNA(p.Type))
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Responses:")
	for _, resp := range r.Responses {
		if resp.SchemaName == "" {
			fmt.Fprintf(w, "%s (no content)\n", resp.Code)
		} else {
			fmt.Fprintf(w, "%s %s\n", resp.Code, resp.SchemaName)
		}
	}
}

// FormatCurl renders a CurlResult as a runnable multi-line curl command.
func FormatCurl(w io.Writer, r *service.CurlResult) {
	fmt.Fprintln(w, "curl \\")
	fmt.Fprintf(w, "  -X %s \\\n", r.Method)

	if r.Body == "" {
		fmt.Fprintf(w, "  \"%s\"\n", r.URL)
		return
	}

	fmt.Fprintf(w, "  \"%s\" \\\n", r.URL)
	for key, value := range r.Headers {
		fmt.Fprintf(w, "  -H \"%s: %s\" \\\n", key, value)
	}
	fmt.Fprintf(w, "  -d '%s'\n", r.Body)
}
