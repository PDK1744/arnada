package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/PDK1744/gogateway/internal/config"
)

func NewProxy(route *config.RouteConfig) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(route.Upstream)
	//log.Println("PARSED URL: ", url)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return proxy, nil
}
