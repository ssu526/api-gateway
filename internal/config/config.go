package config

import "time"

type Config struct {
	Port            string          `yaml:"port"`
	ReadTimeout     time.Duration   `yaml:"read-timeout"`
	WriteTimeout    time.Duration   `yaml:"write-timeout"`
	IdleTimeout     time.Duration   `yaml:"idle-timeout"`
	MaxHeaderBytes  int             `yaml:"max-header-bytes"`
	TransportConfig TransportConfig `yaml:"transport"`
	Router          Router          `yaml:"router"`
	ServiceRegistry ServiceRegistry `yaml:"service-registry"`
}

// service
type ServiceRegistry struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name                string        `yaml:"name"`
	Urls                []string      `yaml:"urls"`
	LbStrategy          string        `yaml:"lb-strategy"` // round robin, least connection
	HealthCheckPath     string        `yaml:"health-check-path"`
	HealthCheckInterval time.Duration `yaml:"health-interval"`
}

// Route and resiliency
type Router struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path        string     `yaml:"path"`
	ServiceName string     `yaml:"service-name"`
	Methods     []string   `yaml:"methods,omitempty"`
	Middlewares []string   `yaml:"middleware"`
	Resilience  Resilience `yaml:"resilience"`
}

type RateLimit struct {
	RPS   int    `yaml:"rps"`
	Burst int    `yaml:"burst"`
	Scope string `yaml:"scope"` // ip, api-key
}

type CircuitBreaker struct {
	FailureThreshold int           `yaml:"failure-threshold"`
	RecoveryTimeout  time.Duration `yaml:"recovery-timeout"`
}

type Resilience struct {
	RateLimit      *RateLimit      `yaml:"rate-limit,omitempty"`
	Timeout        time.Duration   `yaml:"timeout"`
	RetryCount     int             `yaml:"retry-count"`
	RetryBackoff   time.Duration   `yaml:"retry-backoff"`
	CircuitBreaker *CircuitBreaker `yaml:"circuit-breaker"`
}

// Transport
type TransportConfig struct {
	DialTimeout           time.Duration `yaml:"dialTimeout"`
	KeepAlive             time.Duration `yaml:"keepAlive"`
	MaxIdleConns          int           `yaml:"maxIdleConns"`
	MaxIdleConnsPerHost   int           `yaml:"maxIdleConnsPerHost"`
	IdleConnTimeout       time.Duration `yaml:"idleConnTimeout"`
	TLSHandshakeTimeout   time.Duration `yaml:"tlsHandshakeTimeout"`
	ExpectContinueTimeout time.Duration `yaml:"expectedContinueTimeout"`
}
