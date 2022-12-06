package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/mail"
	"net/url"
	"reflect"
	"strings"
	"time"

	"code.olapie.com/sugar/conv"
	"code.olapie.com/sugar/timing"
)

// M is a special map which provides convenient methods
type M map[string]any

func (m M) Slice(key string) []any {
	value := m[key]
	if value == nil {
		return []any{}
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		length := v.Len()
		var values = make([]any, length)
		for i := 0; i < length; i++ {
			values[i] = v.Index(i).Interface()
		}
		return values
	default:
		return []any{value}
	}
}

func (m M) Get(key string) any {
	value := m[key]
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() > 0 {
			return v.Index(0).Interface()
		}
		return nil
	default:
		return value
	}
}

func (m M) Contains(key string) bool {
	return m.Get(key) != nil
}

func (m M) ContainsString(key string) bool {
	switch m.Get(key).(type) {
	case string:
		return true
	default:
		return false
	}
}

func (m M) String(key string) string {
	switch v := m.Get(key).(type) {
	case string:
		return v
	case json.Number:
		return string(v)
	default:
		return ""
	}
}

func (m M) TrimmedString(key string) string {
	switch v := m.Get(key).(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return string(v)
	default:
		return ""
	}
}

func (m M) DefaultString(key string, defaultValue string) string {
	switch v := m.Get(key).(type) {
	case string:
		return v
	default:
		return defaultValue
	}
}

func (m M) MustString(key string) string {
	switch v := m.Get(key).(type) {
	case string:
		return v
	default:
		panic("No string value for key:" + key)
	}
}

func (m M) StringSlice(key string) []string {
	_, found := m[key]
	if !found {
		return nil
	}

	values := m.Slice(key)
	var result []string
	for _, v := range values {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}

	return result
}

func (m M) ContainsBool(key string) bool {
	_, err := conv.ToBool(m.Get(key))
	return err == nil
}

func (m M) Bool(key string) bool {
	v, _ := conv.ToBool(m.Get(key))
	return v
}

func (m M) DefaultBool(key string, defaultValue bool) bool {
	if v, err := conv.ToBool(m.Get(key)); err == nil {
		return v
	}
	return defaultValue
}

func (m M) MustBool(key string) bool {
	if v, err := conv.ToBool(m.Get(key)); err == nil {
		return v
	}
	panic("No bool value for key:" + key)
}

func (m M) Int(key string) int {
	v, _ := conv.ToInt64(m.Get(key))
	return int(v)
}

func (m M) DefaultInt(key string, defaultVal int) int {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return int(v)
	}
	return defaultVal
}

func (m M) MustInt(key string) int {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return int(v)
	}
	panic("No int value for key:" + key)
}

func (m M) IntSlice(key string) []int {
	l, _ := conv.ToIntSlice(m.Slice(key))
	return l
}

func (m M) ContainsInt64(key string) bool {
	_, err := conv.ToInt64(m.Get(key))
	return err == nil
}

func (m M) Int64(key string) int64 {
	v, _ := conv.ToInt64(m.Get(key))
	return v
}

func (m M) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return v
	}
	return defaultVal
}

func (m M) MustInt64(key string) int64 {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return v
	}
	panic("No int64 value for key:" + key)
}

func (m M) Int64Slice(key string) []int64 {
	values := m.Slice(key)
	var result []int64
	for _, v := range values {
		i, e := conv.ToInt64(v)
		if e == nil {
			result = append(result, i)
		}
	}

	return result
}

func (m M) ContainsFloat64(key string) bool {
	_, err := conv.ToFloat64(m.Get(key))
	return err == nil
}

func (m M) Float64(key string) float64 {
	v, _ := conv.ToFloat64(m.Get(key))
	return v
}

func (m M) DefaultFloat64(key string, defaultValue float64) float64 {
	if v, err := conv.ToFloat64(m.Get(key)); err == nil {
		return v
	}
	return defaultValue
}

func (m M) MustFloat64(key string) float64 {
	if v, err := conv.ToFloat64(m.Get(key)); err == nil {
		return v
	}
	panic("No float64 value for key:" + key)
}

