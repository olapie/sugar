package olasec_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/hashutil"
	"code.olapie.com/sugar/v2/olasec"
	"code.olapie.com/sugar/v2/testutil"
)

func TestDecryptedWriter(t *testing.T) {
	raw := []byte(hashutil.SHA1(time.Now().String()))
	enc, err := olasec.Encrypt(raw, "123")
	testutil.NoError(t, err)
	dec := &bytes.Buffer{}
	w := olasec.NewDecryptedWriter(dec, "123")
	n, err := io.Copy(w, bytes.NewReader(enc))
	t.Log(n)
	testutil.NoError(t, err)
	testutil.Equal(t, raw, dec.Bytes())
}
