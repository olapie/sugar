package checking

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"code.olapie.com/sugar/rtx"

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

func IsDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return false
	}
	return true
}

type Validator interface {
	Validate() error
}

func Validate(i any) error {
	if v, ok := i.(Validator); ok {
		return v.Validate()
	}

	v := reflect.ValueOf(i)
	if v.IsValid() && (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
		v = v.Elem()
		if v.CanInterface() {
			if va, ok := v.Interface().(Validator); ok {
				return va.Validate()
			}
		}
	}

	v = rtx.IndirectReadableValue(v)
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for j := 0; j < v.NumField(); j++ {
			if !rtx.IsExported(t.Field(j).Name) {
				continue
			}
			if err := Validate(v.Field(j).Interface()); err != nil {
				return fmt.Errorf("%s:%w", t.Field(j).Name, err)
			}
		}
	}
	return nil
}
