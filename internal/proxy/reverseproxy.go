package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/PDK1744/gogateway/internal/config"
	"github.com/PDK1744/gogateway/internal/middleware"
)

func NewProxy(route *config.RouteConfig) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(route.Upstream)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	ogDirector := proxy.Director

	proxy.Director = func(r *http.Request) {
		ogDirector(r)
		r.Header.Del("X-Forward-For")
		r.Header.Set("X-Forward-For", r.RemoteAddr)
		r.Header.Del("X-Real-IP")
		rc, _ := middleware.GetRequestContext(r.Context())
		rc.Upstream = string(r.URL.Host)
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if rc, ok := middleware.GetRequestContext(r.Context()); ok {
			rc.Error = err.Error()
		}
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Server")
		resp.Header.Del("X-Powered-By")
		return nil
	}

	return proxy, nil
}
