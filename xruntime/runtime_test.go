package xruntime_test

import (
	"reflect"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/xruntime"
	"code.olapie.com/sugar/v2/xtest"
)

func TestIndirectKind(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		k := xruntime.IndirectKind(nil)
		xtest.Equal(t, reflect.Invalid, k)
	})

	t.Run("Struct", func(t *testing.T) {
		var p time.Time
		k := xruntime.IndirectKind(p)
		xtest.Equal(t, reflect.Struct, k)
	})

	t.Run("PointerToStruct", func(t *testing.T) {
		var p *time.Time
		k := xruntime.IndirectKind(p)
		xtest.Equal(t, reflect.Struct, k)
	})

	t.Run("PointerToPointerToStruct", func(t *testing.T) {
		var p **time.Time
		k := xruntime.IndirectKind(p)
		xtest.Equal(t, reflect.Struct, k)
	})

	t.Run("Map", func(t *testing.T) {
		var p map[string]any
		k := xruntime.IndirectKind(p)
		xtest.Equal(t, reflect.Map, k)
	})

	t.Run("PointerToMap", func(t *testing.T) {
		var p map[string]any
		k := xruntime.IndirectKind(p)
		xtest.Equal(t, reflect.Map, k)
	})
}
