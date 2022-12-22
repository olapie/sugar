package olasec_test

import (
	"bytes"
	"code.olapie.com/sugar/olasec"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
)

func TestEncryptedWriter(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc := &bytes.Buffer{}
	w := olasec.NewEncryptedWriter(enc, "123")
	n, err := io.Copy(w, bytes.NewReader(raw))
	testx.NoError(t, err)
	t.Log(n)
	data, err := olasec.Encrypt(raw, "123")
	testx.NoError(t, err)
	testx.Equal(t, enc.Bytes(), data)
}
