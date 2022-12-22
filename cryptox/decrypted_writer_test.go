package cryptox_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"

	"code.olapie.com/sugar/cryptox"
)

func TestDecryptedWriter(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc, err := cryptox.Encrypt(raw, "123")
	testx.NoError(t, err)
	dec := &bytes.Buffer{}
	w := cryptox.NewDecryptedWriter(dec, "123")
	n, err := io.Copy(w, bytes.NewReader(enc))
	t.Log(n)
	testx.NoError(t, err)
	testx.Equal(t, raw, dec.Bytes())
}
