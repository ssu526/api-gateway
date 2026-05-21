package main

import (
	"log"
	"net/http"

	"github.com/ssu526/api-gateway/internal/config"
	"github.com/ssu526/api-gateway/internal/domain/proxy"
	"github.com/ssu526/api-gateway/internal/domain/route"
	"github.com/ssu526/api-gateway/internal/domain/service"
	"github.com/ssu526/api-gateway/internal/gateway"
	"github.com/ssu526/api-gateway/internal/infrastructure/loadBalancer"
)

func main() {
	// Load config file
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
		return
	}

	// Build routes
	routes, err := route.BuildRoutes(cfg.Router.Routes)
	if err != nil {
		log.Fatalf("Failed to build routes: %v", err)
	}
	router := route.NewRouter(routes)

	// Build service registry and proxy
	serviceReg := service.NewServiceRegistry()
	p := proxy.NewProxy(&cfg.TransportConfig)

	for _, svcCfg := range cfg.ServiceRegistry.Services {

		var instances []*service.Instance
		for _, urlStr := range svcCfg.Urls {
			inst := &service.Instance{
				URL:   urlStr,
				State: &service.InstanceState{},
			}
			inst.State.Healthy.Store(true)
			instances = append(instances, inst)
		}

		svc := &service.Service{
			Name:                svcCfg.Name,
			Instances:           instances,
			LoadBalancer:        loadBalancer.NewLoadBalancer(svcCfg.LbStrategy),
			HealthCheckPath:     svcCfg.HealthCheckPath,
			HealthCheckInterval: svcCfg.HealthCheckInterval,
		}

		if err := serviceReg.RegisterService(svc); err != nil {
			log.Fatalf("failed to register service %q: %v", svc.Name, err)
			return
		}
	}

	// Create gateway
	gw := gateway.NewGateway(cfg.Port, router, serviceReg, p)

	log.Printf("API Gateway booting up on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(gw.Port, gw))
}
