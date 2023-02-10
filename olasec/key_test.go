package olasec

import (
	"bytes"
	"testing"

	"code.olapie.com/sugar/v2/slices"

	"code.olapie.com/sugar/v2/testutil"
)

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("123", []byte("abc"))
	t.Log(key)
	hash := hashKey(key)
	t.Log(hash)
}

func TestStream(t *testing.T) {
	raw := testutil.RandomBytes(8)
	password := testutil.RandomString(8)
	stream1 := getCipherStream(password)
	stream2 := getCipherStream(password)

	data := slices.Clone(raw)
	testutil.Equal(t, raw, data)

	// encrypt
	stream1.XORKeyStream(data, data)
	testutil.NotEqual(t, raw, data)

	// decrypt
	stream2.XORKeyStream(data, data)
	testutil.Equal(t, raw, data)
}

func TestStream2(t *testing.T) {
	n := 10 * 500
	raw := testutil.RandomBytes(n)
	password := testutil.RandomString(8)
	var encrypted []byte
	{
		stream1 := getCipherStream(password)
		stream2 := getCipherStream(password)

		var buf1 = bytes.NewBuffer(nil)
		step := 10
		for i := 0; i < n; i += step {
			var data = make([]byte, step)
			stream1.XORKeyStream(data[:], raw[i:i+step])
			buf1.Write(data[:])
		}

		var buf2 = bytes.NewBuffer(nil)
		step = 50
		for i := 0; i < n; i += step {
			var data = make([]byte, step)
			stream2.XORKeyStream(data[:], raw[i:i+step])
			buf2.Write(data[:])
		}

		testutil.Equal(t, buf1.Bytes(), buf2.Bytes())
		encrypted = buf1.Bytes()
	}

	{
		stream1 := getCipherStream(password)
		stream2 := getCipherStream(password)

		var buf1 = bytes.NewBuffer(nil)
		step := 20
		for i := 0; i < n; i += step {
			var data = make([]byte, step)
			stream1.XORKeyStream(data[:], encrypted[i:i+step])
			buf1.Write(data[:])
		}

		var buf2 = bytes.NewBuffer(nil)
		step = 100
		for i := 0; i < n; i += step {
			var data = make([]byte, step)
			stream2.XORKeyStream(data[:], encrypted[i:i+step])
			buf2.Write(data[:])
		}

		testutil.Equal(t, buf1.Bytes(), buf2.Bytes())
		testutil.Equal(t, raw, buf1.Bytes())
	}

}
