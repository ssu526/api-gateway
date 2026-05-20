package service

import (
	"fmt"
	"sync"
)

type Services struct {
	mu       sync.RWMutex
	services map[string]*Service // [name, service]
}

// compile-time interface verification ensuring all required methods are implemented.
var _ ServiceRegistry = (*Services)(nil)

func NewServiceRegistry() *Services {
	return &Services{
		services: make(map[string]*Service),
	}
}

func (s *Services) RegisterService(service *Service) error {
	if err := service.validate(); err != nil {
		return fmt.Errorf("invalid service: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.services[service.Name]; exist {
		return fmt.Errorf("service already registered: %s", service.Name)
	}
	s.services[service.Name] = service
	return nil
}

func (s *Services) DeregisterService(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[name]; !exists {
		return fmt.Errorf("service %q not found", name)
	}

	delete(s.services, name)
	return nil
}

func (s *Services) GetService(serviceName string) (*Service, error) {
	// allow concurrent reads, reads are blocked if a writer holds the lock
	s.mu.RLock()

	service, ok := s.services[serviceName]
	if !ok {
		s.mu.RUnlock() // unlock early if not found
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	// Allocate a brand-new slice header and an entirely separate underlying array.
	// This stops downstream callers from structurally corrupting the central registry
	instances := make([]*Instance, len(service.Instances))
	for i, originalInstance := range service.Instances {
		instances[i] = &Instance{
			// copy-by-value string
			URL: originalInstance.URL,
			// Intentional pointer sharing: reassign the exact same State memory address.
			// This shares the underlying sync/atomic elements (Healthy, ActiveConn).
			State: originalInstance.State,
		}
	}

	name := service.Name
	strategy := service.LbStrategy
	healthCheckPath := service.HealthCheckPath
	healthCheckInterval := service.HealthCheckInterval

	s.mu.RUnlock() // release the lock immediately after we finish copying

	cp := &Service{
		Name:                name,
		Instances:           instances, // Bind our newly isolated slice array
		LbStrategy:          strategy,
		HealthCheckPath:     healthCheckPath,
		HealthCheckInterval: healthCheckInterval,
	}

	return cp, nil
}

/*
func (s *Services) GetService(serviceName string) (*Service, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    service, ok := s.services[serviceName]
    if !ok {
        return nil, fmt.Errorf("service %s not found", serviceName)
    }

    // Zero allocations. Returns the direct reference instantly.
    // Extremely fast, but relies on your developers not mutating service fields.
    return service, nil
}
*/
