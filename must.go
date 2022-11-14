package sugar

import (
	"fmt"
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
		panicWithMessages(msgAndArgs...)
	}
}

// MustFalse panics if b is not true
func MustFalse(b bool, msgAndArgs ...any) {
	if b {
		panicWithMessages(msgAndArgs...)
	}
}

// MustError panics if b is not nil
func MustError(err error, msgAndArgs ...any) {
	if err == nil {
		panicWithMessages(msgAndArgs...)
	}
}

// MustNoError panics if b is not nil
func MustNoError(err error, msgAndArgs ...any) {
	if err != nil {
		panicWithMessages(msgAndArgs...)
	}
}

// MustNil panics if v is not nil
func MustNil[T any](v *T, msgAndArgs ...any) {
	if v != nil {
		panicWithMessages(msgAndArgs...)
	}
}

// MustNotNil panics if v is nil
func MustNotNil[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		panicWithMessages(msgAndArgs...)
	}
}

func MustEmptySlice[T any](a []T, msgAndArgs ...any) {
	if len(a) == 0 {
		panicWithMessages(msgAndArgs...)
	}
}

func MustNotEmptySlice[T any](a []T, msgAndArgs ...any) {
	if len(a) == 0 {
		panicWithMessages(msgAndArgs...)
	}
}

func MustEmptyMap[K comparable, V any](m map[K]V, msgAndArgs ...any) {
	if len(m) != 0 {
		panicWithMessages(msgAndArgs...)
	}
}

func MustNotEmptyMap[K comparable, V any](m map[K]V, msgAndArgs ...any) {
	if len(m) == 0 {
		panicWithMessages(msgAndArgs...)
	}
}

func MustEmptyString[S ~string](s S, msgAndArgs ...any) {
	if len(s) != 0 {
		panicWithMessages(msgAndArgs...)
	}
}

func MustNotEmptyString[S ~string](s S, msgAndArgs ...any) {
	if len(s) == 0 {
		panicWithMessages(msgAndArgs...)
	}
}

// Recover recovers from panic and assign message to outErr
// outErr usually is a pointer to return error
// E.g.
//
//	func doSomething() (err error) {
//	    defer Recover(&err)
//	    ...
//	}
func Recover(outErr *error) {
	if v := recover(); v != nil {
		if err, ok := v.(error); ok {
			*outErr = err
		} else {
			*outErr = fmt.Errorf("panic: %v", v)
		}
	}
}

func panicWithMessages(msgAndArgs ...any) {
	n := len(msgAndArgs)
	switch n {
	case 0:
		panic("")
	case 1:
		panic(msgAndArgs[0])
	default:
		if format, ok := msgAndArgs[0].(string); ok {
			panic(fmt.Sprintf(format, msgAndArgs[1:]...))
		}
		panic(fmt.Sprint(msgAndArgs))
	}
}
