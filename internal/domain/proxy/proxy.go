package proxy

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/ssu526/api-gateway/internal/config"
	"golang.org/x/sync/singleflight"
)

type Proxy struct {
	proxies   sync.Map //[upstream url, reverse proxy]
	transport *http.Transport
	group     singleflight.Group
}

func NewProxy(cfg *config.TransportConfig) *Proxy {
	return &Proxy{
		transport: NewSharedTransport(cfg),
	}
}

func (p *Proxy) getOrCreateProxy(upstream string) (*httputil.ReverseProxy, error) {
	if upstream == "" {
		return nil, errors.New("upstream URL must not be empty")
	}

	u, err := url.Parse(upstream)
	if err != nil {
		return nil, fmt.Errorf("invalid upstream URL %q: %w", upstream, err)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("upstream URL %q has no host", upstream)
	}

	targetBase := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}
	cacheKey := targetBase.String()

	// Fast path: proxy already exists, skip singleflight overhead
	if v, ok := p.proxies.Load(cacheKey); ok {
		return v.(*httputil.ReverseProxy), nil
	}

	// Slow path: deduplicate concurrent creation for the same upstream
	v, err, _ := p.group.Do(cacheKey, func() (interface{}, error) {
		// Re-check inside the group in case another goroutine created it
		// between the fast-path check above and acquiring the group slot
		if v, ok := p.proxies.Load(cacheKey); ok {
			return v, nil
		}

		proxy := &httputil.ReverseProxy{
			Transport: p.transport,
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				slog.Error("upstream error", "url", r.URL, "err", err)
				http.Error(w, "bad gateway", http.StatusBadGateway)
			},
			Rewrite: func(pr *httputil.ProxyRequest) {
				target := *targetBase                             // make a copy of the target
				pr.SetURL(&target)                                // Changes the delivery address
				pr.SetXForwarded()                                // service only knows the request coming from gateway, X-Forwarded-For: user ip, X-Forwarded-Proto: original protocol
				pr.Out.Host = pr.In.Host                          // original domain name requested by the user
				pr.Out.Header.Set("X-Forwarded-Host", pr.In.Host) // original domain name requested by the user
			},
		}
		p.proxies.Store(cacheKey, proxy)

		p.group.Forget(cacheKey)
		return proxy, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*httputil.ReverseProxy), nil
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request, upstream string) {
	proxy, err := p.getOrCreateProxy(upstream)

	if err != nil {
		slog.Error("failed to clear proxy routing", "upstream", upstream, "err", err)
		http.Error(w, "bad gateway configuration", http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
}
