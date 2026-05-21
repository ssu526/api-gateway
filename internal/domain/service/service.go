package service

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ssu526/api-gateway/internal/err"
	"github.com/ssu526/api-gateway/internal/infrastructure/loadBalancer"
)

type Instance struct {
	URL   string
	State *InstanceState
}

type InstanceState struct {
	Healthy    atomic.Bool
	ActiveConn atomic.Int64
}

func (i *Instance) IsHealthy() bool {
	if i == nil || i.State == nil {
		return false
	}
	return i.State.Healthy.Load()
}

func (i *Instance) GetActiveConnections() int64 {
	if i == nil || i.State == nil {
		return 0
	}
	return i.State.ActiveConn.Load()
}

type Service struct {
	Name                string
	Instances           []*Instance
	LoadBalancer        loadBalancer.LoadBalancer
	HealthCheckPath     string
	HealthCheckInterval time.Duration
}

func (s *Service) validate() error {
	if s.Name == "" {
		return err.ErrEmptyServerName
	}

	if len(s.Instances) == 0 {
		return fmt.Errorf("service %q must have at least one instance", s.Name)
	}

	for i, inst := range s.Instances {
		if inst == nil || inst.URL == "" {
			return fmt.Errorf("service %q: instance %d has an empty URL", s.Name, i)
		}
	}

	return nil
}
