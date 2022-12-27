package olasec

import (
	"io"
)

var _ io.ReadSeeker = (*DecryptedReadSeeker)(nil)

// DecryptedReadSeeker reads and decrypt data from original reader
type DecryptedReadSeeker struct {
	r          io.ReadSeeker
	stream     *cipherStream
	readHeader bool
}

func NewDecryptedReadSeeker(r io.ReadSeeker, password string) *DecryptedReadSeeker {
	return &DecryptedReadSeeker{
		r:      r,
		stream: getCipherStream(password),
	}
}

func (r *DecryptedReadSeeker) Read(p []byte) (n int, err error) {
	if err := r.validateHeader(); err != nil {
		return 0, err
	}

	n, err = r.r.Read(p)
	if err != nil {
		return n, err
	}
	r.stream.XORKeyStream(p, p)
	return n, nil
}

func (r *DecryptedReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if err := r.validateHeader(); err != nil {
		return 0, err
	}

	pos, err := r.r.Seek(offset, whence)
	if err != nil {
		return 0, err
	}
	return pos + int64(HeaderSize), nil
}

func (r *DecryptedReadSeeker) validateHeader() error {
	if r.readHeader {
		return nil
	}
	r.readHeader = true
	var header [HeaderSize]byte
	_, err := r.r.Read(header[:])
	if err != nil {
		return err
	}
	if !r.stream.ValidatePassword(header[:]) {
		return ErrKey
	}
	return nil
}
