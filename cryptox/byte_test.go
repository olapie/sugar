package cryptox_test

import (
	"crypto/rand"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/cryptox"
	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
)

func TestEncryptBytes(t *testing.T) {
	password := hashing.SHA1(time.Now().String())
	testEncryptBytes(t, 1<<4+9, password)
	testEncryptBytes(t, 1<<24, password)
}

func testEncryptBytes(t *testing.T, size int, password string) {
	raw := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, raw[:])
	if err != nil {
		t.Fatal(err)
	}

	enc, err := cryptox.Encrypt(raw[:], password)
	t.Log(enc[:30])
	testx.NoError(t, err)
	testx.True(t, cryptox.IsEncrypted(enc))
	dec, err := cryptox.Decrypt(enc, password)
	testx.NoError(t, err)
	testx.False(t, cryptox.IsEncrypted(dec), dec[:cryptox.HeaderSize])
	testx.Equal(t, raw, dec)
}
