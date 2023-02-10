package randoms

import (
	"testing"
)

func TestVerificationService_DigitCode(t *testing.T) {
	for n := 6; n < 20; n++ {
		t.Log(DigitCode(n))
	}
}

func TestVerificationService_AlphaDigitCode(t *testing.T) {
	for n := 4; n < 20; n++ {
		t.Log(AlphaDigitCode(n))
	}
}
