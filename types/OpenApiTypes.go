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
	Operation  Operation           `yaml:"operation" json:"operation"`
}
type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}
type PathItem struct {
	Get    *Operation `yaml:"get"`
	Post   *Operation `yaml:"post"`
	Put    *Operation `yaml:"put"`
	Delete *Operation `yaml:"delete"`
	//Patch
	//Head
	//Options
	//Trace
}
type Operation struct {
	Summary string `yaml:"summary"`
	// add more fields as needed
}
type ComponentsObject struct {
	Schemas map[string]interface{} `yaml:"schemas"`
}
