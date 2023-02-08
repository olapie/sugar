package xcontext

import (
	"context"
	"fmt"
	"reflect"

	"code.olapie.com/sugar/v2/xruntime"
)

type keyType int

// Context keys
const (
	keyStart keyType = iota
	keyLogin
	keySudo
	keyLogger
	keyRequestMetadata

	keyEnd
)

func Detach(ctx context.Context) context.Context {
	newCtx := context.Background()
	for k := keyStart; k < keyEnd; k++ {
		if v := ctx.Value(k); v != nil {
			newCtx = context.WithValue(newCtx, k, v)
		}
	}
	return newCtx
}

func GetLogin[T comparable](ctx context.Context) T {
	v := ctx.Value(keyLogin)
	if v == nil {
		var zero T
		return zero
	}

	if actual, ok := v.(T); ok {
		return actual
	}

	actualVal := reflect.ValueOf(v)
	var expect T
	if xruntime.IsInt(actualVal) || xruntime.IsUint(actualVal) || xruntime.IsFloat(actualVal) {
		if reflect.ValueOf(expect).Kind() == reflect.String {
			reflect.ValueOf(&expect).Elem().SetString(fmt.Sprint(v))
			return expect
		}
	}

	expectType := reflect.TypeOf(expect)
	if actualVal.Type().ConvertibleTo(expectType) {
		defer func() {
			if msg := recover(); msg != nil {
				fmt.Printf("[sugar/v2/contexts] GetLogin: %v\n", msg)
			}
		}()
		reflect.ValueOf(&expect).Elem().Set(reflect.ValueOf(v).Convert(expectType))
	}
	return expect
}

func WithLogin[T comparable](ctx context.Context, v T) context.Context {
	var zero T
	if v == zero {
		if ctx.Value(keyLogin) == nil {
			return ctx
		}
		return context.WithValue(ctx, keyLogin, nil)
	}
	return context.WithValue(ctx, keyLogin, v)
}

func HasLogin(ctx context.Context) bool {
	return ctx.Value(keyLogin) != nil
}

func WithSudo(ctx context.Context) context.Context {
	return context.WithValue(ctx, keySudo, true)
}

func IsSudo(ctx context.Context) bool {
	b, _ := ctx.Value(keySudo).(bool)
	return b
}

func GetLogger[T any](ctx context.Context) T {
	v, _ := ctx.Value(keyLogger).(T)
	return v
}

func WithLogger[T any](ctx context.Context, logger T) context.Context {
	var l T
	if any(logger) == any(l) {
		return ctx
	}

	v := ctx.Value(keyLogger)
	if v != nil {
		if reflect.TypeOf(v).AssignableTo(reflect.TypeOf(l)) {
			panic(fmt.Sprintf("cannot override existing logger[%T] with different type [%T]", v, l))
		}
	}
	return context.WithValue(ctx, keyLogger, logger)
}
