package cryptox

import (
	"code.olapie.com/sugar/errorx"
	"fmt"
	"io"
	"os"
)

type Source string
type Destination string

type FileEncrypter interface {
	Encrypt(dst Destination, src Source) error
}

type FileDecrypter interface {
	Decrypt(dst Destination, src Source) error
}

func MakeFileEncrypter[K string | Key](k K) FileEncrypter {
	return &fileEncrypterImpl{
		key: getKey(k),
	}
}

func MakeFileDecrypter[K string | Key](k K) FileDecrypter {
	return &fileDecrypterImpl{
		key: getKey(k),
	}
}

func EncryptFile[K string | Key](dst Destination, src Source, k K) error {
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
	w := NewEncryptedWriter(df, k)
	_, err = io.Copy(w, sf)
	return errorx.Or(w.Close(), err)
}

func DecryptFile[K string | Key](dst Destination, src Source, k K) error {
	sf, err := os.Open(string(src))
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", src, err)
	}
	defer sf.Close()
	r := NewDecryptedReader(sf, k)
	df, err := os.OpenFile(string(dst), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", dst, err)
	}
	defer df.Close()
	_, err = io.Copy(df, r)
	return err
}

func DecryptFileChunks[K string | Key](dst Destination, chunks []Source, k K) error {
	if len(chunks) == 0 {
		return nil
	}

	err := DecryptFile(dst, chunks[0], k)
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
		r := NewDecryptedReader(sf, k)
		_, err = io.Copy(df, r)
		sf.Close()
		if err != nil {
			return fmt.Errorf("io.Copy:%s, %w", chunk, err)
		}
	}

	return nil
}

func ChangeFileKey[K string | Key](dst Destination, src Source, oldKey, newKey K) error {
	if !ValidateKey(string(src), getKey(oldKey)) {
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
	key Key
}

func (e *fileEncrypterImpl) Encrypt(dst Destination, src Source) error {
	return EncryptFile(dst, src, e.key)
}

type fileDecrypterImpl struct {
	key Key
}

func (e *fileDecrypterImpl) Decrypt(dst Destination, src Source) error {
	return DecryptFile(dst, src, e.key)
}
