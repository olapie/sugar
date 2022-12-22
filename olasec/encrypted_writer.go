package olasec

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
	return &EncryptedWriter{
		w:      w,
		stream: getCipherStream(password),
	}
}

func (w *EncryptedWriter) Write(p []byte) (int, error) {
	if !w.headerWritten {
		_, err := w.w.Write([]byte(MagicNumberV1))
		if err != nil {
			return 0, err
		}
		_, err = w.w.Write(w.stream.keyHash[:])
		if err != nil {
			return 0, err
		}
		w.headerWritten = true
	}

	w.buf.Write(p)
	w.stream.XORKeyStream(w.buf.Bytes(), w.buf.Bytes())
	n, err := w.w.Write(w.buf.Bytes())
	w.buf.Reset()
	if err != nil {
		return n, fmt.Errorf("cannot write: %w", err)
	}
	return n, nil
}
