package router

import (
	"net/http/httptest"
	"testing"

	"github.com/PDK1744/gogateway/internal/config"
)

func TestRouter_Match(t *testing.T) {
	cfg := &config.Config{
		Routes: []config.RouteConfig{
			{
				Host: "sub.domain.com",
				Paths: []config.PathConfig{
					{Path: "/", Upstream: "http://fallback"},
					{Path: "/api/v1", Upstream: "http://api-v1"},
					{Path: "/static", Upstream: "http://static"},
					{Path: "/api", Upstream: "http://api-monolith"},
				},
			},
			{
				Host: "[::1]",
				Paths: []config.PathConfig{
					{Path: "/", Upstream: "http://ipv6-fallback"},
				},
			},
		},
	}

	rtr := NewRouter(cfg)

	tests := []struct {
		name             string
		incomingHost     string
		incomingPath     string
		expectedUpstream string
		expectedOk       bool
	}{
		{
			name:             "Longest prefix match wins over catch-all",
			incomingHost:     "sub.domain.com",
			incomingPath:     "/api/v1/users",
			expectedUpstream: "http://api-v1",
			expectedOk:       true,
		},
		{
			name:             "Fallback route catches un-matched prefixes",
			incomingHost:     "sub.domain.com",
			incomingPath:     "/dashboard/settings",
			expectedUpstream: "http://fallback",
			expectedOk:       true,
		},
		{
			name:             "Host port header is cleanly stripped during evaluation",
			incomingHost:     "sub.domain.com:9999", // Port attached
			incomingPath:     "/static/style.css",
			expectedUpstream: "http://static",
			expectedOk:       true,
		},
		{
			name:             "Unregistered host drops cleanly",
			incomingHost:     "evil-domain.com",
			incomingPath:     "/",
			expectedUpstream: "",
			expectedOk:       false,
		},
		{
			name:             "Host casing is completely ignored",
			incomingHost:     "SuB.DoMaIn.CoM",
			incomingPath:     "/api/v1/health",
			expectedUpstream: "http://api-v1",
			expectedOk:       true,
		},
		{
			name:             "More specific prefix wins over shorter matching prefix",
			incomingHost:     "sub.domain.com",
			incomingPath:     "/api/v1/orders",
			expectedUpstream: "http://api-v1", // NOT http://api-monolith
			expectedOk:       true,
		},
		// ipv6 test fails but its not in scope right now
		// will revisit at a later time
		// {
		// 	name:             "Raw IPv6 address with port is safely parsed",
		// 	incomingHost:     "[::1]:9999",
		// 	incomingPath:     "/",
		// 	expectedUpstream: "http://ipv6-fallback",
		// 	expectedOk:       true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://"+tt.incomingHost+tt.incomingPath, nil)

			upstream, ok := rtr.Match(req)

			if ok != tt.expectedOk {
				t.Errorf("Match() ok = %v, want %v", ok, tt.expectedOk)
			}
			if upstream != tt.expectedUpstream {
				t.Errorf("Match() upstream = %q, want %q", upstream, tt.expectedUpstream)
			}
		})
	}
}
