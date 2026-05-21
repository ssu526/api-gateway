package loadBalancer

const (
	LbRoundRobin = "round-robin"
	LbLeastConn  = "least-conn"
)

func NewLoadBalancer(strategy string) LoadBalancer {
	if strategy == LbLeastConn {
		return NewLeastConnection()
	} else {
		return NewRoundRobin()
	}
}
