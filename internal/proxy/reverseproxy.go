package proxy

import (
	"net/http/httputil"
	"net/url"
)

func Proxy(backendURL string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return proxy, nil
}
