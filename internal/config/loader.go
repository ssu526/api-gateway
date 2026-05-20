package config

import (
	"errors"
	"fmt"
	"os"

	"go.yaml.in/yaml/v2"
)

func Load(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", filePath, err)
	}
	return parseConfig(data)
}

func parseConfig(data []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	applyDefaults(&config)

	if err := validate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func applyDefaults(c *Config) {
	defaults := DefaultConfig()
	serviceDefaults := DefaultServiceConfig()
	resilienceDefaults := DefaultResilienceConfig()
	transportDefault := DefaultTransportConfig()

	// top level config defaults
	if c.Port == "" {
		c.Port = defaults.Port
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaults.ReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaults.WriteTimeout
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = defaults.IdleTimeout
	}
	if c.MaxHeaderBytes == 0 {
		c.MaxHeaderBytes = defaults.MaxHeaderBytes
	}

	// transport defaults
	if c.TransportConfig.DialTimeout == 0 {
		c.TransportConfig.DialTimeout = transportDefault.DialTimeout
	}
	if c.TransportConfig.KeepAlive == 0 {
		c.TransportConfig.KeepAlive = transportDefault.KeepAlive
	}
	if c.TransportConfig.MaxIdleConns == 0 {
		c.TransportConfig.MaxIdleConns = transportDefault.MaxIdleConns
	}
	if c.TransportConfig.MaxIdleConnsPerHost == 0 {
		c.TransportConfig.MaxIdleConnsPerHost = transportDefault.MaxIdleConnsPerHost
	}
	if c.TransportConfig.IdleConnTimeout == 0 {
		c.TransportConfig.IdleConnTimeout = transportDefault.IdleConnTimeout
	}
	if c.TransportConfig.TLSHandshakeTimeout == 0 {
		c.TransportConfig.TLSHandshakeTimeout = transportDefault.TLSHandshakeTimeout
	}
	if c.TransportConfig.ExpectContinueTimeout == 0 {
		c.TransportConfig.ExpectContinueTimeout = transportDefault.ExpectContinueTimeout
	}

	// service level defaults
	for i := range c.ServiceRegistry.Services {
		s := &c.ServiceRegistry.Services[i]

		if s.LbStrategy == "" {
			s.LbStrategy = serviceDefaults.LbStrategy
		}
		if s.HealthCheckInterval == 0 {
			s.HealthCheckInterval = serviceDefaults.HealthCheckInterval
		}
		if s.HealthCheckPath == "" {
			s.HealthCheckPath = serviceDefaults.HealthCheckPath
		}
	}

	for i := range c.Router.Routes {
		r := &c.Router.Routes[i]

		if r.Resilience.Timeout == 0 {
			r.Resilience.Timeout = resilienceDefaults.Timeout
		}
		if r.Resilience.RetryCount == 0 {
			r.Resilience.RetryCount = resilienceDefaults.RetryCount
		}
		if r.Resilience.RetryBackoff == 0 {
			r.Resilience.RetryBackoff = resilienceDefaults.RetryBackoff
		}

		if r.Resilience.CircuitBreaker == nil {
			r.Resilience.CircuitBreaker = &CircuitBreaker{
				FailureThreshold: resilienceDefaults.CircuitBreaker.FailureThreshold,
				RecoveryTimeout:  resilienceDefaults.CircuitBreaker.RecoveryTimeout,
			}
		} else {
			if r.Resilience.CircuitBreaker.FailureThreshold == 0 {
				r.Resilience.CircuitBreaker.FailureThreshold = resilienceDefaults.CircuitBreaker.FailureThreshold
			}
			if r.Resilience.CircuitBreaker.RecoveryTimeout == 0 {
				r.Resilience.CircuitBreaker.RecoveryTimeout = resilienceDefaults.CircuitBreaker.RecoveryTimeout
			}
		}
	}
}

func validate(c *Config) error {
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
			return fmt.Errorf("service %q: must have at least one URL", s.Name)
		}
		serviceMap[s.Name] = true
	}

	for _, r := range c.Router.Routes {
		if r.Path == "" {
			return errors.New("route path cannot be empty")
		}
		if !serviceMap[r.ServiceName] {
			return fmt.Errorf("route %q: references undefined service %q", r.Path, r.ServiceName)
		}
	}

	return nil
}
