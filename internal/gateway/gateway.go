package gateway

import (
	"net/http"

	"github.com/ssu526/api-gateway/internal/domain/proxy"
	"github.com/ssu526/api-gateway/internal/domain/route"
	"github.com/ssu526/api-gateway/internal/domain/service"
)

type Gateway struct {
	Port            string
	ReadTimeoutMs   int
	WriteTimeoutMs  int
	IdleTimeoutMs   int
	MaxHeaderBytes  int
	Router          route.Router
	ServiceRegistry service.ServiceRegistry
	Proxy           proxy.ReverseProxy
}

func NewGateway(
	port string,
	router route.Router,
	serviceRegistry service.ServiceRegistry,
	proxy proxy.ReverseProxy,
) *Gateway {
	return &Gateway{
		Port:            port,
		Router:          router,
		ServiceRegistry: serviceRegistry,
		Proxy:           proxy,
	}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := g.Router.FindRoute(r)

	if !ok {
		http.NotFound(w, r)
		return
	}

	svc, err := g.ServiceRegistry.GetService(route.ServiceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if len(svc.Instances) == 0 {
		http.Error(w, "No upstream targets available", http.StatusBadGateway)
	}

	upstream := svc.Instances[0].URL
	g.Proxy.Forward(w, r, upstream)
}
