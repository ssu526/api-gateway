package loadBalancer

// import (
// 	"errors"

// 	"github.com/ssu526/api-gateway/internal/domain/service"
// )

// type LeastConnection struct {
// }

// func NewLeastConnection() *LeastConnection {
// 	return &LeastConnection{}
// }

// func (l *LeastConnection) Next(instances []*service.Instance) (*service.Instance, error) {
// 	var selected *service.Instance

// 	for _, inst := range instances {
// 		if !inst.Healthy.Load() {
// 			continue
// 		}

// 		if selected == nil || inst.ActiveConn.Load() < selected.ActiveConn.Load() {
// 			selected = inst
// 		}
// 	}

// 	if selected == nil {
// 		return nil, errors.New("no healthy instances")
// 	}

// 	return selected, nil
// }
