package middleware

import (
	"net/http"
	"slices"
)

func BuildChain(final http.Handler, middlewares ...Middleware) http.Handler {
	handler := final
	for _, mw := range slices.Backward(middlewares) {
		handler = mw(handler)
	}
	return handler
}
