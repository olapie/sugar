package cryptox

import (
	"fmt"
	"io"
	"os"

	"code.olapie.com/sugar/errorx"
)

type Source string
type Destination string

type FileEncrypter interface {
	Encrypt(dst Destination, src Source) error
}

type FileDecrypter interface {
	Decrypt(dst Destination, src Source) error
}

func MakeFileEncrypter(password string) FileEncrypter {
	return &fileEncrypterImpl{
		password: password,
	}
}

func MakeFileDecrypter(password string) FileDecrypter {
	return &fileDecrypterImpl{
		password: password,
	}
}

func EncryptFile(dst Destination, src Source, password string) error {
	sf, err := os.Open(string(src))
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", src, err)
	}
	defer sf.Close()

	df, err := os.OpenFile(string(dst), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", dst, err)
	}
	defer df.Close()
	w := NewEncryptedWriter(df, password)
	_, err = io.Copy(w, sf)
	return errorx.Or(w.Close(), err)
}

func DecryptFile(dst Destination, src Source, password string) error {
	sf, err := os.Open(string(src))
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", src, err)
	}
	defer sf.Close()
	r := NewDecryptedReader(sf, password)
	df, err := os.OpenFile(string(dst), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", dst, err)
	}
	defer df.Close()
	_, err = io.Copy(df, r)
	return err
}

func DecryptFileChunks(dst Destination, chunks []Source, password string) error {
	if len(chunks) == 0 {
		return nil
	}

	err := DecryptFile(dst, chunks[0], password)
	if err != nil {
		return fmt.Errorf("cryptox.DecryptFile:%s, %w", chunks[0], err)
	}

	chunks = chunks[1:]
	if len(chunks) == 0 {
		return nil
	}

	df, err := os.Open(string(dst))
	if err != nil {
		return fmt.Errorf("os.Open:%s, %w", dst, err)
	}
	defer df.Close()

	for _, chunk := range chunks {
		sf, err := os.Open(string(chunk))
		if err != nil {
			return fmt.Errorf("os.Open:%s, %w", chunk, err)
		}
		r := NewDecryptedReader(sf, password)
		_, err = io.Copy(df, r)
		sf.Close()
		if err != nil {
			return fmt.Errorf("io.Copy:%s, %w", chunk, err)
		}
	}

	return nil
}

func ChangeFileKey(dst Destination, src Source, oldKey, newKey string) error {
	if !ValidateKey(string(src), oldKey) {
		return ErrKey
	}

	sf, err := os.Open(string(src))
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", src, err)
	}
	defer sf.Close()

	df, err := os.OpenFile(string(dst), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", dst, err)
	}
	defer df.Close()

	dr := NewDecryptedReader(sf, oldKey)
	ew := NewEncryptedWriter(df, newKey)
	_, err = io.Copy(ew, dr)
	return errorx.Or(ew.Close(), err)
}

type fileEncrypterImpl struct {
	password string
}

func (e *fileEncrypterImpl) Encrypt(dst Destination, src Source) error {
	return EncryptFile(dst, src, e.password)
}

type fileDecrypterImpl struct {
	password string
}

func (e *fileDecrypterImpl) Decrypt(dst Destination, src Source) error {
	return DecryptFile(dst, src, e.password)
}
