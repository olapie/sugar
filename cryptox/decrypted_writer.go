package cryptox

import (
	"bytes"
	"fmt"
	"io"
)

var _ io.WriteCloser = (*DecryptedWriter)(nil)

// DecryptedWriter decrypts and writes data into original Writer
type DecryptedWriter struct {
	w     io.Writer
	key   Key
	block [encryptionBlockSize]byte
	buf   bytes.Buffer

	headerDecrypted bool
}

func NewDecryptedWriter[K string | Key](w io.Writer, k K) *DecryptedWriter {
	Writer := &DecryptedWriter{
		w:   w,
		key: getKey(k),
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
		if !ValidateKey(header, w.key) {
			return n, ErrKey
		}
	}

	for w.buf.Len() >= encryptionBlockSize {
		next := w.buf.Next(encryptionBlockSize)
		if err := w.key.AES(next, next); err != nil {
			return n, fmt.Errorf("aes: %w", err)
		}

		if _, err := w.w.Write(next); err != nil {
			return n, fmt.Errorf("cannot write: %w", err)
		}
	}

	return n, nil
}

func (w *DecryptedWriter) Close() error {
	for w.buf.Len() > 0 {
		next := w.buf.Next(encryptionBlockSize)
		if err := w.key.AES(next, next); err != nil {
			return fmt.Errorf("aes: %w", err)
		}

		if _, err := w.w.Write(next); err != nil {
			return fmt.Errorf("cannot write: %w", err)
		}
	}
	return nil
}
