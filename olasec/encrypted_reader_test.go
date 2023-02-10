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

func TestEncryptedReader(t *testing.T) {
	raw := []byte(hashutil.SHA1(time.Now().String()))
	enc := &bytes.Buffer{}
	{
		r := olasec.NewEncryptedReader(bytes.NewReader(raw), "123")
		n, err := io.Copy(enc, r)
		testutil.NoError(t, err)
		t.Log(n)
	}
	{
		data, err := olasec.Encrypt(raw, "123")
		testutil.NoError(t, err)
		testutil.Equal(t, enc.Bytes(), data)
	}
}
