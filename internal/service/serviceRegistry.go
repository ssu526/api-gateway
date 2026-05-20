package service

import (
	"fmt"
)

type ServiceRegistry struct {
	services map[string]*Service // [service name, Service instance]
}

func New() *ServiceRegistry {
	m := make(map[string]*Service)
	return &ServiceRegistry{
		services: m,
	}
}

func (s *ServiceRegistry) RegisterService(service *Service) error {
	if _, exist := s.services[service.Name]; exist {
		return fmt.Errorf("service already registered: %s", service.Name)
	}
	s.services[service.Name] = service
	return nil
}

func (s *ServiceRegistry) UpdateService(service *Service) error {
	if _, exist := s.services[service.Name]; !exist {
		return fmt.Errorf("service not found: %s", service.Name)
	}
	s.services[service.Name] = service
	return nil
}

func (s *ServiceRegistry) GetService(serviceName string) (*Service, error) {
	service, ok := s.services[serviceName]

	if !ok {
		return nil, fmt.Errorf("Service %s not found", serviceName)
	}

	return service, nil
}
