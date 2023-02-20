package httpkit

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
)

type statusGetter interface {
	Status() int
}

// http.Flusher doesn't return error, however gzip.Writer/deflate.WrapResponseWriter only implement `Flush() error`
type flusher interface {
	Flush() error
}

var (
	_ statusGetter        = (*WrapResponseWriter)(nil)
	_ http.Hijacker       = (*WrapResponseWriter)(nil)
	_ http.Flusher        = (*WrapResponseWriter)(nil)
	_ http.ResponseWriter = (*WrapResponseWriter)(nil)
)

// WrapResponseWriter is a wrapper of http.ResponseWriter to make sure write status code only one time
type WrapResponseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func NewWrapResponseWriter(rw http.ResponseWriter) *WrapResponseWriter {
	return &WrapResponseWriter{
		ResponseWriter: rw,
	}
}

func (w *WrapResponseWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusContinue {
		log.Println("cannot write invalid status code", statusCode)
		statusCode = http.StatusInternalServerError
	}
	if w.status > 0 {
		log.Println("status code already written")
		return
	}
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *WrapResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body = data
	return w.ResponseWriter.Write(data)
}

func (w *WrapResponseWriter) Status() int {
	return w.status
}

func (w *WrapResponseWriter) Body() []byte {
	return w.body
}

func (w *WrapResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *WrapResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *WrapResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("hijack not supported")
}

func (w *WrapResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
