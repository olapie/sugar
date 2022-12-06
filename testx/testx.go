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
