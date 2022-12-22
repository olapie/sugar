package cryptox

import (
	"bytes"
	"fmt"
	"io"
)

var _ io.WriteCloser = (*DecryptedWriter)(nil)

// DecryptedWriter decrypts and writes data into original Writer
type DecryptedWriter struct {
	w      io.Writer
	stream *cipherStream
	block  [encryptionBlockSize]byte
	buf    bytes.Buffer

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

	if !w.headerDecrypted && w.buf.Len() >= HeaderSize {
		w.headerDecrypted = true
		header := w.buf.Next(HeaderSize)
		if !w.stream.ValidatePassword(header) {
			return n, ErrKey
		}
	}

	for w.buf.Len() >= encryptionBlockSize {
		next := w.buf.Next(encryptionBlockSize)
		w.stream.XORKeyStream(next, next)

		if _, err := w.w.Write(next); err != nil {
			return n, fmt.Errorf("cannot write: %w", err)
		}
	}

	return n, nil
}

func (w *DecryptedWriter) Close() error {
	for w.buf.Len() > 0 {
		next := w.buf.Next(encryptionBlockSize)
		w.stream.XORKeyStream(next, next)
		if _, err := w.w.Write(next); err != nil {
			return fmt.Errorf("cannot write: %w", err)
		}
	}
	return nil
}
