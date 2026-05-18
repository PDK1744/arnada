package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PDK1744/gogateway/internal/config"
	"github.com/PDK1744/gogateway/internal/proxy"
	"github.com/PDK1744/gogateway/internal/router"
)

func StartServer() {
	cfg, err := config.LoadConfig("/home/kobeb/KobeCodes/gogateway/config.yaml")
	if err != nil {
		panic(err)
	}

	router, err := router.NewRouter(cfg)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		upstream, _ := router.Match(r)
		if upstream == "" {
			http.Error(w, "Response not found", http.StatusNotFound)
		}
		proxy, err := proxy.Proxy(upstream)
		if err != nil {
			log.Fatal("PROXY ERROR: ", err)
		}
		proxy.ServeHTTP(w, r)
	})
	fmt.Println("Gateway listening on: ", cfg.Server.Listen)
	http.ListenAndServe(cfg.Server.Listen, nil)
}
