package stringx

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"code.olapie.com/sugar/conv"

	"code.olapie.com/sugar/rtx"
)

func Join(a any, sep string) (string, error) {
	l, err := conv.ToStringSlice(a)
	if err != nil {
		return "", err
	}
	return strings.Join(l, sep), nil
}

func TrimSpace[T ~string](s T) T {
	return T(strings.TrimSpace(string(s)))
}

var whitespaceRegexp = regexp.MustCompile(`[ \t\n\r]+`)
var bulletRegexp = regexp.MustCompile(`[\d\.\*]*`)

// Squish returns the string
// first removing all whitespace on both ends of the string,
// and then changing remaining consecutive whitespace groups into one space each.
func Squish[T ~string](s T) T {
	str := strings.TrimSpace(string(s))
	str = whitespaceRegexp.ReplaceAllString(str, " ")
	return T(str)
}

func SquishFields(i any) {
	squishFields(reflect.ValueOf(i))
}

func squishFields(v reflect.Value) {
	v = rtx.IndirectReadableValue(v)
	switch v.Kind() {
	case reflect.Struct:
		squishStructFields(v)
	case reflect.String:
		fmt.Println(v.CanSet(), v.String())
		if v.CanSet() {
			v.SetString(Squish(v.String()))
		}
	default:
		break
	}
}

func squishStructFields(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(Squish(fv.String()))
		case reflect.Struct:
			squishStructFields(fv)
		case reflect.Ptr, reflect.Interface:
			if fv.IsNil() {
				break
			}
			squishFields(fv.Elem())
		default:
			break
		}
	}
}

func RemoveAllSpaces[T ~string](s T) T {
	return T(strings.ReplaceAll(string(Squish(s)), " ", ""))
}

func RemoveBullet[T ~string](s T) T {
	s = Squish(s)
	a := strings.Split(string(s), " ")

	if len(a) == 0 {
		return ""
	}
	a[0] = bulletRegexp.ReplaceAllString(a[0], "")
	if a[0] == "" {
		a = a[1:]
	}
	return T(strings.Join(a, " "))
}

func FromVarargs(keyValues ...any) (keys []string, values []any, err error) {
	n := len(keyValues)
	if n%2 != 0 {
		err = errors.New("keyValues should be pairs of (string, any)")
		return
	}

	keys, values = make([]string, 0, n/2), make([]any, 0, n/2)
	for i := 0; i < n/2; i++ {
		if k, ok := keyValues[2*i].(string); !ok {
			err = fmt.Errorf("keyValues[%d] isn't convertible to string", i)
			return
		} else if keyValues[2*i+1] == nil {
			err = fmt.Errorf("keyValues[%d] is nil", 2*i+1)
			return
		} else {
			keys = append(keys, k)
			values = append(values, keyValues[2*i+1])
		}
	}
	return
}
