package cryptox

import (
	"bytes"
	"io"

	"code.olapie.com/sugar/errorx"
)

type Encrypter interface {
	Encrypt(raw []byte) ([]byte, error)
}

type Decrypter interface {
	Decrypt(raw []byte) ([]byte, error)
}

func MakeEncrypter[K string | Key](k K) Encrypter {
	return &encrypterImpl{
		key: getKey(k),
	}
}

func MakeDecrypter[K string | Key](k K) Decrypter {
	return &decrypterImpl{
		key: getKey(k),
	}
}

func Encrypt[K string | Key](raw []byte, k K) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := NewEncryptedWriter(buf, k)
	_, err := io.Copy(w, bytes.NewReader(raw))
	err = errorx.Or(w.Close(), err)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decrypt[K string | Key](data []byte, k K) ([]byte, error) {
	r := NewDecryptedReader(bytes.NewReader(data), k)
	w := bytes.NewBuffer(nil)
	_, err := io.Copy(w, r)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func ChangeKey[K string | Key](data []byte, oldKey, newKey K) ([]byte, error) {
	raw, err := Decrypt(data, oldKey)
	if err != nil {
		return nil, err
	}

	return Encrypt(raw, newKey)
}

type encrypterImpl struct {
	key Key
}

func (e *encrypterImpl) Encrypt(raw []byte) ([]byte, error) {
	return Encrypt(raw, e.key)
}

type decrypterImpl struct {
	key Key
}

func (e *decrypterImpl) Decrypt(raw []byte) ([]byte, error) {
	return Decrypt(raw, e.key)
}
