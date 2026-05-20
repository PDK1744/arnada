package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statWriter := &StatWriter{ResponseWriter: w, Status: http.StatusOK}
		defer func() {
			rc, ok := GetRequestContext(r.Context())
			if !ok {
				fmt.Println("RequestContext not found")
				return
			}

			// sync data into the context object
			rc.Status = statWriter.Status
			rc.RespBodySize = int(statWriter.Bytes)
			rc.DurationMs = int64(time.Since(rc.StartTime).Milliseconds())

			rcJson, _ := json.Marshal(rc)
			fmt.Println(string(rcJson))
		}()

		next.ServeHTTP(statWriter, r)

	})
}
