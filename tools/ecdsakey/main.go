package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"syscall"
	"time"

	"code.olapie.com/sugar/cryptox"
	"code.olapie.com/sugar/must"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	pass := readConfirmedSecret("file password")
	if len(pass) < 8 {
		fmt.Println("Password is too short")
		return
	}
	pk := must.Get(cryptox.GeneratePrivateKey())
	pri := must.Get(cryptox.EncodePrivateKey(pk, pass))
	pub := must.Get(cryptox.EncodePublicKey(&pk.PublicKey))
	name := time.Now().Format("20060102")
	must.NoError(os.WriteFile(name+"-key.png", pri, 0644))
	must.NoError(os.WriteFile(name+"-pub.png", pub, 0644))

	pubKey := must.Get(cryptox.DecodePublicKey(pub))
	priKey := must.Get(cryptox.DecodePrivateKey(pri, pass))

	// Test
	hash := sha256.Sum256([]byte("message: hello"))
	sign := must.Get(ecdsa.SignASN1(rand.Reader, priKey, hash[:]))
	ok := ecdsa.VerifyASN1(pubKey, hash[:], sign)
	if !ok {
		fmt.Println("Test failed")
	}

	hash[0] = 20
	ok = ecdsa.VerifyASN1(pubKey, hash[:], sign)
	if ok {
		fmt.Println("Test failed")
	}
	fmt.Println("Test succeeded")
}

func readConfirmedSecret(name string) string {
	pass1 := readNonEmptyPassword(fmt.Sprintf("Enter %s: ", name))
	pass2 := readNonEmptyPassword(fmt.Sprintf("Repeat %s: ", name))
	if pass1 != pass2 {
		fmt.Println("Inputs mismatch")
		return ""
	}
	return pass1
}

func readNonEmptyPassword(msg ...any) string {
	var pass []byte
	for len(pass) == 0 {
		fmt.Print(msg...)
		pass = must.Get(terminal.ReadPassword(syscall.Stdin))
		fmt.Println()
	}
	return string(pass)
}
