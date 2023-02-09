package rt_test

import (
	"reflect"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/rt"
	"code.olapie.com/sugar/v2/testutil"
)

func TestIndirectKind(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		k := rt.IndirectKind(nil)
		testutil.Equal(t, reflect.Invalid, k)
	})

	t.Run("Struct", func(t *testing.T) {
		var p time.Time
		k := rt.IndirectKind(p)
		testutil.Equal(t, reflect.Struct, k)
	})

	t.Run("PointerToStruct", func(t *testing.T) {
		var p *time.Time
		k := rt.IndirectKind(p)
		testutil.Equal(t, reflect.Struct, k)
	})

	t.Run("PointerToPointerToStruct", func(t *testing.T) {
		var p **time.Time
		k := rt.IndirectKind(p)
		testutil.Equal(t, reflect.Struct, k)
	})

	t.Run("Map", func(t *testing.T) {
		var p map[string]any
		k := rt.IndirectKind(p)
		testutil.Equal(t, reflect.Map, k)
	})

	t.Run("PointerToMap", func(t *testing.T) {
		var p map[string]any
		k := rt.IndirectKind(p)
		testutil.Equal(t, reflect.Map, k)
	})
}
