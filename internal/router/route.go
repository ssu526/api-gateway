package router

type Route struct {
	Path        string   `yaml:"path"`
	ServiceName string   `yaml:"service-name"`
	Methods     []string `yaml:"methods,omitempty"`
	Middlewares []string `yaml:"middleware"`
}
