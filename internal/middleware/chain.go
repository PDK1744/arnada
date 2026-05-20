package middleware

import (
	"net/http"
	"slices"
)

func BuildChain(final http.Handler, middlewares ...Middleware) http.Handler {
	handler := final
	for _, midware := range slices.Backward(middlewares) {
		handler = midware(handler)
	}
	return handler
}
