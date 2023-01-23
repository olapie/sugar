package olasec_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/olasec"

	"code.olapie.com/sugar/v2/xhash"
	"code.olapie.com/sugar/v2/xtest"
)

func TestDecryptedWriter(t *testing.T) {
	raw := []byte(xhash.SHA1(time.Now().String()))
	enc, err := olasec.Encrypt(raw, "123")
	xtest.NoError(t, err)
	dec := &bytes.Buffer{}
	w := olasec.NewDecryptedWriter(dec, "123")
	n, err := io.Copy(w, bytes.NewReader(enc))
	t.Log(n)
	xtest.NoError(t, err)
	xtest.Equal(t, raw, dec.Bytes())
}
