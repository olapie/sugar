package olasec_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/v2/olasec"
	"code.olapie.com/sugar/v2/testutil"
	"code.olapie.com/sugar/v2/types"
)

func TestDecryptedReader(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	enc, err := olasec.Encrypt(raw, "123")
	testutil.NoError(t, err)
	r := olasec.NewDecryptedReader(bytes.NewReader(enc), "123")
	dec := &bytes.Buffer{}
	n, err := io.Copy(dec, r)
	testutil.NoError(t, err)
	t.Log(n)
	testutil.Equal(t, raw, dec.Bytes())
}

func BenchmarkDecryptedReader(b *testing.B) {
	raw := testutil.RandomBytes(int(4 * types.MB))
	enc, err := olasec.Encrypt(raw, "123")
	testutil.NoError(b, err)
	for i := 0; i < b.N; i++ {
		olasec.Decrypt(enc, "123")
	}
}
