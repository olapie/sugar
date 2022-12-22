package cryptox

import (
	"io"
)

var _ io.Reader = (*DecryptedReader)(nil)

// DecryptedReader reads and decrypt data from original reader
type DecryptedReader struct {
	r          io.Reader
	stream     *cipherStream
	readHeader bool
}

func NewDecryptedReader(r io.Reader, password string) *DecryptedReader {
	reader := &DecryptedReader{
		r:      r,
		stream: getCipherStream(password),
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

		if !r.stream.ValidatePassword(header[:]) {
			return 0, ErrKey
		}
	}

	n, err = r.r.Read(p)
	if err != nil {
		return n, err
	}
	r.stream.XORKeyStream(p, p)
	return n, nil
}
