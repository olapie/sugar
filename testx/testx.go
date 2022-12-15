package testx

import (
	"crypto/rand"
	"fmt"
	rand2 "math/rand"
	"reflect"
	"testing"
	"time"
)

func Equal(t *testing.T, expected, result any) {
	if reflect.DeepEqual(expected, result) {
		return
	}
	t.Errorf("Expect: %v, got: %v", expected, result)
}

func NotEqual(t *testing.T, expected, result any) {
	if !reflect.DeepEqual(expected, result) {
		return
	}
	t.FailNow()
}

func True(t *testing.T, b bool, args ...any) {
	if !b {
		args = append([]any{"Expected true, got false"}, args...)
		t.Fatal(args...)
	}
}

func False(t *testing.T, b bool, args ...any) {
	if b {
		args = append([]any{"Expected false, got true"}, args...)
		t.Fatal(args...)
	}
}

func NoError(t *testing.T, err error, args ...any) {
	if err != nil {
		args = append([]any{err}, args...)
		t.Fatal(args...)
	}
}

func Error(t *testing.T, err error, args ...any) {
	if err == nil {
		args = append([]any{"Expected error, actually got nil"}, args...)
		t.Fatal(args...)
	}
}

func EmptySlice[E any](t *testing.T, a []E, args ...any) {
	if len(a) != 0 {
		msg := fmt.Sprintf("Expect empty, got %d", len(a))
		args = append([]any{msg}, args...)
		t.Fatal(args...)
	}
}

func NotEmptySlice[E any](t *testing.T, a []E, args ...any) {
	if len(a) == 0 {
		msg := "Expect not empty, got empty"
		args = append([]any{msg}, args...)
		t.Fatal(args...)
	}
}

func EmptyMap[K comparable, V any](t *testing.T, m map[K]V, args ...any) {
	if len(m) != 0 {
		msg := fmt.Sprintf("Expect empty, got %d", len(m))
		args = append([]any{msg}, args...)
		t.Fatal(args...)
	}
}

func NotEmptyMap[K comparable, V any](t *testing.T, m map[K]V, args ...any) {
	if len(m) == 0 {
		msg := "Expect not empty, got empty"
		args = append([]any{msg}, args...)
		t.Fatal(args...)
	}
}

func RandomBytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err == nil {
		return b
	}

	rand2.Seed(time.Now().Unix())
	_, err = rand2.Read(b)
	if err == nil {
		return b
	}

	panic(err)
}

func RandomString(n int) string {
	return RandomStringT[string](n)
}

func RandomStringT[T ~string](n int) T {
	if n <= 0 {
		return ""
	}

	b := RandomBytes(n + 1/2)
	s := fmt.Sprintf("%x", b)
	return T(s[:n])
}
