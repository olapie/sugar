package testx

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, expected, result any) {
	if reflect.DeepEqual(expected, result) {
		return
	}

	t.Errorf("expect: %v, got: %v", expected, result)
}

func NotEqual(t *testing.T, expected, result any) {
	if !reflect.DeepEqual(expected, result) {
		return
	}
	t.FailNow()
}

func True(t *testing.T, b bool) {
	if !b {
		t.Fatal("expected true")
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func False(t *testing.T, b bool, msgs ...any) {
	if b {
		t.Fatal(msgs...)
	}
}

func Error(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected error")
	}
}
