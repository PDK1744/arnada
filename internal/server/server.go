package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PDK1744/gogateway/internal/config"
	"github.com/PDK1744/gogateway/internal/middleware"
	"github.com/PDK1744/gogateway/internal/proxy"
	"github.com/PDK1744/gogateway/internal/router"
)

func StartServer() {
	cfg, err := config.LoadConfig("/home/kobeb/KobeCodes/gogateway/config.yaml")
	if err != nil {
		panic(err)
	}

	rtr, err := router.NewRouter(cfg)
	if err != nil {
		panic(err)
	}

	proxyManager, err := proxy.NewProxyManager(cfg)
	if err != nil {
		panic(err)
	}
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstream, _ := rtr.Match(r)
		if upstream == "" {
			http.Error(w, "Response not found", http.StatusNotFound)
			return
		}

		proxy, err := proxyManager.GetProxy(upstream)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		proxy.Transport = &http.Transport{
			DisableKeepAlives: true,
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Println("PROXY ERROR:", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		log.Println("HOST:", r.Host)
		log.Println("UPSTREAM:", upstream)

		log.Println("ABOUT TO CALL PROXY")
		proxy.ServeHTTP(w, r)
		log.Println("PROXY RETURNED")
	})

	wrapped := middleware.BuildChain(finalHandler, middleware.ReqId, middleware.Logger)

	fmt.Println("Gateway listening on: ", cfg.Server.Listen)
	http.ListenAndServe(cfg.Server.Listen, wrapped)
}
