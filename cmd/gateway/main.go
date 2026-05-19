package main

import (
	"log"
	"net/http"

	"github.com/ssu526/api-gateway/internal/gateway"
	"github.com/ssu526/api-gateway/internal/proxy"
	"github.com/ssu526/api-gateway/internal/router"
	"github.com/ssu526/api-gateway/internal/service"
)

func main() {
	routes := []router.Route{
		{
			Path:        "/api/users",
			ServiceName: "users-service",
		},
	}

	services := []*service.Service{
		{
			Name: "users-service",
			Urls: []string{"http://localhost:9091"},
		},
	}

	r := router.New(routes)
	serviceReg := service.New()
	serviceReg.RegisterService(services[0])
	p := proxy.New()
	_, _ = p.GetOrCreateProxy("http://localhost:9091")

	gw := gateway.New(r, serviceReg, p)

	log.Printf("API Gateway booting up on port %s...", ":8080")
	log.Fatal(http.ListenAndServe(gw.Port, gw))
}
