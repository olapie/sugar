package xrandom

import (
	"code.olapie.com/sugar/v2/must"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
)

var defaultOTPGenerator *otpGenerator

const digitAlphaString = "23456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

type otpGenerator struct {
	mu       sync.Mutex
	lenToMax map[int]*big.Int
}

func newOTPGenerator() *otpGenerator {
	return &otpGenerator{
		lenToMax: map[int]*big.Int{},
	}
}

func (s *otpGenerator) generateDigitCode(length int) (string, error) {
	max := s.getMaxRand(length)
	var code string
	for len(code) < length {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("rand.Int: %w", err)
		}
		code += n.String()
	}
	return code[:length], nil
}

func (s *otpGenerator) generateAlphaDigitCode(length int) (string, error) {
	max := big.NewInt(int64(len(digitAlphaString)))
	a := make([]byte, length)
	for i := range a {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("rand.Int: %w", err)
		}
		a[i] = digitAlphaString[index.Int64()]
	}
	return string(a), nil
}

func (s *otpGenerator) getMaxRand(length int) *big.Int {
	max := s.lenToMax[length]
	if max != nil {
		return max
	}
	s.mu.Lock()
	max = s.lenToMax[length]
	if max != nil {
		return max
	}
	max = big.NewInt(1)
	n10 := big.NewInt(10)
	// length=6 => max=1e7
	// length=8 => max=1e9
	for i := 0; i < length; i++ {
		max.Mul(max, n10)
	}
	s.lenToMax[length] = max
	s.mu.Unlock()
	return max
}

func DigitCode(n int) string {
	if defaultOTPGenerator == nil {
		defaultOTPGenerator = newOTPGenerator()
	}
	return must.Get(defaultOTPGenerator.generateDigitCode(n))
}

func AlphaDigitCode(n int) string {
	if defaultOTPGenerator == nil {
		defaultOTPGenerator = newOTPGenerator()
	}
	return must.Get(defaultOTPGenerator.generateAlphaDigitCode(n))
}
