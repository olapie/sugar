package cryptox_test

import (
	"bytes"
	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/cryptox"
)

func TestEncryptedWriter(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc := &bytes.Buffer{}
	w := cryptox.NewEncryptedWriter(enc, "123")
	n, err := io.Copy(w, bytes.NewReader(raw))
	w.Close()
	testx.NoError(t, err)
	t.Log(n)
	data, err := cryptox.Encrypt(raw, "123")
	testx.NoError(t, err)
	testx.Equal(t, enc.Bytes(), data)
}
