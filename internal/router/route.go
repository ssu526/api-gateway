package router

import "time"

type Route struct {
	Path        string
	ServiceName string
	Methods     []string
	Middlewares []string
	Resilience  Resilience
}

type Resilience struct {
	Timeout        time.Duration
	RetryCount     int
	RetryBackoff   time.Duration
	RateLimiter    *RateLimiter
	CircuitBreaker *CircuitBreaker
}

type RateLimiter struct {
	RequestPerSecond int
	Burst            int
	Per              RateLimitScope
}

type CircuitBreaker struct {
	FailureThreshold int
	RecoveryTimeout  time.Duration
}

type RateLimitScope int

const (
	RateLimit_IP RateLimitScope = iota
	RateLimit_User
	RateLimit_APIKey
)
