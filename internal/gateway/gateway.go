package gateway

import (
	"net/http"
	"time"
)

type Gateway struct {
	Port            string
	ReadTimeoutMs   time.Duration
	WriteTimeoutMs  time.Duration
	IdleTimeoutMs   time.Duration
	MaxHeaderBytes  int
	Router          Router
	ServiceRegistry ServiceRegistry
	Proxy           Proxy
}

func New(port string,
	readTimeoutMs time.Duration,
	writeTimeoutMs time.Duration,
	idleTimeoutMs time.Duration,
	maxHeaderBytes int,
	router Router,
	serviceRegistry ServiceRegistry,
	proxy Proxy) *Gateway {
	return &Gateway{
		Port:            port,
		ReadTimeoutMs:   readTimeoutMs,
		WriteTimeoutMs:  writeTimeoutMs,
		IdleTimeoutMs:   idleTimeoutMs,
		MaxHeaderBytes:  maxHeaderBytes,
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

	if len(svc.Urls) == 0 {
		http.Error(w, "No upstream targets available", http.StatusBadGateway)
	}

	upstream := svc.Urls[0]
	err = g.Proxy.Forward(w, r, upstream)

	if err != nil {
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
}
