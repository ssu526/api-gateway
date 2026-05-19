package service

import "fmt"

type ServiceRegistry struct {
	services map[string]*Service
}

func New() *ServiceRegistry {
	m := make(map[string]*Service)
	return &ServiceRegistry{
		services: m,
	}
}

func (s *ServiceRegistry) RegisterService(service *Service) {
	s.services[service.Name] = service
}

func (s *ServiceRegistry) GetService(serviceName string) (*Service, error) {
	service, ok := s.services[serviceName]

	if !ok {
		return nil, fmt.Errorf("Service %s not found", serviceName)
	}

	return service, nil
}
