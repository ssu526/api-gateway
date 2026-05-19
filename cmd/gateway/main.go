package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ssu526/api-gateway/internal/config"
	"github.com/ssu526/api-gateway/internal/gateway"
	"github.com/ssu526/api-gateway/internal/proxy"
	"github.com/ssu526/api-gateway/internal/router"
	"github.com/ssu526/api-gateway/internal/service"
)

func main() {
	// load config file
	cfg, err := config.Load("../../internal/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	fmt.Println(cfg)

	// build routes
	routes := router.BuildRoutes(cfg.Router.Routes)
	r := router.New(routes)

	// Build service registry and proxy
	serviceReg := service.New()
	p := proxy.New()
	for _, svcCfg := range cfg.ServiceRegistry.Services {
		svc := &service.Service{
			Name: svcCfg.Name,
			Urls: svcCfg.Urls,
		}
		for i := 0; i < len(svc.Urls); i++ {
			p.GetOrCreateProxy(svc.Urls[i])
		}
		serviceReg.RegisterService(svc)
	}

	gw := gateway.New(cfg.Port, r, serviceReg, p)

	log.Printf("API Gateway booting up on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(gw.Port, gw))
}
