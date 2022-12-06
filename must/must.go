package must

import "code.olapie.com/sugar/rtx"

// Get eliminates nil err and panics if err isn't nil
func Get[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// GetTwo eliminates nil err and panics if err isn't nil
func GetTwo[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// True panics if b is not true
func True(b bool, msgAndArgs ...any) {
	if !b {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// False panics if b is not true
func False(b bool, msgAndArgs ...any) {
	if b {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// Error panics if b is not nil
func Error(err error, msgAndArgs ...any) {
	if err == nil {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// NoError panics if b is not nil
func NoError(err error, msgAndArgs ...any) {
	if err != nil {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// Nil panics if v is not nil
func Nil[T any](v *T, msgAndArgs ...any) {
	if v != nil {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// NotNil panics if v is nil
func NotNil[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func EmptySlice[T any](a []T, msgAndArgs ...any) {
	if len(a) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func NotEmptySlice[T any](a []T, msgAndArgs ...any) {
	if len(a) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func EmptyString[S ~string](s S, msgAndArgs ...any) {
	if len(s) != 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func NotEmptyString[S ~string](s S, msgAndArgs ...any) {
	if len(s) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func EmptyMap[K comparable, V any](m map[K]V, msgAndArgs ...any) {
	if len(m) != 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func NotEmptyMap[K comparable, V any](m map[K]V, msgAndArgs ...any) {
	if len(m) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}
