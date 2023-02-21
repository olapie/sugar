package urlutil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"code.olapie.com/sugar/v2/types"
)

func ToM(v url.Values) types.M {
	return ToMap(v)
}

func ToMap(values url.Values) map[string]any {
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

func ToMapString(values url.Values) map[string]string {
	m := ToMap(values)
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
