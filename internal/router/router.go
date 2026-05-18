package router

import (
	"fmt"
	"net/http"

	"github.com/PDK1744/gogateway/internal/config"
)

type Router struct {
	routes []Route
}

type Route struct {
	Host     string
	Upstream string
}

func NewRouter(cfg *config.Config) (*Router, error) {
	var routes []Route
	for _, r := range cfg.Routes {
		route := &Route{Host: r.Host, Upstream: r.Upstream}
		routes = append(routes, *route)
	}
	return &Router{routes: routes}, nil
}

func (r *Router) Match(req *http.Request) (backend string, ok bool) {
	for _, r := range r.routes {
		if r.Host == req.Host {
			fmt.Println("Match found!")
			return r.Upstream, ok
		}
	}
	return "", ok
}
