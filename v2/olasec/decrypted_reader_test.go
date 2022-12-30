package olasec_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/olasec"
	"code.olapie.com/sugar/xtype"

	"code.olapie.com/sugar/xhash"
	"code.olapie.com/sugar/xtest"
)

func TestDecryptedReader(t *testing.T) {
	raw := []byte(xhash.SHA1(time.Now().String()))
	enc, err := olasec.Encrypt(raw, "123")
	xtest.NoError(t, err)
	r := olasec.NewDecryptedReader(bytes.NewReader(enc), "123")
	dec := &bytes.Buffer{}
	n, err := io.Copy(dec, r)
	xtest.NoError(t, err)
	t.Log(n)
	xtest.Equal(t, raw, dec.Bytes())
}

func BenchmarkDecryptedReader(b *testing.B) {
	raw := xtest.RandomBytes(int(4 * xtype.MB))
	enc, err := olasec.Encrypt(raw, "123")
	xtest.NoError(b, err)
	for i := 0; i < b.N; i++ {
		olasec.Decrypt(enc, "123")
	}
}
