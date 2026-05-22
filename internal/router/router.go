package router

import (
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/PDK1744/gogateway/internal/config"
)

type SortedPaths []config.PathConfig

type Router struct {
	hostRoutes map[string]SortedPaths
}

func NewRouter(cfg *config.Config) *Router {
	rtr := &Router{
		hostRoutes: make(map[string]SortedPaths),
	}
	for _, route := range cfg.Routes {
		paths := SortedPaths(route.Paths)

		sort.SliceStable(paths, func(i, j int) bool {
			return len(paths[i].Path) > len(paths[j].Path)
		})

		hostKey := strings.ToLower(strings.TrimSpace(route.Host))
		rtr.hostRoutes[hostKey] = paths
	}
	return rtr
}

func (r *Router) Match(req *http.Request) (upstream string, ok bool) {
	host := req.Host
	if strings.Contains(host, ":") {
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
	}
	host = strings.ToLower(host)

	paths, hostExists := r.hostRoutes[host]
	if !hostExists {
		return "", false
	}
	for _, p := range paths {
		if strings.HasPrefix(req.URL.Path, p.Path) {
			return p.Upstream, true
		}
	}
	return "", false
}
