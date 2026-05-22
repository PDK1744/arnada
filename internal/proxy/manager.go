package proxy

import (
	"fmt"
	"net/http/httputil"

	"github.com/PDK1744/gogateway/internal/config"
)

type ProxyManager struct {
	proxies map[string]*httputil.ReverseProxy
}

func NewProxyManager(cfg *config.Config) (*ProxyManager, error) {
	proxies := make(map[string]*httputil.ReverseProxy)
	for _, r := range cfg.Routes {
		for _, path := range r.Paths {
			if _, ok := proxies[path.Upstream]; !ok {
				proxy, err := NewProxy(path.Upstream)
				if err != nil {
					return nil, err
				}
				proxies[path.Upstream] = proxy
			}
		}
	}
	return &ProxyManager{proxies: proxies}, nil
}

func (p *ProxyManager) GetProxy(upstream string) (*httputil.ReverseProxy, error) {
	proxUp, ok := p.proxies[upstream]
	if !ok {
		return nil, fmt.Errorf("Invalid or missing upstream: %v", upstream)
	}
	return proxUp, nil
}
