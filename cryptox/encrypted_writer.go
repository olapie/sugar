package cryptox

import (
	"bytes"
	"fmt"
	"io"
)

var _ io.Writer = (*EncryptedWriter)(nil)

// EncryptedWriter encrypts and write data into original writer
type EncryptedWriter struct {
	w             io.Writer
	stream        *cipherStream
	buf           bytes.Buffer
	headerWritten bool
}

func NewEncryptedWriter(w io.Writer, password string) *EncryptedWriter {
	Writer := &EncryptedWriter{
		w:      w,
		stream: getCipherStream(password),
	}

	return Writer
}

func (w *EncryptedWriter) Write(p []byte) (n int, err error) {
	if !w.headerWritten {
		nWrite, err := w.w.Write([]byte(MagicNumber))
		if err != nil {
			return nWrite, err
		}
		nWrite, err = w.w.Write(w.stream.keyHash[:])
		if err != nil {
			return nWrite + MagicNumberSize, err
		}
		w.headerWritten = true
	}

	w.buf.Write(p)
	w.stream.XORKeyStream(w.buf.Bytes(), w.buf.Bytes())
	n, err = w.w.Write(w.buf.Bytes())
	w.buf.Reset()
	if err != nil {
		return n, fmt.Errorf("cannot write: %w", err)
	}
	return n, nil
}
