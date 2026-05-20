package service

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	LbRoundRobin = "round-robin"
	LbLeastConn  = "least-conn"
)

var validStrategies = map[string]bool{
	LbRoundRobin: true,
	LbLeastConn:  true,
}

type Instance struct {
	URL   string
	State *InstanceState
}

type InstanceState struct {
	Healthy    atomic.Bool
	ActiveConn atomic.Int64
}

type Service struct {
	Name                string
	Instances           []*Instance
	LbStrategy          string
	HealthCheckPath     string
	HealthCheckInterval time.Duration
}

func (s *Service) validate() error {
	if s.Name == "" {
		return errors.New("service name must not be empty")
	}

	if len(s.Instances) == 0 {
		return fmt.Errorf("service %q must have at least one instance", s.Name)
	}

	for i, inst := range s.Instances {
		if inst == nil || inst.URL == "" {
			return fmt.Errorf("service %q: instance %d has an empty URL", s.Name, i)
		}
	}

	if s.LbStrategy != "" {
		s.LbStrategy = LbRoundRobin
	}

	if !validStrategies[s.LbStrategy] {
		return fmt.Errorf("service %q: unknown LB strategy %q", s.Name, s.LbStrategy)
	}

	return nil
}
