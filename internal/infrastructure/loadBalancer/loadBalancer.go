package loadBalancer

type Peer interface {
	IsHealthy() bool
	GetActiveConnections() int64
}

type LoadBalancer interface {
	Next(instances []Peer) (Peer, error)
}
