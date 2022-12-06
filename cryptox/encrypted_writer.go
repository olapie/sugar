package cryptox

import (
	"bytes"
	"fmt"
	"io"
)

var _ io.WriteCloser = (*EncryptedWriter)(nil)

// EncryptedWriter encrypts and write data into original writer
type EncryptedWriter struct {
	w     io.Writer
	key   Key
	block [encryptionBlockSize]byte
	buf   bytes.Buffer

	headerWritten bool
}

func NewEncryptedWriter[K string | Key](w io.Writer, k K) *EncryptedWriter {
	Writer := &EncryptedWriter{
		w:   w,
		key: getKey(k),
	}

	return Writer
}

func (w *EncryptedWriter) Write(p []byte) (n int, err error) {
	if !w.headerWritten {
		nWrite, err := w.w.Write([]byte(MagicNumber))
		if err != nil {
			return nWrite, err
		}
		hash := w.key.Hash()
		nWrite, err = w.w.Write(hash[:])
		if err != nil {
			return nWrite + MagicNumberSize, err
		}
		w.headerWritten = true
	}

	n, err = w.buf.Write(p)
	if err != nil {
		return 0, err
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

func (w *EncryptedWriter) Close() error {
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
