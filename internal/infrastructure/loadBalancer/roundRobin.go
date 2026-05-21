package loadBalancer

import (
	"sync/atomic"

	"github.com/ssu526/api-gateway/internal/err"
)

type RoundRobin struct {
	counter atomic.Uint64
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (r *RoundRobin) Next(instances []Peer) (Peer, error) {
	n := len(instances)

	if n == 0 {
		return nil, err.ErrNoServerRegistered
	}

	baseIdx := r.counter.Add(1)

	for i := 0; i < n; i++ {
		targetIdx := (baseIdx + uint64(i)) % uint64(n)
		inst := instances[targetIdx]

		if inst.IsHealthy() {
			return inst, nil
		}
	}
	return nil, err.ErrNoHealthyServer
}
