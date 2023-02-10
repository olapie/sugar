package olasec_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/hashutil"
	"code.olapie.com/sugar/v2/olasec"
	"code.olapie.com/sugar/v2/testutil"
)

func TestEncodePrivateKey(t *testing.T) {
	pk, err := olasec.GeneratePrivateKey()
	testutil.NoError(t, err)
	data, err := olasec.EncodePrivateKey(pk, "hello")
	testutil.NoError(t, err)
	_, err = olasec.DecodePrivateKey(data, "hi")
	testutil.Error(t, err)
	pk2, err := olasec.DecodePrivateKey(data, "hello")
	testutil.NoError(t, err)
	digest := []byte(hashutil.SHA1(time.Now().String()))
	sign1, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	testutil.NoError(t, err)
	sign2, err := ecdsa.SignASN1(rand.Reader, pk2, digest[:])
	testutil.NoError(t, err)

	testutil.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign1))
	testutil.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign2))
	testutil.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign1))
	testutil.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign2))

	testutil.Equal(t, pk.X, pk2.X)
	testutil.Equal(t, pk.Y, pk2.Y)
	testutil.Equal(t, pk.D, pk2.D)
}

func TestEncodePublicKey(t *testing.T) {
	pk, err := olasec.GeneratePrivateKey()
	testutil.NoError(t, err)
	data, err := olasec.EncodePublicKey(&pk.PublicKey)
	testutil.NoError(t, err)
	pub, err := olasec.DecodePublicKey(data)
	testutil.NoError(t, err)
	digest := []byte(hashutil.SHA1(time.Now().String()))
	sign, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	testutil.NoError(t, err)
	testutil.True(t, ecdsa.VerifyASN1(pub, digest[:], sign))
	testutil.True(t, pk.PublicKey.Equal(pub))
}
