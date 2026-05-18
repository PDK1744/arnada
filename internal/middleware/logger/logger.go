package logger

import (
	"log"
	"net/http"
	"time"

	"github.com/PDK1744/gogateway/internal/middleware"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("LOGGER HIT")
		statWriter := &middleware.StatWriter{ResponseWriter: w}
		start := time.Now()

		next.ServeHTTP(statWriter, r)
		log.Println("AFTER HANDLER")

		log.Printf(
			"Method: %s | Path: %s | Duration: %v | Status: %v | Bytes: %v | RemoteAddr: %v",
			r.Method, r.URL.Path, time.Since(start), statWriter.Status, statWriter.Bytes, r.RemoteAddr,
		)
	})
}
