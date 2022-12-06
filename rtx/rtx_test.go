package rtx_test

import (
	"reflect"
	"testing"
	"time"

	"code.olapie.com/sugar/rtx"
	"code.olapie.com/sugar/testx"
)

func TestIndirectKind(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		k := rtx.IndirectKind(nil)
		testx.Equal(t, reflect.Invalid, k)
	})

	t.Run("Struct", func(t *testing.T) {
		var p time.Time
		k := rtx.IndirectKind(p)
		testx.Equal(t, reflect.Struct, k)
	})

	t.Run("PointerToStruct", func(t *testing.T) {
		var p *time.Time
		k := rtx.IndirectKind(p)
		testx.Equal(t, reflect.Struct, k)
	})

	t.Run("PointerToPointerToStruct", func(t *testing.T) {
		var p **time.Time
		k := rtx.IndirectKind(p)
		testx.Equal(t, reflect.Struct, k)
	})

	t.Run("Map", func(t *testing.T) {
		var p map[string]any
		k := rtx.IndirectKind(p)
		testx.Equal(t, reflect.Map, k)
	})

	t.Run("PointerToMap", func(t *testing.T) {
		var p map[string]any
		k := rtx.IndirectKind(p)
		testx.Equal(t, reflect.Map, k)
	})
}
