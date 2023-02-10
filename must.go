package sugar

import (
	"fmt"

	"code.olapie.com/sugar/v2/rt"
)

// MustGet eliminates nil err and panics if err isn't nil
func MustGet[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetTwo eliminates nil err and panics if err isn't nil
func MustGetTwo[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// MustTrue panics if b is not true
func MustTrue(b bool, msgAndArgs ...any) {
	if !b {
		rt.PanicWithMessages(msgAndArgs...)
	}
}

// MustFalse panics if b is not true
func MustFalse(b bool, msgAndArgs ...any) {
	if b {
		rt.PanicWithMessages(msgAndArgs...)
	}
}

// MustNil panics if v is not nil
func MustNil(v any, msgAndArgs ...any) {
	if v == nil {
		return
	}

	s := fmt.Sprintf("%#v", v)
	if len(msgAndArgs) == 0 {
		rt.PanicWithMessages()
	}
	msgAndArgs[0] = s + " " + fmt.Sprint(msgAndArgs[0])
	rt.PanicWithMessages(msgAndArgs...)
}

// MustNotNil panics if v is nil
func MustNotNil(v any, msgAndArgs ...any) {
	if v == nil {
		rt.PanicWithMessages(msgAndArgs...)
	}
}

// MustNilPointer panics if v is not nil
func MustNilPointer[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		return
	}

	s := fmt.Sprintf("%#v", v)
	if len(msgAndArgs) == 0 {
		rt.PanicWithMessages()
	}
	msgAndArgs[0] = s + " " + fmt.Sprint(msgAndArgs[0])
	rt.PanicWithMessages(msgAndArgs...)
}

// MustNotNilPointer panics if v is nil
func MustNotNilPointer[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		rt.PanicWithMessages(msgAndArgs...)
	}
}
