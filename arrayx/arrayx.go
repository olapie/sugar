package arrayx

import (
	"reflect"

	"code.olapie.com/sugar/rtx"
)

func Reverse(a any) bool {
	a = rtx.Indirect(a)
	if a == nil {
		return false
	}
	v := reflect.ValueOf(a)
	if v.IsNil() || !v.IsValid() || v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return false
	}

	for i, j := 0, v.Len()-1; i < j; i, j = i+1, j-1 {
		vi, vj := v.Index(i), v.Index(j)
		tmp := vi.Interface()
		if !vi.CanSet() {
			return false
		}
		vi.Set(vj)
		vj.Set(reflect.ValueOf(tmp))
	}
	return true
}
