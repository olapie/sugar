package cryptox

import (
	"code.olapie.com/sugar/slicing"
	"code.olapie.com/sugar/testx"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("123", []byte("abc"))
	t.Log(key)
	hash := key.Hash()
	t.Log(hash)
}

func TestAES(t *testing.T) {
	raw := testx.RandomBytes(8)
	password := testx.RandomString(8)
	stream1 := getCipherStream(password)
	stream2 := getCipherStream(password)

	data := slicing.Clone(raw)
	testx.Equal(t, raw, data)

	// encrypt
	stream1.XORKeyStream(data, data)
	testx.NotEqual(t, raw, data)

	// decrypt
	stream2.XORKeyStream(data, data)
	testx.Equal(t, raw, data)
}
