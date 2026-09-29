package parser

import (
	"encoding/json"
	"fmt"

	"github.com/yourusername/oas-cli/types"
	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*api_types.OpenAPISpec, error) {
	var spec api_types.OpenAPISpec

	if json.Valid(data) {
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, fmt.Errorf("unable to parse specification as JSON: %w", err)
		}
		return &spec, nil
	}

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("unable to parse specification as YAML: %w", err)
	}
	return &spec, nil
}
