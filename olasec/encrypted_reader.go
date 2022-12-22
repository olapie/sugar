package olasec

import (
	"io"
)

var _ io.Reader = (*EncryptedReader)(nil)

// EncryptedReader reads and encrypts data from original reader
type EncryptedReader struct {
	r      io.Reader
	stream *cipherStream
	header []byte
}

func NewEncryptedReader(r io.Reader, password string) *EncryptedReader {
	reader := &EncryptedReader{
		r:      r,
		stream: getCipherStream(password),
	}

	reader.header = []byte(MagicNumberV1)
	reader.header = append(reader.header, reader.stream.keyHash[:]...)
	return reader
}

func (r *EncryptedReader) Read(p []byte) (int, error) {
	offset := 0
	data := p
	if len(r.header) > 0 {
		n := copy(p, r.header)
		r.header = r.header[n:]
		if n == len(p) {
			return n, nil
		}
		data = p[n:]
		offset = n
	}

	n, err := r.r.Read(data)
	if n > 0 {
		r.stream.XORKeyStream(data[:n], data[:n])
	}
	return offset + n, err
}
