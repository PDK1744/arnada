package middleware

import "time"

type RequestContext struct {
	RequestID string
	Method    string
	Path      string
	Host      string

	Upstream string

	StartTime  time.Time
	DurationMs int64

	Status int
	Bytes  int64

	ReqBodySize  int
	RespBodySize int

	Error string
}
