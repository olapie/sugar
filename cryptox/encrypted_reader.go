package cryptox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var _ io.Reader = (*EncryptedReader)(nil)

// EncryptedReader reads and encrypts data from original reader
type EncryptedReader struct {
	r      io.Reader
	stream *cipherStream
	block  [encryptionBlockSize]byte
	srcBuf bytes.Buffer
	dstBuf bytes.Buffer
	eof    bool
}

func NewEncryptedReader(r io.Reader, password string) *EncryptedReader {
	reader := &EncryptedReader{
		r:      r,
		stream: getCipherStream(password),
	}

	reader.dstBuf.Write([]byte(MagicNumber))
	reader.dstBuf.Write(reader.stream.keyHash[:])
	return reader
}

func (r *EncryptedReader) Read(p []byte) (n int, err error) {
	size := len(p)
	for n < size {
		nRead, err := r.dstBuf.Read(p[n:])
		n += nRead
		if n >= size {
			return n, err
		}

		if r.eof {
			return n, io.EOF
		}

		err = r.readBlock()
		if err != nil {
			r.eof = errors.Is(err, io.EOF)
			if !r.eof {
				return n, err
			}
		}
	}

	return n, nil
}

func (r *EncryptedReader) readBlock() error {
	n, err := r.r.Read(r.block[:])
	if _, wErr := r.srcBuf.Write(r.block[:n]); wErr != nil {
		err = wErr
	}

	if r.srcBuf.Len() == 0 {
		return err
	}

	if encErr := r.encrypt(err == io.EOF); encErr != nil {
		err = encErr
	}

	return err
}

func (r *EncryptedReader) encrypt(all bool) error {
	for r.srcBuf.Len() >= encryptionBlockSize {
		next := r.srcBuf.Next(encryptionBlockSize)
		r.stream.XORKeyStream(next, next)

		if _, err := r.dstBuf.Write(next); err != nil {
			return fmt.Errorf("cannot write: %w", err)
		}
	}

	if all {
		for r.srcBuf.Len() > 0 {
			next := r.srcBuf.Next(encryptionBlockSize)
			r.stream.XORKeyStream(next, next)
			if _, err := r.dstBuf.Write(next); err != nil {
				return fmt.Errorf("cannot write: %w", err)
			}
		}
	}
	return nil
}
