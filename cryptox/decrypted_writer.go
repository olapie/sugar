package cryptox

import (
	"bytes"
	"fmt"
	"io"
)

var _ io.Writer = (*DecryptedWriter)(nil)

// DecryptedWriter decrypts and writes data into original Writer
type DecryptedWriter struct {
	w               io.Writer
	stream          *cipherStream
	buf             bytes.Buffer
	headerDecrypted bool
}

func NewDecryptedWriter(w io.Writer, password string) *DecryptedWriter {
	Writer := &DecryptedWriter{
		w:      w,
		stream: getCipherStream(password),
	}
	return Writer
}

func (w *DecryptedWriter) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	if err != nil {
		return 0, err
	}

	if !w.headerDecrypted {
		if w.buf.Len() < HeaderSize {
			return len(p), nil
		}

		w.headerDecrypted = true
		header := w.buf.Next(HeaderSize)
		if !w.stream.ValidatePassword(header) {
			return len(p), ErrKey
		}
	}

	w.stream.XORKeyStream(w.buf.Bytes(), w.buf.Bytes())
	n, err = w.w.Write(w.buf.Bytes())
	w.buf.Reset()
	if err != nil {
		return n, fmt.Errorf("cannot write: %w", err)
	}
	return len(p), nil
}
