package olasec_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
	"time"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/olasec"
	"code.olapie.com/sugar/testx"
)

func TestEncodePrivateKey(t *testing.T) {
	pk, err := olasec.GeneratePrivateKey()
	testx.NoError(t, err)
	data, err := olasec.EncodePrivateKey(pk, "hello")
	testx.NoError(t, err)
	_, err = olasec.DecodePrivateKey(data, "hi")
	testx.Error(t, err)
	pk2, err := olasec.DecodePrivateKey(data, "hello")
	testx.NoError(t, err)
	digest := []byte(hashing.SHA1(time.Now().String()))
	sign1, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	testx.NoError(t, err)
	sign2, err := ecdsa.SignASN1(rand.Reader, pk2, digest[:])
	testx.NoError(t, err)

	testx.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign1))
	testx.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign2))
	testx.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign1))
	testx.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign2))

	testx.Equal(t, pk.X, pk2.X)
	testx.Equal(t, pk.Y, pk2.Y)
	testx.Equal(t, pk.D, pk2.D)
}

func TestEncodePublicKey(t *testing.T) {
	pk, err := olasec.GeneratePrivateKey()
	testx.NoError(t, err)
	data, err := olasec.EncodePublicKey(&pk.PublicKey)
	testx.NoError(t, err)
	pub, err := olasec.DecodePublicKey(data)
	testx.NoError(t, err)
	digest := []byte(hashing.SHA1(time.Now().String()))
	sign, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	testx.NoError(t, err)
	testx.True(t, ecdsa.VerifyASN1(pub, digest[:], sign))
	testx.True(t, pk.PublicKey.Equal(pub))
}
