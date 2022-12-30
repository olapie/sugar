package xslice_test

import (
	"testing"

	"code.olapie.com/sugar/xslice"
	"code.olapie.com/sugar/xtest"
)

func TestReverse(t *testing.T) {
	t.Run("IntArray", func(t *testing.T) {
		a := []int{1, 2, 3, -9, 10, 1, 101}
		xslice.Reverse(a)
		xtest.Equal(t, []int{101, 1, 10, -9, 3, 2, 1}, a)

		a = []int{1}
		xslice.Reverse(a)
		xtest.Equal(t, []int{1}, a)

		a = []int{}
		xslice.Reverse(a)
		xtest.Equal(t, []int{}, a)

		a = []int{1, 3}
		xslice.Reverse(a)
		xtest.Equal(t, []int{3, 1}, a)
	})

	t.Run("StringArray", func(t *testing.T) {
		a := []string{"a", "b", "c", "d"}
		xslice.Reverse(a)
		xtest.Equal(t, []string{"d", "c", "b", "a"}, a)
	})
}
