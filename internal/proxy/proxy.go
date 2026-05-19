package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Proxy struct {
	Proxies   sync.Map //[upstream url, reverse proxy]
	Transport *http.Transport
	Group     singleflight.Group
}

func New() *Proxy {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   50,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}
	return &Proxy{
		Transport: transport,
	}
}

func (p *Proxy) GetOrCreateProxy(upstream string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}

	targetBase := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}
	cacheKey := targetBase.String()

	if v, ok := p.Proxies.Load(cacheKey); ok {
		return v.(*httputil.ReverseProxy), nil
	}

	v, err, _ := p.Group.Do(cacheKey, func() (interface{}, error) {
		if v, ok := p.Proxies.Load(cacheKey); ok {
			return v, nil
		}

		proxy := &httputil.ReverseProxy{
			Transport: p.Transport,
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				http.Error(w, "bad gateway", http.StatusBadGateway)
			},
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetURL(targetBase)                             // Sets up the routing destination
				pr.SetXForwarded()                                // updates X-Forwarded-For, X-Forwarded-Host, and X-Forwarded-Proto
				pr.Out.Header.Set("X-Forwarded-Host", pr.In.Host) // reset host to original host
			},
		}
		p.Proxies.Store(cacheKey, proxy)
		return proxy, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*httputil.ReverseProxy), nil
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request, upstream string) error {
	proxy, err := p.GetOrCreateProxy(upstream)

	if err != nil {
		return err
	}

	proxy.ServeHTTP(w, r)
	return nil
}
