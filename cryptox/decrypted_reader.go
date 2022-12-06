package cryptox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var _ io.Reader = (*DecryptedReader)(nil)

// DecryptedReader reads and decrypt data from original reader
type DecryptedReader struct {
	r      io.Reader
	key    Key
	block  [encryptionBlockSize]byte
	srcBuf bytes.Buffer
	dstBuf bytes.Buffer
	eof    bool

	readHeader bool
}

func NewDecryptedReader[K string | Key](r io.Reader, k K) *DecryptedReader {
	reader := &DecryptedReader{
		r:   r,
		key: getKey(k),
	}

	return reader
}

func (r *DecryptedReader) Read(p []byte) (n int, err error) {
	if !r.readHeader {
		r.readHeader = true
		var header [HeaderSize]byte
		_, err = r.r.Read(header[:])
		if err != nil {
			return 0, err
		}

		if !ValidateKey(header[:], r.key) {
			return 0, ErrKey
		}
	}
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

func (r *DecryptedReader) readBlock() error {
	n, readErr := r.r.Read(r.block[:])
	if _, err := r.srcBuf.Write(r.block[:n]); err != nil {
		return err
	}

	if readErr == io.EOF {
		for r.srcBuf.Len() > 0 {
			next := r.srcBuf.Next(encryptionBlockSize)
			if err := r.key.AES(next, next); err != nil {
				return fmt.Errorf("aes: %w", err)
			}

			if _, err := r.dstBuf.Write(next); err != nil {
				return fmt.Errorf("cannot write: %w", err)
			}
		}
	} else {
		for r.srcBuf.Len() >= encryptionBlockSize {
			next := r.srcBuf.Next(encryptionBlockSize)
			if err := r.key.AES(next, next); err != nil {
				return fmt.Errorf("aes: %w", err)
			}

			if _, err := r.dstBuf.Write(next); err != nil {
				return fmt.Errorf("cannot write: %w", err)
			}
		}
	}
	return readErr
}
