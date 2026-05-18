package proxy

import (
	"fmt"
	"log"
	"net/http/httputil"

	"github.com/PDK1744/gogateway/internal/config"
)

type ProxyManager struct {
	proxies map[string]*httputil.ReverseProxy
}

func NewProxyManager(cfg *config.Config) (*ProxyManager, error) {
	proxies := make(map[string]*httputil.ReverseProxy)
	for _, r := range cfg.Routes {
		proxy, err := NewProxy(&r)
		if err != nil {
			return nil, err
		}
		proxies[r.Upstream] = proxy
	}
	return &ProxyManager{proxies: proxies}, nil
}

func (p *ProxyManager) GetProxy(upstream string) (*httputil.ReverseProxy, error) {
	proxUp, ok := p.proxies[upstream]
	if !ok {
		return nil, fmt.Errorf("Invalid or missing upstream: %v", upstream)
	}
	fmt.Println("Proxy found!")
	log.Println("GET PROXY FOR:", upstream)
	log.Printf("PROXY INSTANCE: %+v\n", p)
	return proxUp, nil
}
