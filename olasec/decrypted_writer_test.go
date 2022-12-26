package olasec_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/olasec"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
)

func TestDecryptedWriter(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc, err := olasec.Encrypt(raw, "123")
	testx.NoError(t, err)
	dec := &bytes.Buffer{}
	w := olasec.NewDecryptedWriter(dec, "123")
	n, err := io.Copy(w, bytes.NewReader(enc))
	t.Log(n)
	testx.NoError(t, err)
	testx.Equal(t, raw, dec.Bytes())
}
