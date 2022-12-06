package mathx_test

import (
	"testing"

	"code.olapie.com/sugar/mathx"
	"code.olapie.com/sugar/testx"
)

func TestMax(t *testing.T) {
	t.Run("N0", func(t *testing.T) {
		v := mathx.Max[int]()
		testx.Equal(t, 0, v)
	})

	t.Run("N1", func(t *testing.T) {
		v := mathx.Max(10)
		testx.Equal(t, 10, v)
	})

	t.Run("N2", func(t *testing.T) {
		v := mathx.Max(-0.3, 10.9)
		testx.Equal(t, 10.9, v)
	})

	t.Run("N3", func(t *testing.T) {
		v := mathx.Max(-0.3, 10.9, 3.8)
		testx.Equal(t, 10.9, v)
	})
}
