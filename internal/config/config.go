package config

type Config struct {
	Port            string          `yaml:"port"`
	ReadTimeoutMs   int             `yaml:"read-timeout-ms"`
	WriteTimeoutMs  int             `yaml:"write-timeout-ms"`
	IdleTimeoutMs   int             `yaml:"idle-timeout-ms"`
	MaxHeaderBytes  int             `yaml:"max-header-bytes"`
	Router          Router          `yaml:"router"`
	ServiceRegistry ServiceRegistry `yaml:"service-registry"`
}

// service
type ServiceRegistry struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name                  string   `yaml:"service"`
	Urls                  []string `yaml:"urls"`
	LbStrategy            string   `yaml:"lb-strategy"`
	HealthCheckPath       string   `yaml:"health-check-path"`
	HealthCheckIntervalMs int      `yaml:"health-interval"`
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
	RequestPerSecond int    `yaml:"rps"`
	Burst            int    `yaml:"burst"`
	Per              string `yaml:"per"` // ip, user, api-key
}

type CircuitBreaker struct {
	FailureThreshold  int `yaml:"failure-threshold"`
	RecoveryTimeoutMs int `yaml:"recovery-timeout-ms"`
}

type Resilience struct {
	RateLimit      *RateLimit      `yaml:"rate-limit,omitempty"`
	TimeoutMs      int             `yaml:"timeout-ms"`
	RetryCount     int             `yaml:"retry-count"`
	RetryBackoffMs int             `yaml:"retry-backoff-ms"`
	CircuitBreaker *CircuitBreaker `yaml:"circuit-breaker"`
}
