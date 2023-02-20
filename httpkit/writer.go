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

// http.Flusher doesn't return error, however gzip.Writer/deflate.Writer only implement `Flush() error`
type flusher interface {
	Flush() error
}

var (
	_ statusGetter  = (*Writer)(nil)
	_ http.Hijacker = (*Writer)(nil)
	_ http.Flusher  = (*Writer)(nil)
)

// Writer is a wrapper of http.ResponseWriter to make sure write status code only one time
type Writer struct {
	http.ResponseWriter
	status int
	body   []byte
}

func NewWriter(rw http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: rw,
	}
}

func (w *Writer) WriteHeader(statusCode int) {
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

func (w *Writer) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body = data
	return w.ResponseWriter.Write(data)
}

func (w *Writer) Status() int {
	return w.status
}

func (w *Writer) Body() []byte {
	return w.body
}

func (w *Writer) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *Writer) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *Writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("hijack not supported")
}

func (w *Writer) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
