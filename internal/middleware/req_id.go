package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func ReqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", reqID)

		rc := &RequestContext{
			RequestID:   reqID,
			Method:      r.Method,
			Path:        r.URL.Path,
			Host:        r.Host,
			StartTime:   time.Now(),
			ReqBodySize: int(r.ContentLength),
		}
		ctx := WithRequestContext(r.Context(), rc)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
