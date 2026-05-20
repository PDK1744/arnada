package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("LOGGER HIT")
		statWriter := &StatWriter{ResponseWriter: w, Status: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(statWriter, r)
		log.Println("AFTER HANDLER")

		log.Printf(
			"Method: %s | Path: %s | Duration: %v | Status: %v | Bytes: %v | RemoteAddr: %v",
			r.Method, r.URL.Path, time.Since(start).Milliseconds(), statWriter.Status, statWriter.Bytes, r.RemoteAddr,
		)
	})
}
