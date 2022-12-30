package xmath_test

import (
	"testing"

	"code.olapie.com/sugar/xmath"
	"code.olapie.com/sugar/xtest"
)

func TestMax(t *testing.T) {
	t.Run("N0", func(t *testing.T) {
		v := xmath.Max[int]()
		xtest.Equal(t, 0, v)
	})

	t.Run("N1", func(t *testing.T) {
		v := xmath.Max(10)
		xtest.Equal(t, 10, v)
	})

	t.Run("N2", func(t *testing.T) {
		v := xmath.Max(-0.3, 10.9)
		xtest.Equal(t, 10.9, v)
	})

	t.Run("N3", func(t *testing.T) {
		v := xmath.Max(-0.3, 10.9, 3.8)
		xtest.Equal(t, 10.9, v)
	})
}
