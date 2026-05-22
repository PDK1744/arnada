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

		reqCtx, ok := GetRequestContext(r.Context())
		if !ok {
			return
		}
		reqCtx.SetRequest(reqID, r.Method, r.URL.Path, r.Host, time.Now(), int(r.ContentLength))

		ctx := WithRequestContext(r.Context(), reqCtx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ContextInitializer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := &RequestContext{}

		if sw, ok := w.(*StatWriter); ok {
			sw.RequestContext = reqCtx
		}

		ctx := WithRequestContext(r.Context(), reqCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
