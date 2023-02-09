package hashutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

func MD5[T ~string | ~[]byte](data T) string {
	sum := md5.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}

func SHA1[T ~string | ~[]byte](data T) string {
	sha1er := sha1.New()
	b := []byte(data)
	for len(b) > 0 {
		n, err := sha1er.Write(b)
		if err != nil {
			panic(err)
		}
		b = b[n:]
	}
	return hex.EncodeToString(sha1er.Sum(nil))
}

func SHA256[T ~string | ~[]byte](data T) string {
	sha256er := sha256.New()
	b := []byte(data)
	for len(b) > 0 {
		n, err := sha256er.Write(b)
		if err != nil {
			panic(err)
		}
		b = b[n:]
	}
	return hex.EncodeToString(sha256er.Sum(nil))
}

func Hash32[T ~string | ~[]byte](data T) [32]byte {
	// enforce allocation for v
	v := append([]byte{}, []byte(data)...)
	sum := sha256.Sum256(v)
	for i := 0; i < 3; i++ {
		v = append(sum[:], v...)
		sum = sha256.Sum256(v)
	}
	return sum
}

func FileMD5(filename string) (string, error) {
	return hashFile(filename, md5.New())
}

func FileSHA1(filename string) (string, error) {
	return hashFile(filename, sha1.New())
}

func FileSHA256(filename string) (string, error) {
	return hashFile(filename, sha256.New())
}

func hashFile(filename string, h hash.Hash) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
