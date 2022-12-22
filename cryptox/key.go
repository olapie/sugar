package cryptox

import (
	"bytes"
	"code.olapie.com/sugar/errorx"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	KeySize     = 32
	KeyHashSize = 16
)

const ErrKey errorx.String = "invalid key"

type Key [KeySize]byte

//
//func (k *Key) AES(dst, src []byte) error {
//	block, err := aes.NewCipher((*k)[:KeySize/2])
//	if err != nil {
//		return err
//	}
//	stream := cipher.NewCTR(block, (*k)[KeySize/2:])
//	stream.XORKeyStream(dst, src)
//	return nil
//}

func (k *Key) Hash() [KeyHashSize]byte {
	md5Sum := md5.Sum((*k)[:])
	sha1Sum := sha1.Sum((*k)[:])
	hash := argon2.IDKey(sha1Sum[:], md5Sum[:], 1, 64, 4, KeyHashSize)
	var res KeyHash
	copy(res[:], hash)
	return res
}

// KeyHash is used to validate key before conducting decryption
type KeyHash [KeyHashSize]byte

func DeriveKey(password string, salt []byte) Key {

	if len(salt) == 0 {
		md5Sum := md5.Sum([]byte(strings.Repeat("ola"+password, 3)))
		salt = md5Sum[:]

	}
	k := argon2.IDKey([]byte(password), salt, 1, 64, 1, KeySize)

	if len(k) != 32 {
		panic(fmt.Errorf("key length is %d instead of %d", len(k), KeySize))
	}
	var res Key
	copy(res[:], k)
	return res
}

func ValidateKey[S string | []byte](s S, password string) bool {
	switch v := any(s).(type) {
	case string:
		sf, err := os.Open(v)
		if err != nil {
			return false
		}
		defer sf.Close()

		var header [HeaderSize]byte
		_, err = io.ReadFull(sf, header[:])

		if err != nil {
			return false
		}
		return getCipherStream(password).ValidatePassword(header[:])
	case []byte:
		return getCipherStream(password).ValidatePassword(v)
	default:
		return false
	}
}

func getCipherStream(password string) *cipherStream {
	key := DeriveKey(password, nil)
	block, err := aes.NewCipher((key)[:KeySize/2])
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, (key)[KeySize/2:])
	return &cipherStream{
		keyHash: key.Hash(),
		Stream:  stream,
	}
}

type cipherStream struct {
	keyHash KeyHash
	cipher.Stream
}

func (i *cipherStream) ValidatePassword(data []byte) bool {
	var header [HeaderSize]byte
	if len(data) < HeaderSize {
		return false
	}

	copy(header[:], data[:HeaderSize])

	if string(header[:MagicNumberSize]) != MagicNumber {
		return false
	}
	return bytes.Compare(header[MagicNumberSize:], i.keyHash[:]) == 0
}
