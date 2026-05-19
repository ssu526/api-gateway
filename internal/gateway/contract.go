package gateway

import (
	"net/http"
	"net/http/httputil"

	"github.com/ssu526/api-gateway/internal/router"
	"github.com/ssu526/api-gateway/internal/service"
)

type Router interface {
	FindRoute(*http.Request) (*router.Route, bool)
}

type ServiceRegistry interface {
	GetService(string) (*service.Service, error)
}

type Proxy interface {
	GetOrCreateProxy(string) (*httputil.ReverseProxy, error)
	Forward(http.ResponseWriter, *http.Request, string) error
}
