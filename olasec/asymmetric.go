package olasec

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
)

func GeneratePrivateKey[K *ecdsa.PrivateKey]() (K, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func MustGeneratePrivateKey[K *ecdsa.PrivateKey]() K {
	k, err := GeneratePrivateKey[K]()
	if err != nil {
		panic(err)
	}
	return k
}

func EncodePublicKey[K *ecdsa.PublicKey](k K) ([]byte, error) {
	pk := (*ecdsa.PublicKey)(k)
	return x509.MarshalPKIXPublicKey(pk)
}

func MustEncodePublicKey[K *ecdsa.PublicKey](k K) []byte {
	data, err := EncodePublicKey(k)
	if err != nil {
		panic(err)
	}
	return data
}

func EncodePrivateKey[K *ecdsa.PrivateKey](k K, passphrase string) ([]byte, error) {
	pk := (*ecdsa.PrivateKey)(k)
	data, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return nil, err
	}

	return EncryptBytes(data, passphrase)
}

func MustEncodePrivateKey[K *ecdsa.PrivateKey](k K, passphrase string) []byte {
	data, err := EncodePrivateKey(k, passphrase)
	if err != nil {
		panic(err)
	}
	return data
}

func DecodePrivateKey[K *ecdsa.PrivateKey](data []byte, passphrase string) (K, error) {
	size := len(data) - HeaderSize
	if size < 0 {
		return nil, fmt.Errorf("invalid data")
	}
	raw, err := DecryptBytes(data, passphrase)
	if err != nil {
		return nil, err
	}
	return x509.ParseECPrivateKey(raw)
}

func MustDecodePrivateKey[K *ecdsa.PrivateKey](data []byte, passphrase string) K {
	k, err := DecodePrivateKey(data, passphrase)
	if err != nil {
		panic(err)
	}
	return k
}

func DecodePublicKey[K *ecdsa.PublicKey](data []byte) (K, error) {
	anyPub, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	pub := anyPub.(*ecdsa.PublicKey)
	return pub, nil
}

func MustDecodePublicKey[K *ecdsa.PublicKey](data []byte) K {
	k, err := DecodePublicKey(data)
	if err != nil {
		panic(err)
	}
	return k
}
