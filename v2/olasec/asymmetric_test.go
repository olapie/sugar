package olasec_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
	"time"

	"code.olapie.com/sugar/olasec"
	"code.olapie.com/sugar/xhash"
	"code.olapie.com/sugar/xtest"
)

func TestEncodePrivateKey(t *testing.T) {
	pk, err := olasec.GeneratePrivateKey()
	xtest.NoError(t, err)
	data, err := olasec.EncodePrivateKey(pk, "hello")
	xtest.NoError(t, err)
	_, err = olasec.DecodePrivateKey(data, "hi")
	xtest.Error(t, err)
	pk2, err := olasec.DecodePrivateKey(data, "hello")
	xtest.NoError(t, err)
	digest := []byte(xhash.SHA1(time.Now().String()))
	sign1, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	xtest.NoError(t, err)
	sign2, err := ecdsa.SignASN1(rand.Reader, pk2, digest[:])
	xtest.NoError(t, err)

	xtest.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign1))
	xtest.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign2))
	xtest.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign1))
	xtest.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign2))

	xtest.Equal(t, pk.X, pk2.X)
	xtest.Equal(t, pk.Y, pk2.Y)
	xtest.Equal(t, pk.D, pk2.D)
}

func TestEncodePublicKey(t *testing.T) {
	pk, err := olasec.GeneratePrivateKey()
	xtest.NoError(t, err)
	data, err := olasec.EncodePublicKey(&pk.PublicKey)
	xtest.NoError(t, err)
	pub, err := olasec.DecodePublicKey(data)
	xtest.NoError(t, err)
	digest := []byte(xhash.SHA1(time.Now().String()))
	sign, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	xtest.NoError(t, err)
	xtest.True(t, ecdsa.VerifyASN1(pub, digest[:], sign))
	xtest.True(t, pk.PublicKey.Equal(pub))
}
