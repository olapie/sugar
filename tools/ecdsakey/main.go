package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"syscall"
	"time"

	"code.olapie.com/sugar/v2"
	"code.olapie.com/sugar/v2/olasec"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	pass := readConfirmedSecret("file password")
	if len(pass) < 8 {
		fmt.Println("Password is too short")
		return
	}
	pk := sugar.MustGet(olasec.GeneratePrivateKey())
	pri := sugar.MustGet(olasec.EncodePrivateKey(pk, pass))
	pub := sugar.MustGet(olasec.EncodePublicKey(&pk.PublicKey))
	name := time.Now().Format("20060102")
	sugar.MustNil(os.WriteFile(name+"-key.png", pri, 0644))
	sugar.MustNil(os.WriteFile(name+"-pub.png", pub, 0644))

	pubKey := sugar.MustGet(olasec.DecodePublicKey(pub))
	priKey := sugar.MustGet(olasec.DecodePrivateKey(pri, pass))

	// Test
	hash := sha256.Sum256([]byte("message: hello"))
	sign := sugar.MustGet(ecdsa.SignASN1(rand.Reader, priKey, hash[:]))
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
		pass = sugar.MustGet(terminal.ReadPassword(syscall.Stdin))
		fmt.Println()
	}
	return string(pass)
}
