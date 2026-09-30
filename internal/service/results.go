package service

// InspectResult holds the computed summary data for the inspect command.
type InspectResult struct {
	OpenAPIVersion string
	Title          string
	Version        string
	Description    string
	Endpoints      int
	Operations     int
	Schemas        int
	Servers        []string
}

// OperationSummary is a single method+path entry for the list command.
type OperationSummary struct {
	Method string
	Path   string
}

// ListResult holds every operation defined in the spec, in display order.
type ListResult struct {
	Operations []OperationSummary
}

// ParameterView is a single parameter's display data for the show command.
type ParameterView struct {
	Name     string
	In       string
	Required bool
	Type     string
}

// ResponseView is a single response's display data for the show command.
type ResponseView struct {
	Code       string
	SchemaName string
}

// ShowResult holds the computed detail data for one operation.
type ShowResult struct {
	Method     string
	Path       string
	Summary    string
	Parameters []ParameterView
	Responses  []ResponseView
}

// CurlResult holds the assembled pieces of an example curl request.
type CurlResult struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string // pretty-printed JSON, empty if the operation has no body
}
