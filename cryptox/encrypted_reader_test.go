package cryptox_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/cryptox"
	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
)

func TestEncryptedReader(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc := &bytes.Buffer{}
	r := cryptox.NewEncryptedReader(bytes.NewReader(raw), "123")
	n, err := io.Copy(enc, r)
	testx.NoError(t, err)
	t.Log(n)
	data, err := cryptox.Encrypt(raw, "123")
	testx.NoError(t, err)
	testx.Equal(t, enc.Bytes(), data)
}
