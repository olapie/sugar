package olasec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"golang.org/x/crypto/argon2"
)

// keyHashType is used to validate key before conducting decryption
type keyHashType [KeyHashSize]byte

var defaultSaltPart = []byte("ola-sec")

func deriveKey(password, salt []byte, keyLen int) []byte {
	if len(password) == 0 {
		panic("password is empty")
	}

	if len(salt) == 0 {
		salt = append(defaultSaltPart, password...)
		salt = bytes.Repeat(salt, 3)
		sum := md5.Sum(salt)
		salt = sum[:]
	}

	if keyLen <= 0 || keyLen > 128 {
		panic(fmt.Sprintf("invalid keyLen: %d", keyLen))
	}
	k := argon2.IDKey(password, salt, 1, 128, 1, uint32(keyLen))
	if len(k) != keyLen {
		panic(fmt.Errorf("key length is %d instead of %d", len(k), keyLen))
	}
	return k
}

func hashKey(k Key) [KeyHashSize]byte {
	md5Sum := md5.Sum(k[:])
	sha1Sum := sha1.Sum(k[:])
	b := deriveKey(sha1Sum[:], md5Sum[:], KeyHashSize)
	var hash keyHashType
	copy(hash[:], b)
	return hash
}

func getCipherStream(password string) *cipherStream {
	key := DeriveKey(password, nil)
	block, err := aes.NewCipher((key)[:KeySize/2])
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, (key)[KeySize/2:])
	return &cipherStream{
		keyHash: hashKey(key),
		Stream:  stream,
	}
}

type cipherStream struct {
	keyHash keyHashType
	cipher.Stream
}

func (i *cipherStream) ValidatePassword(data []byte) bool {
	var header [HeaderSize]byte
	if len(data) < HeaderSize {
		return false
	}

	copy(header[:], data[:HeaderSize])

	if string(header[:MagicNumberSize]) != MagicNumberV1 {
		return false
	}
	return bytes.Compare(header[MagicNumberSize:], i.keyHash[:]) == 0
}
