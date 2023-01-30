package mob

import (
	"code.olapie.com/sugar/mob/nomobile"
	"code.olapie.com/sugar/v2/xtype"
)

type Map struct {
	m xtype.M
}

func NewMap() *Map {
	return &Map{m: xtype.M{}}
}

func (m *Map) GetInt64(key string) int64 {
	return m.m.Int64(key)
}

func (m *Map) GetFloat64(key string) float64 {
	return m.m.Float64(key)
}

func (m *Map) GetString(key string) string {
	return m.m.String(key)
}

func (m *Map) GetBool(key string) bool {
	return m.m.Bool(key)
}

func (m *Map) GetInt64List(key string) *Int64List {
	switch v := m.m[key].(type) {
	case []int64:
		l := nomobile.NewList(v)
		return &Int64List{
			List: *l,
		}
	case *Int64List:
		return v
	default:
		return nil
	}
}

func (m *Map) GetStringList(key string) *StringList {
	switch v := m.m[key].(type) {
	case []string:
		l := nomobile.NewList(v)
		return &StringList{
			List: *l,
		}
	case *StringList:
		return v
	default:
		return nil
	}
}

func (m *Map) SetInt64(key string, val int64) {
	m.m[key] = val
}

func (m *Map) SetFloat64(key string, val float64) {
	m.m[key] = val
}

func (m *Map) SetString(key string, val string) {
	m.m[key] = val
}

func (m *Map) SetBool(key string, val bool) {
	m.m[key] = val
}

func (m *Map) SetInt64List(key string, val *Int64List) {
	m.m[key] = val.Elements()
}

func (m *Map) SetStringList(key string, val *StringList) {
	m.m[key] = val.Elements()
}
