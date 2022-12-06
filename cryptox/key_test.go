package cryptox_test

import (
	"bytes"
	"code.olapie.com/sugar/cryptox"
	"code.olapie.com/sugar/hashing"
	"code.olapie.com/sugar/slicing"
	"testing"
	"time"
)

func TestDeriveKey(t *testing.T) {
	key := cryptox.DeriveKey("123", []byte("abc"))
	t.Log(key)
	hash := key.Hash()
	t.Log(hash)
}

func TestAES(t *testing.T) {
	raw := []byte(hashing.SHA1(time.Now().String()))
	password := hashing.SHA1(time.Now().String())
	salt := []byte(hashing.SHA1(time.Now().String()))
	key := cryptox.DeriveKey(password, salt)

	data := slicing.Clone(raw)
	if bytes.Compare(raw, data) != 0 {
		t.FailNow()
	}

	err := key.AES(data, data)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(raw, data) == 0 {
		t.FailNow()
	}

	err = key.AES(data, data)
	if bytes.Compare(raw, data) != 0 {
		t.FailNow()
	}
}
