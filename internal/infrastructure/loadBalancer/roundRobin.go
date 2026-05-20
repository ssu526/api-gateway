package loadBalancer

// import (
// 	"errors"
// 	"sync/atomic"

// 	"github.com/ssu526/api-gateway/internal/domain/service"
// )

// type RoundRobin struct {
// 	counter atomic.Uint64
// }

// func NewRoundRobin() *RoundRobin {
// 	return &RoundRobin{}
// }

// func (r *RoundRobin) Next(instances []*service.Instance) (*service.Instance, error) {
// 	n := len(instances)

// 	if n == 0 {
// 		return nil, errors.New("no instances registered")
// 	}

// 	baseIdx := r.counter.Add(1)

// 	for i := 0; i < n; i++ {
// 		targetIdx := (baseIdx + uint64(i)%uint64(n))
// 		inst := instances[targetIdx]

// 		if inst.Healthy.Load() {
// 			return inst, nil
// 		}
// 	}
// 	return nil, errors.New("no healthy intances avialable")
// }
