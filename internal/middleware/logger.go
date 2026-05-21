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
			reqCtx := statWriter.RequestContext
			if reqCtx == nil {
				fmt.Println("RequestContext not found on StatWriter")
				return
			}

			// sync data into the context object
			duration := time.Since(reqCtx.StartTime).Milliseconds()
			reqCtx.SetResponse(statWriter.Status, statWriter.Bytes, duration)

			rcJson, _ := json.Marshal(reqCtx)
			fmt.Println(string(rcJson))
		}()

		next.ServeHTTP(statWriter, r)

	})
}
