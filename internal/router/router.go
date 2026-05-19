package router

import (
	"net/http"
	"slices"
	"strings"

	"github.com/ssu526/api-gateway/internal/config"
)

// Router holds a list of routes used to match incoming HTTP requests to backend services
type Router struct {
	routes []Route
}

// Creates a Router with the provided route definitions
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

	for i := 0; i < len(cfgRoutes); i++ {
		routes = append(routes, Route{
			Path:        cfgRoutes[i].Path,
			ServiceName: cfgRoutes[i].ServiceName,
			Methods:     cfgRoutes[i].Methods,
			Middlewares: cfgRoutes[i].Middlewares,
		})
	}
	return routes
}
