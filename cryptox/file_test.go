package cryptox_test

import (
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/testx"

	"code.olapie.com/sugar/cryptox"
)

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
	err := cryptox.EncryptFile(cryptox.Destination(encFilename), cryptox.Source(rawFilename), password)
	testx.NoError(t, err)

	testx.True(t, cryptox.IsEncrypted(encFilename))
	testx.False(t, cryptox.IsEncrypted(rawFilename))
	err = cryptox.DecryptFile(cryptox.Destination(decFilename), cryptox.Source(encFilename), password)
	testx.NoError(t, err)
	raw, err := os.ReadFile(rawFilename)
	testx.NoError(t, err)

	enc, err := os.ReadFile(encFilename)
	testx.NoError(t, err)
	testx.NotEqual(t, raw, enc)

	dec, err := os.ReadFile(decFilename)
	testx.NoError(t, err)
	testx.Equal(t, raw, dec)
}
