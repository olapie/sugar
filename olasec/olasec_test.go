package olasec_test

import (
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"code.olapie.com/sugar/olasec"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"
)

func TestEncrypt(t *testing.T) {
	password := hashing.SHA1(time.Now().String())
	testEncrypt(t, 1<<4+9, password)
	testEncrypt(t, 1<<24, password)
}

func testEncrypt(t *testing.T, size int, password string) {
	raw := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, raw[:])
	if err != nil {
		t.Fatal(err)
	}

	enc, err := olasec.Encrypt(raw[:], password)
	t.Log(enc[:30])
	testx.NoError(t, err)
	testx.True(t, olasec.IsEncrypted(enc))
	dec, err := olasec.DecryptBytes(enc, password)
	testx.NoError(t, err)
	testx.False(t, olasec.IsEncrypted(dec), dec[:olasec.HeaderSize])
	testx.Equal(t, raw, dec)
}

func TestEncryptFile(t *testing.T) {
	err := os.MkdirAll("testdata", 0755)
	if err != nil {
		t.Fatal(err)
	}

	rawFilename := "testdata/rawfile"
	largeFilename := "testdata/largefile"
	t.Cleanup(func() {
		os.RemoveAll(rawFilename)
		os.RemoveAll(largeFilename)
	})

	password := hashing.SHA1(time.Now().String())
	var raw [32]byte
	n, err := io.ReadFull(rand.Reader, raw[:])
	testx.NoError(t, err)
	t.Log(n, raw)
	f, err := os.OpenFile(rawFilename, os.O_CREATE|os.O_WRONLY, 0644)
	testx.NoError(t, err)

	_, err = f.Write(raw[:])
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rawFilename)
	testEncryptFile(t, rawFilename, password)

	var large [32 * 1024 * 1024]byte
	n, err = io.ReadFull(rand.Reader, large[:])
	testx.NoError(t, err)

	f, err = os.OpenFile(largeFilename, os.O_CREATE|os.O_WRONLY, 0644)
	testx.NoError(t, err)

	_, err = f.Write(large[:])
	f.Close()
	testx.NoError(t, err)

	testEncryptFile(t, largeFilename, password)
}

func testEncryptFile(t *testing.T, rawFilename string, password string) {
	encFilename := rawFilename + ".enc" + filepath.Ext(rawFilename)
	decFilename := rawFilename + ".dec" + filepath.Ext(rawFilename)
	t.Cleanup(func() {
		os.RemoveAll(decFilename)
		os.RemoveAll(encFilename)
	})
	err := olasec.EncryptFile(olasec.SF(rawFilename), olasec.DF(encFilename), password)
	testx.NoError(t, err)

	testx.True(t, olasec.IsEncryptedFile(encFilename))
	testx.False(t, olasec.IsEncryptedFile(rawFilename))
	err = olasec.DecryptFile(olasec.SF(encFilename), olasec.DF(decFilename), password)
	testx.NoError(t, err)
	raw, err := os.ReadFile(rawFilename)
	testx.NoError(t, err)

	enc, err := os.ReadFile(encFilename)
	testx.NoError(t, err)
	testx.NotEqual(t, raw, enc)

	dec, err := os.ReadFile(decFilename)
	testx.NoError(t, err)
	testx.Equal(t, raw, dec)

	valid := olasec.ValidateFilePassword(encFilename, password)
	testx.True(t, valid)
}
