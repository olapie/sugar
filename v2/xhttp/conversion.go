package xhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"code.olapie.com/sugar/xtype"
)

type MapOutput interface {
	map[string]string | map[string]any | xtype.M
}

type MapInput interface {
	http.Header | []*http.Cookie | url.Values
}

func ToMap[IN MapInput, OUT MapOutput](input IN) OUT {
	var out OUT
	var res any
	_, isMapString := any(out).(map[string]string)
	switch v := any(input).(type) {
	case http.Header:
		if isMapString {
			res = headerToMapString(v)
		} else {
			res = headerToMap(v)
		}
	case []*http.Cookie:
		if isMapString {
			res = cookiesToMapString(v)
		} else {
			res = cookiesToMap(v)
		}
	case url.Values:
		if isMapString {
			res = valuesToMapString(v)
		} else {
			res = valuesToMap(v)
		}
	}

	if res != nil {
		reflect.ValueOf(&out).Elem().Set(reflect.ValueOf(res))
	}
	return out
}

func ToMapAny[T MapInput](v T) map[string]any {
	return ToMap[T, map[string]any](v)
}

func ToMapString[T MapInput](v T) map[string]string {
	return ToMap[T, map[string]string](v)
}

func ToM[T MapInput](v T) xtype.M {
	return ToMap[T, xtype.M](v)
}

func headerToMap(h http.Header) map[string]any {
	m := make(map[string]any, len(h))
	for k, v := range h {
		if len(v) == 0 {
			continue
		}

		if len(v) == 1 {
			m[k] = v[0]
			continue
		}
		m[k] = v
		k = strings.ToLower(k)
		if strings.HasPrefix(k, "x-") {
			k = k[2:]
			k = strings.Replace(k, "-", "_", -1)
			if _, ok := m[k]; !ok {
				m[k] = v
			}
		}
	}
	return m
}

func headerToMapString(h http.Header) map[string]string {
	m := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) == 0 {
			continue
		}
		m[k] = v[0]
		k = strings.ToLower(k)
		if strings.HasPrefix(k, "x-") {
			k = k[2:]
			k = strings.Replace(k, "-", "_", -1)
			if _, ok := m[k]; !ok {
				m[k] = v[0]
			}
		}
	}
	return m
}

func cookiesToMapString(cookies []*http.Cookie) map[string]string {
	m := make(map[string]string, len(cookies))
	for _, c := range cookies {
		m[c.Name] = c.Value
	}
	return m
}

func cookiesToMap(cookies []*http.Cookie) map[string]any {
	m := make(map[string]any, len(cookies))
	for _, c := range cookies {
		m[c.Name] = c.Value
	}
	return m
}

func valuesToMap(values url.Values) map[string]any {
	m := map[string]any{}
	for k, va := range values {
		isArray := strings.HasSuffix(k, "[]")
		if isArray {
			k = k[0 : len(k)-2]
			if k == "" {
				continue
			}

			if len(va) == 1 {
				va = strings.Split(va[0], ",")
			}
		}

		if len(va) == 0 {
			continue
		}

		k = strings.ToLower(k)
		if isArray || len(va) > 1 {
			// value is an array or expected to be an array
			m[k] = va
		} else {
			m[k] = va[0]
		}
	}

	if jsonStr, _ := m["json"].(string); jsonStr != "" {
		var j map[string]any
		err := json.Unmarshal([]byte(jsonStr), &j)
		if err == nil {
			for k, v := range m {
				j[k] = v
			}
			m = j
		}
	}
	return m
}

func valuesToMapString(values url.Values) map[string]string {
	m := valuesToMap(values)
	res := make(map[string]string, len(m))
	for k, v := range m {
		if s, ok := v.(string); ok {
			res[k] = s
		} else {
			res[k] = fmt.Sprint(v)
		}
	}
	return res
}
