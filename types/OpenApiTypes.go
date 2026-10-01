package api_types

type Server struct {
	URL string `yaml:"url" json:"url"`
}
type OpenAPISpec struct {
	OpenAPI    string              `yaml:"openapi" json:"openapi"`
	Info       Info                `yaml:"info"`
	Servers    []Server            `yaml:"servers" json:"servers"`
	Paths      map[string]PathItem `yaml:"paths"`
	Components ComponentsObject    `yaml:"components"`
}
type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}
type PathItem struct {
	Get     *Operation `yaml:"get"`
	Post    *Operation `yaml:"post"`
	Put     *Operation `yaml:"put"`
	Delete  *Operation `yaml:"delete"`
	Patch   *Operation `yaml:"patch"`
	Head    *Operation `yaml:"head"`
	Options *Operation `yaml:"options"`
	Trace   *Operation `yaml:"trace"`
}

type Parameter struct {
	Name     string `yaml:"name"`
	In       string `yaml:"in"`
	Required bool   `yaml:"required"`
	Schema   struct {
		Type string   `yaml:"type"`
		Enum []string `yaml:"enum"`
	} `yaml:"schema"`
}
type Response struct {
	Description string `yaml:"description"`
	Content     map[string]struct {
		Schema map[string]interface{} `yaml:"schema"` // holds $ref
	} `yaml:"content"`
}
type Operation struct {
	Summary     string              `yaml:"summary"`
	Description string              `yaml:"description"`
	Parameters  []Parameter         `yaml:"parameters"`
	Responses   map[string]Response `yaml:"responses"`
	RequestBody *RequestBody        `yaml:"requestBody"`
}

type RequestBody struct {
	Required bool `yaml:"required"`
	Content  map[string]struct {
		Schema map[string]interface{} `yaml:"schema"`
	} `yaml:"content"`
}
type ComponentsObject struct {
	Schemas map[string]interface{} `yaml:"schemas"`
}
