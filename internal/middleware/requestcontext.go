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

func (rc *RequestContext) SetResponse(status, respBodySize int, duration int64) {
	rc.Status = status
	rc.RespBodySize = respBodySize
	rc.DurationMs = duration
}

func (rc *RequestContext) SetRequest(reqId, method, path, host string, startTime time.Time, reqBodySize int) {
	rc.RequestID = reqId
	rc.Method = method
	rc.Path = path
	rc.Host = host
	rc.StartTime = startTime
	rc.ReqBodySize = reqBodySize
}
