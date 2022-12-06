package hashing

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func MD5(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}

func SHA1(str string) string {
	sha1er := sha1.New()
	b := []byte(str)
	for len(b) > 0 {
		n, err := sha1er.Write(b)
		if err != nil {
			panic(err)
		}
		b = b[n:]
	}
	return hex.EncodeToString(sha1er.Sum(nil))
}

func SHA256(str string) string {
	sha256er := sha256.New()
	b := []byte(str)
	for len(b) > 0 {
		n, err := sha256er.Write(b)
		if err != nil {
			panic(err)
		}
		b = b[n:]
	}
	return hex.EncodeToString(sha256er.Sum(nil))
}

func Hash32(b []byte) [32]byte {
	v := b
	sum := sha256.Sum256(v)
	for i := 0; i < 3; i++ {
		v = append(sum[:], v...)
		sum = sha256.Sum256(v)
	}
	return sum
}
