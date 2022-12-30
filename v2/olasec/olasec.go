package olasec

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"code.olapie.com/sugar/xerror"
)

// MagicNumberV1 is a defined 4-byte number to identify file type
// refer to https://en.wikipedia.org/wiki/List_of_file_signatures
// Header layout: magic number | key checksum
const (
	MagicNumberV1 = "\xFE\xF1\xFD\x01"

	MagicNumberSize = len(MagicNumberV1)
	KeySize         = 32
	KeyHashSize     = 16
	HeaderSize      = MagicNumberSize + KeyHashSize

	ErrKey xerror.String = "invalid key"
)

type Key [KeySize]byte
type SourceFile string
type DestFile string
type Filename string

func SF(filename string) SourceFile {
	return SourceFile(filename)
}

func DF(filename string) DestFile {
	return DestFile(filename)
}

func DeriveKey(password string, salt []byte) Key {
	k := deriveKey([]byte(password), salt, KeySize)
	var key Key
	copy(key[:], k)
	return key
}

func IsEncrypted(data []byte) bool {
	if len(data) < HeaderSize {
		return false
	}
	return string(data[:MagicNumberSize]) == MagicNumberV1
}

func IsEncryptedFile(filename string) bool {
	f, err := os.Open(string(filename))
	if err != nil {
		return false
	}
	defer f.Close()

	var header [HeaderSize]byte
	_, err = io.ReadFull(f, header[:])
	if err != nil {
		return false
	}
	return string(header[:MagicNumberSize]) == MagicNumberV1
}

func ValidatePassword(data []byte, password string) bool {
	if len(data) < HeaderSize {
		return false
	}
	return getCipherStream(password).ValidatePassword(data)
}

func ValidateFilePassword(filename, password string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	var header [HeaderSize]byte
	_, err = io.ReadFull(f, header[:])
	if err != nil {
		return false
	}
	return getCipherStream(password).ValidatePassword(header[:])
}

func Encrypt(raw []byte, password string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := NewEncryptedWriter(buf, password)
	_, err := io.Copy(w, bytes.NewReader(raw))
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

// DecryptInPlace will write decrypted data into parameter data
// input data parameter will be modified after decrypting
func DecryptInPlace(data []byte, password string) ([]byte, error) {
	if len(data) < HeaderSize {
		return nil, ErrKey
	}
	stream := getCipherStream(password)
	if !stream.ValidatePassword(data[:HeaderSize]) {
		return nil, ErrKey
	}
	stream.XORKeyStream(data[HeaderSize:], data[HeaderSize:])
	return data[HeaderSize:], nil
}

func ReEncrypt(data []byte, oldPassword, newPassword string) ([]byte, error) {
	raw, err := Decrypt(data, oldPassword)
	if err != nil {
		return nil, err
	}
	return Encrypt(raw, newPassword)
}

func EncryptFile(src SourceFile, dst DestFile, password string) error {
	sf, err := os.Open(string(src))
	if err != nil {
		return fmt.Errorf("os.Open: %s, %w", src, err)
	}
	defer sf.Close()

	df, err := os.OpenFile(string(dst), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s, %w", dst, err)
	}
	defer df.Close()
	w := NewEncryptedWriter(df, password)
	_, err = io.Copy(w, sf)
	return err
}

func DecryptFile(src SourceFile, dst DestFile, password string) error {
	sf, err := os.Open(string(src))
	if err != nil {
		return fmt.Errorf("os.Open: %s, %w", src, err)
	}
	defer sf.Close()
	r := NewDecryptedReader(sf, password)
	df, err := os.OpenFile(string(dst), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s, %w", dst, err)
	}
	defer df.Close()
	_, err = io.Copy(df, r)
	return err
}

func DecryptFileChunks(chunks []SourceFile, dst DestFile, password string) error {
	if len(chunks) == 0 {
		return nil
	}

	err := DecryptFile(chunks[0], dst, password)
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

func ReEncryptFile(src SourceFile, dst DestFile, srcPassword, dstPassword string) error {
	if !ValidateFilePassword(string(src), srcPassword) {
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

	dr := NewDecryptedReader(sf, srcPassword)
	ew := NewEncryptedWriter(df, dstPassword)
	_, err = io.Copy(ew, dr)
	return err
}