func (m M) Float64Slice(key string) []float64 {
	values := m.Slice(key)
	var result []float64
	for _, val := range values {
		i, e := conv.ToFloat64(val)
		if e == nil {
			result = append(result, i)
		}
	}

	return result
}

func (m M) BigInt(key string) *big.Int {
	s := m.String(key)
	n := big.NewInt(0)
	_, ok := n.SetString(s, 10)
	if !ok {
		_, ok = n.SetString(s, 16)
	}

	if ok {
		return n
	}

	if k, err := conv.ToInt64(m[key]); err == nil {
		return big.NewInt(k)
	}

	return nil
}

func (m M) DefaultBigInt(key string, defaultVal *big.Int) *big.Int {
	if n := m.BigInt(key); n != nil {
		return n
	}
	return defaultVal
}

func (m M) MustBigInt(key string) *big.Int {
	if n := m.BigInt(key); n != nil {
		return n
	}
	panic("No big.Int64 value for key:" + key)
}

func (m M) BigFloat(key string) *big.Float {
	s := m.String(key)
	n := big.NewFloat(0)
	_, ok := n.SetString(s)
	if ok {
		return n
	}

	if k, err := conv.ToFloat64(m[key]); err == nil {
		return big.NewFloat(k)
	}

	return nil
}

func (m M) DefaultBigFloat(key string, defaultVal *big.Float) *big.Float {
	if n := m.BigFloat(key); n != nil {
		return n
	}
	return defaultVal
}

func (m M) MustBigFloat(key string) *big.Float {
	if n := m.BigFloat(key); n != nil {
		return n
	}
	panic("No big.Float64 value for key:" + key)
}

func (m M) Map(key string) M {
	switch v := m.Get(key).(type) {
	case M:
		return v
	case map[string]any:
		return v
	default:
		return M{}
	}
}

func (m M) Time(key string) (time.Time, bool) {
	return m.TimeInLocation(key, time.UTC)
}

func (m M) TimeInLocation(key string, loc *time.Location) (time.Time, bool) {
	s := strings.TrimSpace(m.String(key))
	if s == "" {
		return time.Time{}, false
	}
	d, err := timing.ToTimeInLocation(s, loc)
	if err != nil {
		return time.Time{}, false
	}
	return d, true
}

func (m M) EmailAddress(key string) *mail.Address {
	s := m.String(key)
	s = strings.TrimSpace(s)
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return nil
	}
	return addr
}

func (m M) URL(key string) string {
	s := m.String(key)
	s = strings.TrimSpace(s)
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}

	if u.Scheme == "" || u.Host == "" {
		return ""
	}
	return s
}

func (m M) SetNX(k string, v any) {
	if v == nil {
		return
	}
	if _, ok := m[k]; ok {
		return
	}
	m[k] = v
}

func (m M) AddMap(val M) {
	for k, v := range val {
		m.SetNX(k, v)
	}
}

func (m M) AddStruct(s any) error {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

func (m M) JSONString() string {
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

func (m M) Remove(keys ...string) {
	for k := range m {
		if indexOfStr(keys, k) < 0 {
			delete(m, k)
		}
	}
}

func (m M) Keep(keys ...string) {
	for k := range m {
		if indexOfStr(keys, k) < 0 {
			delete(m, k)
		}
	}
}
func (m M) ContainsID(key string) bool {
	_, err := conv.ToInt64(m.Get(key))
	return err == nil
}

func (m M) ID(key string) ID {
	v, _ := conv.ToInt64(m.Get(key))
	return ID(v)
}

func (m M) DefaultID(key string, defaultVal ID) ID {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return ID(v)
	}
	return defaultVal
}

func (m M) MustID(key string) ID {
	if v, err := conv.ToInt64(m.Get(key)); err == nil {
		return ID(v)
	}
	panic("No ID value for key:" + key)
}

func (m M) IDSlice(key string) []ID {
	values := m.Slice(key)
	var result []ID
	for _, v := range values {
		i, e := conv.ToInt64(v)
		if e == nil {
			result = append(result, ID(i))
		}
	}

	return result
}

func indexOfStr(l []string, s string) int {
	for i, str := range l {
		if s == str {
			return i
		}
	}
	return -1
}
