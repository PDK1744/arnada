package proxy

import (
	"testing"

	"github.com/PDK1744/gogateway/internal/config"
)

func TestProxyManager_Deduplication(t *testing.T) {
	cfg := &config.Config{
		Routes: []config.RouteConfig{
			{
				Host: "a.com",
				Paths: []config.PathConfig{
					{Path: "/api", Upstream: "http://shared-backend:8080"},
				},
			},
			{
				Host: "b.com",
				Paths: []config.PathConfig{
					{Path: "/", Upstream: "http://shared-backend:8080"}, // Duplicated upstream target
					{Path: "/static", Upstream: "http://static-backend:9000"},
				},
			},
		},
	}

	pm, err := NewProxyManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create ProxyManager: %v", err)
	}

	// Assert that even though there is 3 total route paths,
	// only initialized exactly 2 unique reverse proxy pools.
	expectedCount := 2
	if len(pm.proxies) != expectedCount {
		t.Errorf("ProxyManager cached %d proxies, want %d (deduplication failed)", len(pm.proxies), expectedCount)
	}
}
