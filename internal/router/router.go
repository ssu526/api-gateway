package router

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ssu526/api-gateway/internal/config"
)

// Router holds a list of routes used to match incoming HTTP requests to backend services
type Router struct {
	routes []Route
}

func New(routes []Route) *Router {
	return &Router{
		routes: routes,
	}
}

// Returns the first route that matches both the request path and HTTP method
func (r *Router) FindRoute(req *http.Request) (*Route, bool) {
	reqPath := req.URL.Path

	for i := range r.routes {
		if strings.HasPrefix(reqPath, r.routes[i].Path) {
			if len(r.routes[i].Methods) == 0 {
				return &r.routes[i], true
			}
			if slices.Contains(r.routes[i].Methods, req.Method) {
				return &r.routes[i], true
			}
		}
	}
	return nil, false
}

func BuildRoutes(cfgRoutes []config.Route) []Route {
	routes := make([]Route, 0, len(cfgRoutes))

	for i := range cfgRoutes {
		routes = append(routes, Route{
			Path:        cfgRoutes[i].Path,
			ServiceName: cfgRoutes[i].ServiceName,
			Methods:     cfgRoutes[i].Methods,
			Middlewares: cfgRoutes[i].Middlewares,
			Resilience:  buildResiliency(cfgRoutes[i].Resilience),
		})
	}
	return routes
}

func buildResiliency(r config.Resilience) Resilience {
	p := Resilience{
		Timeout:      time.Duration(r.TimeoutMs) * time.Microsecond,
		RetryCount:   r.RetryCount,
		RetryBackoff: time.Duration(r.RetryBackoffMs) * time.Millisecond,
	}

	if r.RateLimit != nil {
		p.RateLimiter = &RateLimiter{
			RequestPerSecond: r.RateLimit.RequestPerSecond,
			Burst:            r.RateLimit.Burst,
			Per:              parseRateLimitScope(r.RateLimit.Per),
		}
	}

	if r.CircuitBreaker != nil {
		p.CircuitBreaker = &CircuitBreaker{
			FailureThreshold: r.CircuitBreaker.FailureThreshold,
			RecoveryTimeout:  time.Duration(r.CircuitBreaker.RecoveryTimeoutMs) * time.Millisecond,
		}
	}
	return p
}

func parseRateLimitScope(s string) RateLimitScope {
	switch s {
	case "user":
		return RateLimit_User
	case "api-key":
		return RateLimit_APIKey
	default:
		return RateLimit_IP
	}
}
