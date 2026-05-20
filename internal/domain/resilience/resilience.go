package resilience

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	RPS   int
	Burst int
	Scope RateLimitScope
}

type CircuitBreaker struct {
	FailureThreshold int
	RecoveryTimeout  time.Duration
}

type Resilience struct {
	Timeout        time.Duration
	Retry          int
	RetryBackoff   time.Duration
	RateLimiter    *RateLimiter // nil = disabled
	CircuitBreaker *CircuitBreaker
}

type RateLimitScope int

func (r RateLimitScope) String() string {
	switch r {
	case RateLimit_IP:
		return "ip"
	case RateLimit_APIKey:
		return "api-key"
	default:
		return fmt.Sprintf("unknown(%d)", int(r))
	}
}

const (
	RateLimit_IP RateLimitScope = iota
	RateLimit_APIKey
)

func ParseRateLimitScope(s string) (RateLimitScope, error) {
	switch s {
	case "ip":
		return RateLimit_IP, nil
	case "api-key":
		return RateLimit_APIKey, nil
	default:
		return 0, fmt.Errorf("unknown rate limit scope %q", s)
	}
}
