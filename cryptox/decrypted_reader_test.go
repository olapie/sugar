package cryptox_test

import (
	"bytes"
	"code.olapie.com/sugar/cryptox"
	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
	"io"
	"testing"
	"time"
)

func TestDecryptedReader(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc, err := cryptox.Encrypt(raw, "123")
	testx.NoError(t, err)
	r := cryptox.NewDecryptedReader(bytes.NewReader(enc), "123")
	dec := &bytes.Buffer{}
	n, err := io.Copy(dec, r)
	testx.NoError(t, err)
	t.Log(n)
	testx.Equal(t, raw, dec.Bytes())
}
