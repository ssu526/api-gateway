package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"go.yaml.in/yaml/v2"
)

type Loader struct {
	once sync.Once
	cfg  *Config
	err  error
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(filePath string) (*Config, error) {
	l.once.Do(func() {
		l.cfg, l.err = load(filePath)
	})

	return l.cfg, l.err
}

func load(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %s: %w", filePath, err)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("Failed to unmarshall yaml: %s: %w", filePath, err)
	}

	applyDefaults(&config)
	applyOverride(&config)

	if err := Validate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func applyDefaults(c *Config) {
	if c.Port == "" {
		c.Port = "8080"
	}
}

func applyOverride(c *Config) {
	if val := os.Getenv("GATEWAY_PORT"); val != "" {
		c.Port = val
	}
}

func Validate(c *Config) error {
	if c.Port == "" {
		return errors.New("port is required")
	}

	if len(c.Router.Routes) == 0 {
		return errors.New("at least one route must be defined")
	}

	serviceMap := make(map[string]bool)
	for _, s := range c.ServiceRegistry.Services {
		if s.Name == "" {
			return errors.New("service name cannot be empty")
		}
		if len(s.Urls) == 0 {
			return fmt.Errorf("service %s must have at least one URL", s.Name)
		}
		serviceMap[s.Name] = true
	}

	for _, r := range c.Router.Routes {
		if r.Path == "" {
			return errors.New("route path cannot be empty")
		}
		if !serviceMap[r.ServiceName] {
			return fmt.Errorf("route references undefined service: %s", r.ServiceName)
		}
	}

	return nil
}
