package rtx

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
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
		panic(fmt.Sprint(msgAndArgs...))
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

func FuncNameOf(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func NameOf(i any) string {
	return reflect.TypeOf(i).Name()
}

func IsNil(i any) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func DeepCopy(dst, src any) error {
	if dst == nil {
		return errors.New("dst cannot be nil")
	}

	if src == nil {
		return errors.New("src cannot be nil")
	}

	dstType := reflect.TypeOf(dst)
	srcType := reflect.TypeOf(src)

	dstKind := IndirectReadableValue(reflect.ValueOf(dst)).Kind()
	srcKind := IndirectReadableValue(reflect.ValueOf(src)).Kind()
	if reflect.PtrTo(srcType).ConvertibleTo(dstType) || (srcKind == dstKind && dstKind == reflect.Struct) {
		err := GobCopy(dst, src)
		if err != nil {
			return fmt.Errorf("json copy: %w", err)
		}
		return nil
	}

	return fmt.Errorf("cannot copy %T to %T", src, dst)
}

func DeepNew(t reflect.Type) reflect.Value {
	v := reflect.New(t)
	e := v.Elem()
	for e.Kind() == reflect.Ptr {
		e.Set(reflect.New(e.Type().Elem()))
		e = e.Elem()
	}

	if e.Kind() != reflect.Struct {
		return v
	}

	for i := 0; i < e.NumField(); i++ {
		ft := e.Type().Field(i)
		if !e.Field(i).CanSet() {
			continue
		}

		switch ft.Type.Kind() {
		case reflect.Ptr:
			e.Field(i).Set(DeepNew(ft.Type.Elem()))
		case reflect.Struct:
			e.Field(i).Set(DeepNew(ft.Type).Elem())
		case reflect.Slice:
			elemVal := DeepNew(ft.Type.Elem()).Elem()
			sliceVal := reflect.New(ft.Type).Elem()
			sliceVal = reflect.Append(sliceVal, elemVal)
			e.Field(i).Set(sliceVal)
		}
	}
	return v
}

func GobCopy(dst, src any) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(src)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	dec := gob.NewDecoder(&b)
	err = dec.Decode(dst)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func JSONCopy(dst, src any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	err = json.Unmarshal(b, dst)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}
