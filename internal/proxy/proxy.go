package proxy

import (
	"net/http"
	"sync"

	"golang.org/x/sync/singleflight"
)

type Proxy struct {
	Proxies   sync.Map //[upstream url, reverse proxy]
	Transport *http.Transport
	Group     singleflight.Group
}
