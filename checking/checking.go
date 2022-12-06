package checking

import (
	"reflect"
	"regexp"

	"code.olapie.com/sugar/conv"
)

type NumberOrString interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string
}

func IsString(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.String
}

func IsBool(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Bool
}

func IsFloat(v any) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func IsUint(v any) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func IsInt(v any) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func IsNumber(v any) bool {
	return IsInt(v) || IsUint(v) || IsFloat(v)
}

func IsEmailAddress(s string) bool {
	addr, _ := conv.ToEmailAddress(s)
	return addr != ""
}

func IsURL(s string) bool {
	u, _ := conv.ToURL(s)
	return u != ""
}

var (
	nickRegexp     = regexp.MustCompile("^[^ \n\r\t\f][^\n\r\t\f]{0,28}[^ \n\r\t\f]$")
	usernameRegexp = regexp.MustCompile("^[a-zA-Z][\\w\\.]{1,19}$")
)

func IsUsername(s string) bool {
	return usernameRegexp.MatchString(s)
}

func IsNickname(s string) bool {
	return nickRegexp.MatchString(s)
}
