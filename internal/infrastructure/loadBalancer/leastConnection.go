package loadBalancer

import (
	"github.com/ssu526/api-gateway/internal/err"
)

type LeastConnection struct {
}

func NewLeastConnection() *LeastConnection {
	return &LeastConnection{}
}

func (l *LeastConnection) Next(instances []Peer) (Peer, error) {
	var selected Peer

	for _, inst := range instances {
		if !inst.IsHealthy() {
			continue
		}

		if selected == nil || (inst.GetActiveConnections() < selected.GetActiveConnections()) {
			selected = inst
		}
	}

	if selected == nil {
		return nil, err.ErrNoHealthyServer
	}

	return selected, nil
}
