package urlx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"reflect"
	"strings"

	"code.olapie.com/sugar/checking"
	"code.olapie.com/sugar/conv"
	"code.olapie.com/sugar/mapping"
	"code.olapie.com/sugar/rtx"
	"code.olapie.com/sugar/stringx"
	"code.olapie.com/sugar/types"
)

func Join(a ...string) string {
	if len(a) == 0 {
		return ""
	}
	// path.Join will convert // to be /
	p := path.Join(a...)
	p = strings.Replace(p, ":/", "://", 1)
	i := strings.Index(p, "://")
	s := p
	if i >= 0 {
		i += 3
		s = p[i:]
		l := strings.Split(s, "/")
		for i, v := range l {
			l[i] = url.PathEscape(v)
		}
		p = p[:i] + path.Join(l...)
	}
	return p
}

func ToValues(i any) (url.Values, error) {
	i = rtx.IndirectToStringerOrError(i)
	if i == nil {
		return nil, errors.New("nil values")
	}
	switch v := i.(type) {
	case url.Values:
		return v, nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("cannot convert %#v of type %T to url.Values", i, i)
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("cannot convert %#v of type %T to url.Values", i, i)
	}
	uv := url.Values{}
	for k, v := range m {
		uv.Set(k, fmt.Sprint(v))
	}
	return uv, nil
}

func MustToValues(i any) url.Values {
	v, err := ToValues(i)
	if err != nil {
		panic(err)
	}
	return v
}

func VarargsToValues(keyAndValues ...any) (url.Values, error) {
	uv := url.Values{}
	keys, vals, err := stringx.FromVarargs(keyAndValues...)
	if err != nil {
		return nil, err
	}
	for i, k := range keys {
		a, err := conv.ToStringSlice(vals[i])
		if err != nil {
			return nil, err
		}
		for _, v := range a {
			if v != "" {
				uv.Add(k, v)
			}
		}
	}
	return uv, nil
}

func MustVarargsToValues(keyAndValues ...any) url.Values {
	v, err := VarargsToValues(keyAndValues...)
	if err != nil {
		panic(err)
	}
	return v
}

func GetPathSegments(endpoint string) (segments []string, paramIndexes []int) {
	segments = strings.Split(endpoint, "/")
	for i, seg := range segments {
		if len(seg) < 3 {
			continue
		}
		n := len(seg)
		if seg[0] == '{' && seg[n-1] == '}' {
			paramIndexes = append(paramIndexes, i)
		}
	}
	return
}

func AppendQuery(urlString string, query url.Values) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return urlString, fmt.Errorf("parse: %w", err)
	}
	q := u.Query()
	for k, vals := range query {
		for _, v := range vals {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func SetPathParams(endpoint string, params any) (string, any) {
	segments, paramIndexes := GetPathSegments(endpoint)
	if len(paramIndexes) == 0 {
		return endpoint, params
	}

	if len(paramIndexes) == 1 && (checking.IsNumber(params) || checking.IsString(params)) {
		str := fmt.Sprint(params)
		if str == "" {
			return endpoint, params
		}
		segments[paramIndexes[0]] = str
		return strings.Join(segments, "/"), nil
	}

	k := rtx.IndirectKind(params)
	if k != reflect.Struct && k != reflect.Map {
		return endpoint, params
	}

	m, ok := params.(types.M)
	if ok {
		m = mapping.Clone(m)
	} else {
		jsonData, err := json.Marshal(params)
		if err != nil {
			return endpoint, params
		}

		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			return endpoint, params
		}
	}

	for _, idx := range paramIndexes {
		name := segments[idx][1 : len(segments[idx])-1]
		param := m.String(name)
		if param != "" {
			segments[idx] = param
			delete(m, name)
		}
	}
	return strings.Join(segments, "/"), m
}
