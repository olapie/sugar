package rtx

import (
	"fmt"
	"reflect"
)

// Indirect From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// Indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func Indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// IndirectToStringerOrError From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// IndirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
func IndirectToStringerOrError(a any) any {
	if a == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func IndirectReadableValue(v reflect.Value) reflect.Value {
	for (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func IndirectWritableValue(v reflect.Value, populate bool) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if populate && v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			} else {
				break
			}
		}
		v = v.Elem()
	}
	if !v.CanSet() {
		panic(fmt.Sprintf("Cannot set %v", v.Kind()))
	}
	return v
}

func IndirectKind(i any) reflect.Kind {
	switch v := i.(type) {
	case reflect.Type:
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		return v.Kind()
	case reflect.Value:
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		return v.Kind()
	case reflect.Kind:
		return v
	case nil:
		return reflect.Invalid
	default:
		return IndirectKind(reflect.TypeOf(i))
	}
}

func Addr[T any](v T) *T {
	return &v
}

func Dereference[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}

func PanicWithMessages(msgAndArgs ...any) {
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
