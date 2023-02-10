package maths_test

import (
	"testing"

	"code.olapie.com/sugar/v2/testutil"
)

func TestMax(t *testing.T) {
	t.Run("N0", func(t *testing.T) {
		v := maths.Max[int]()
		testutil.Equal(t, 0, v)
	})

	t.Run("N1", func(t *testing.T) {
		v := maths.Max(10)
		testutil.Equal(t, 10, v)
	})

	t.Run("N2", func(t *testing.T) {
		v := maths.Max(-0.3, 10.9)
		testutil.Equal(t, 10.9, v)
	})

	t.Run("N3", func(t *testing.T) {
		v := maths.Max(-0.3, 10.9, 3.8)
		testutil.Equal(t, 10.9, v)
	})
}
