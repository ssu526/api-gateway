package config

type Config struct {
	Port            string          `yaml:"port"`
	Router          Router          `yaml:"router"`
	ServiceRegistry ServiceRegistry `yaml:"service-registry"`
}

type Router struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path        string   `yaml:"path"`
	ServiceName string   `yaml:"service-name"`
	Methods     []string `yaml:"methods,omitempty"`
	Middlewares []string `yaml:"middleware"`
}

type ServiceRegistry struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name string   `yaml:"service"`
	Urls []string `yaml:"urls"`
}
