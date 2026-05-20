package middleware

import (
	"context"
	"time"
)

type ctxKey int

const contextKey ctxKey = iota

type RequestContext struct {
	RequestID string
	Method    string
	Path      string
	Host      string

	Upstream string

	StartTime  time.Time
	DurationMs int64

	Status int

	ReqBodySize  int
	RespBodySize int

	Error string
}

// attach *RequestContext to context
func WithRequestContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, contextKey, rc)
}

func GetRequestContext(ctx context.Context) (*RequestContext, bool) {
	rc, ok := ctx.Value(contextKey).(*RequestContext)
	return rc, ok
}
