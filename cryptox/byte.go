package cryptox

import (
	"bytes"
	"crypto/cipher"
	"io"

	"code.olapie.com/sugar/errorx"
)

type Encrypter interface {
	Encrypt(raw []byte) []byte
}

type Decrypter interface {
	Decrypt(raw []byte) []byte
}

func MakeEncrypter(password string) Encrypter {
	return &encrypterImpl{
		stream: getCipherStream(password),
	}
}

func MakeDecrypter(password string) Decrypter {
	return &decrypterImpl{
		stream: getCipherStream(password),
	}
}

func Encrypt(raw []byte, password string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := NewEncryptedWriter(buf, password)
	_, err := io.Copy(w, bytes.NewReader(raw))
	err = errorx.Or(w.Close(), err)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decrypt(data []byte, password string) ([]byte, error) {
	r := NewDecryptedReader(bytes.NewReader(data), password)
	w := bytes.NewBuffer(nil)
	_, err := io.Copy(w, r)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func ChangeKey(data []byte, oldKey, newKey string) ([]byte, error) {
	raw, err := Decrypt(data, oldKey)
	if err != nil {
		return nil, err
	}

	return Encrypt(raw, newKey)
}

type encrypterImpl struct {
	stream cipher.Stream
}

func (e *encrypterImpl) Encrypt(raw []byte) []byte {
	dst := make([]byte, len(raw))
	e.stream.XORKeyStream(dst, raw)
	return dst
}

type decrypterImpl struct {
	stream cipher.Stream
}

func (e *decrypterImpl) Decrypt(raw []byte) []byte {
	dst := make([]byte, len(raw))
	e.stream.XORKeyStream(dst, raw)
	return dst
}
