package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ssu526/api-gateway/internal/config"
	"github.com/ssu526/api-gateway/internal/gateway"
	"github.com/ssu526/api-gateway/internal/proxy"
	"github.com/ssu526/api-gateway/internal/router"
	"github.com/ssu526/api-gateway/internal/service"
)

func main() {
	// Load config file
	loader := config.NewLoader()
	cfg, err := loader.Load("../../internal/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	// Build routes
	routes := router.BuildRoutes(cfg.Router.Routes)
	r := router.New(routes)

	// Build service registry and proxy
	serviceReg := service.New()
	p := proxy.New()
	for _, svcCfg := range cfg.ServiceRegistry.Services {
		svc := &service.Service{
			Name:                svcCfg.Name,
			Urls:                svcCfg.Urls,
			LbStrategy:          svcCfg.LbStrategy,
			HealthCheckPath:     svcCfg.HealthCheckPath,
			HealthCheckInterval: time.Duration(svcCfg.HealthCheckIntervalMs) * time.Microsecond,
		}
		for i := 0; i < len(svc.Urls); i++ {
			p.GetOrCreateProxy(svc.Urls[i])
		}
		serviceReg.RegisterService(svc)
	}

	// Creae gateway
	gw := gateway.New(cfg.Port,
		time.Duration(cfg.ReadTimeoutMs)*time.Microsecond,
		time.Duration(cfg.WriteTimeoutMs)*time.Microsecond,
		time.Duration(cfg.IdleTimeoutMs)*time.Microsecond,
		cfg.MaxHeaderBytes,
		r,
		serviceReg,
		p)

	log.Printf("API Gateway booting up on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(gw.Port, gw))
}
