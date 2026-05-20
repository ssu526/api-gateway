package route

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/ssu526/api-gateway/internal/config"
	"github.com/ssu526/api-gateway/internal/domain/resilience"
)

type Routes struct {
	routes []*Route
}

type Route struct {
	Path        string
	ServiceName string
	Methods     []string
	Middlewares []string
	Resilience  resilience.Resilience
}

func NewRouter(routes []Route) Router {
	sorted := make([]*Route, len(routes))

	for i := range routes {
		sorted[i] = &Route{
			Path:        routes[i].Path,
			ServiceName: routes[i].ServiceName,
			Methods:     routes[i].Methods,
			Middlewares: routes[i].Middlewares,
			Resilience:  routes[i].Resilience,
		}
	}

	slices.SortFunc(sorted, func(a, b *Route) int {
		if len(a.Path) < len(b.Path) {
			return 1
		}
		if len(a.Path) > len(b.Path) {
			return -1
		}
		return 0
	})

	return &Routes{
		routes: sorted,
	}
}

func BuildRoutes(cfgRoutes []config.Route) ([]Route, error) {
	routes := make([]Route, 0, len(cfgRoutes))

	for i := range cfgRoutes {
		if err := validate(cfgRoutes[i]); err != nil {
			return nil, fmt.Errorf("invalid route at index %d: %w", i, err)
		}

		resil, err := buildResiliency(cfgRoutes[i].Resilience)
		if err != nil {
			return nil, fmt.Errorf("building resilience config for route %q: %w", cfgRoutes[i].Path, err)
		}
		routes = append(routes, Route{
			Path:        cfgRoutes[i].Path,
			ServiceName: cfgRoutes[i].ServiceName,
			Methods:     cfgRoutes[i].Methods,
			Middlewares: cfgRoutes[i].Middlewares,
			Resilience:  resil,
		})
	}
	return routes, nil
}

func buildResiliency(r config.Resilience) (resilience.Resilience, error) {
	p := resilience.Resilience{
		Timeout:      r.Timeout,
		Retry:        r.RetryCount,
		RetryBackoff: r.RetryBackoff,
	}

	if r.RateLimit != nil {
		scope, err := resilience.ParseRateLimitScope(r.RateLimit.Scope)
		if err != nil {
			return resilience.Resilience{}, fmt.Errorf("invalid rate limit scope: %w", err)
		}
		p.RateLimiter = &resilience.RateLimiter{
			RPS:   r.RateLimit.RPS,
			Burst: r.RateLimit.Burst,
			Scope: scope,
		}
	}

	if r.CircuitBreaker != nil {
		p.CircuitBreaker = &resilience.CircuitBreaker{
			FailureThreshold: r.CircuitBreaker.FailureThreshold,
			RecoveryTimeout:  r.CircuitBreaker.RecoveryTimeout,
		}
	}
	return p, nil
}

func validate(r config.Route) error {
	if r.Path == "" {
		return errors.New("route path must not be empty")
	}

	if !strings.HasPrefix(r.Path, "/") {
		return fmt.Errorf("route path %q must start with '/'", r.Path)
	}

	if r.ServiceName == "" {
		return fmt.Errorf("route %q: service name must not be empty", r.Path)
	}
	return nil
}

func (r *Routes) FindRoute(req *http.Request) (*Route, bool) {
	for _, route := range r.routes {
		if pathMatches(req.URL.Path, route.Path) {
			if len(route.Methods) == 0 {
				return route, true
			}
			if slices.Contains(route.Methods, req.Method) {
				return route, true
			}
		}
	}
	return nil, false
}

func pathMatches(reqPath, routePath string) bool {
	return reqPath == routePath || strings.HasPrefix(reqPath, strings.TrimRight(routePath, "/")+"/")
}
