package config

import "time"

func DefaultConfig() *Config {
	return &Config{
		Port:            "8080",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     120 * time.Second, // 2 minutes
		MaxHeaderBytes:  1 << 20,           // 1MB
		TransportConfig: DefaultTransportConfig(),
		Router: Router{
			Routes: []Route{},
		},
		ServiceRegistry: ServiceRegistry{
			Services: []Service{},
		},
	}
}

func DefaultServiceConfig() Service {
	return Service{
		LbStrategy:          "round-robin",
		HealthCheckPath:     "/healthz",
		HealthCheckInterval: 10 * time.Second,
	}
}

func DefaultResilienceConfig() Resilience {
	return Resilience{
		Timeout:      3 * time.Second,
		RetryCount:   3,
		RetryBackoff: 200 * time.Millisecond,
		CircuitBreaker: &CircuitBreaker{
			FailureThreshold: 5,
			RecoveryTimeout:  30 * time.Second,
		},
		// RateLimit left as nil by default. Let the user explicitly turn it on.
	}
}

func DefaultTransportConfig() TransportConfig {
	return TransportConfig{
		DialTimeout:           5 * time.Second,
		KeepAlive:             30 * time.Second,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   50,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
