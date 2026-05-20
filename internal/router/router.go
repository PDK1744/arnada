package router

import (
	"fmt"
	"net/http"

	"github.com/PDK1744/gogateway/internal/config"
)

type Router struct {
	routes []config.RouteConfig
}

func NewRouter(cfg *config.Config) (*Router, error) {
	var routes []config.RouteConfig
	for _, r := range cfg.Routes {
		route := &config.RouteConfig{Host: r.Host, Upstream: r.Upstream}
		routes = append(routes, *route)
	}
	return &Router{routes: routes}, nil
}

func (r *Router) Match(req *http.Request) (upstream string, ok bool) {
	for _, r := range r.routes {
		if r.Host == req.Host {
			fmt.Println("Match found!")
			return r.Upstream, ok
		}
	}
	return "", ok
}
