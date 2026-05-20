package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

		proxy.ServeHTTP(w, r)
	})

	wrapped := middleware.BuildChain(finalHandler, middleware.ReqId, middleware.Logger)

	fmt.Println("Gateway listening on: ", cfg.Server.Listen)
	server := &http.Server{Addr: cfg.Server.Listen, Handler: wrapped}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
