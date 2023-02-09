package contacts

import (
	"fmt"
	"reflect"
	"regexp"

	"code.olapie.com/sugar/v2/conv"
	"code.olapie.com/sugar/v2/rt"
)

func IsEmailAddress(s string) bool {
	addr, _ := conv.ToEmailAddress(s)
	return addr != ""
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

	v = rt.IndirectReadableValue(v)
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for j := 0; j < v.NumField(); j++ {
			if !rt.IsExported(t.Field(j).Name) {
				continue
			}
			if err := Validate(v.Field(j).Interface()); err != nil {
				return fmt.Errorf("%s:%w", t.Field(j).Name, err)
			}
		}
	}
	return nil
}
