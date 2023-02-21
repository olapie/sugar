package httpwriter

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

// http.Flusher doesn't return error, however gzip.Writer/deflate.Wrapper only implement `Flush() error`
type flusher interface {
	Flush() error
}

var (
	_ statusGetter        = (*Wrapper)(nil)
	_ http.Hijacker       = (*Wrapper)(nil)
	_ http.Flusher        = (*Wrapper)(nil)
	_ http.ResponseWriter = (*Wrapper)(nil)
)

// Wrapper is a wrapper of http.ResponseWriter to make sure write status code only one time
type Wrapper struct {
	http.ResponseWriter
	status int
	body   []byte
}

func NewWrapper(rw http.ResponseWriter) *Wrapper {
	return &Wrapper{
		ResponseWriter: rw,
	}
}

func (w *Wrapper) WriteHeader(statusCode int) {
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

func (w *Wrapper) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body = data
	return w.ResponseWriter.Write(data)
}

func (w *Wrapper) Status() int {
	return w.status
}

func (w *Wrapper) Body() []byte {
	return w.body
}

func (w *Wrapper) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *Wrapper) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *Wrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("hijack not supported")
}

func (w *Wrapper) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
