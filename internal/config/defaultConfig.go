package config

func DefaultConfig() *Config {
	return &Config{
		Port:           "8080",
		ReadTimeoutMs:  5000,
		WriteTimeoutMs: 10000,
		IdleTimeoutMs:  60000,
		MaxHeaderBytes: 1 << 20, // 1MB
		Router: Router{
			Routes: []Route{},
		},
		ServiceRegistry: ServiceRegistry{
			Services: []Service{},
		},
	}
}
