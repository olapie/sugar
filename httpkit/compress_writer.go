package httpkit

import (
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	_ statusGetter  = (*CompressedWriter)(nil)
	_ http.Hijacker = (*CompressedWriter)(nil)
	_ http.Flusher  = (*CompressedWriter)(nil)
)

type CompressedWriter struct {
	*WrapResponseWriter
	compressWriter io.Writer
	err            error
	hasBody        bool
}

func NewCompressedWriter(w *WrapResponseWriter, encoding string) (*CompressedWriter, error) {
	switch encoding {
	case "gzip":
		cw := &CompressedWriter{}
		cw.WrapResponseWriter = w
		cw.compressWriter = gzip.NewWriter(w)
		SetContentEncoding(w.Header(), encoding)
		return cw, nil
	case "deflate":
		fw, err := flate.NewWriter(w, flate.DefaultCompression)
		if err != nil {
			return nil, fmt.Errorf("new flate writer: %w", err)
		}
		cw := &CompressedWriter{}
		cw.compressWriter = fw
		cw.WrapResponseWriter = w
		SetContentEncoding(w.Header(), encoding)
		return cw, nil
	default:
		return nil, errors.New("unsupported encoding")
	}
}

func (w *CompressedWriter) Write(data []byte) (int, error) {
	if !w.hasBody {
		w.hasBody = len(data) > 0
	}
	return w.compressWriter.Write(data)
}

func (w *CompressedWriter) Flush() {
	// Flush the compressed writer, then flush httpWriter
	if f, ok := w.compressWriter.(flusher); ok {
		if err := f.Flush(); err != nil {
			log.Println("flush", err)
			w.err = err
		}
		w.WrapResponseWriter.Flush()
	}
}

func (w *CompressedWriter) Error() error {
	return w.err
}

func (w *CompressedWriter) Close() error {
	if !w.hasBody {
		w.WrapResponseWriter.Flush()
		return nil
	}
	if closer, ok := w.compressWriter.(io.Closer); ok {
		// Closing a writer without written data will cause an error if Response status is 204 NoContent
		return closer.Close()
	}
	return nil
}

func CompressWriter(w http.ResponseWriter, encodings ...string) (http.ResponseWriter, error) {
	if _, ok := w.(*CompressedWriter); ok {
		log.Println("cannot compress a compressed writer")
		return w, nil
	}

	rw, _ := w.(*WrapResponseWriter)
	if rw == nil {
		rw = NewWrapResponseWriter(w)
	}

	if len(encodings) == 0 {
		return nil, errors.New("missing encodings")
	}

	var err error
	for _, encoding := range encodings {
		cw, er := NewCompressedWriter(rw, encoding)
		if er == nil {
			return cw, nil
		}
		err = errors.Join(err, er)
	}

	return nil, err
}
