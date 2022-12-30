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

func TestEncryptedReader(t *testing.T) {
	raw := []byte(xhash.SHA1(time.Now().String()))
	enc := &bytes.Buffer{}
	{
		r := olasec.NewEncryptedReader(bytes.NewReader(raw), "123")
		n, err := io.Copy(enc, r)
		xtest.NoError(t, err)
		t.Log(n)
	}
	{
		data, err := olasec.Encrypt(raw, "123")
		xtest.NoError(t, err)
		xtest.Equal(t, enc.Bytes(), data)
	}
}
