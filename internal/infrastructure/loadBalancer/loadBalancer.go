package loadBalancer

import "github.com/ssu526/api-gateway/internal/domain/service"

type LoadBalancer interface {
	Next(service []*service.Instance) (*service.Instance, error)
}
