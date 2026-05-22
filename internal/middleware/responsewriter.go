package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type StatWriter struct {
	http.ResponseWriter
	Status         int
	Bytes          int
	RequestContext *RequestContext
}

func (rw *StatWriter) WriteHeader(status int) {
	rw.Status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *StatWriter) Write(b []byte) (int, error) {
	if rw.Status == 0 {
		rw.Status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.Bytes += n
	return n, err
}

func (rw *StatWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (rw *StatWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijacker not supported")
	}
	return h.Hijack()
}
