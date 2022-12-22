package olasec_test

import (
	"bytes"
	"code.olapie.com/sugar/olasec"
	"code.olapie.com/sugar/types"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
)

func TestDecryptedReader(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc, err := olasec.Encrypt(raw, "123")
	testx.NoError(t, err)
	r := olasec.NewDecryptedReader(bytes.NewReader(enc), "123")
	dec := &bytes.Buffer{}
	n, err := io.Copy(dec, r)
	testx.NoError(t, err)
	t.Log(n)
	testx.Equal(t, raw, dec.Bytes())
}

func BenchmarkDecryptedReader(b *testing.B) {
	raw := testx.RandomBytes(int(4 * types.MB))
	enc, err := olasec.Encrypt(raw, "123")
	testx.NoError(b, err)
	for i := 0; i < b.N; i++ {
		olasec.DecryptBytes(enc, "123")
	}
}
