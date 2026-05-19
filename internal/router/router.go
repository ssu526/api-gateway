package router

import (
	"net/http"
	"slices"
	"strings"
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
		route := r.routes[i]
		if strings.HasPrefix(reqPath, route.Path) {
			if len(route.Methods) == 0 {
				return &r.routes[i], true
			}
			if slices.Contains(route.Methods, req.Method) {
				return &route, true
			}
		}
	}
	return nil, false
}
