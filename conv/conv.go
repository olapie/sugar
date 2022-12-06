package conv

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"math"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"code.olapie.com/sugar/errorx"
	"code.olapie.com/sugar/rtx"
)

// ToBool converts i to bool
// i can be bool, integer or string
func ToBool(i any) (bool, error) {
	i = rtx.Indirect(i)
	switch v := i.(type) {
	case bool:
		return v, nil
	case nil:
		return false, errorx.NotExist
	case string:
		return strconv.ParseBool(v)
	}

	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.String:
		return strconv.ParseBool(v.String())
	}

	n, err := parseInt64(i)
	if err != nil {
		return false, fmt.Errorf("cannot convert %#v of type %T to bool", i, i)
	}
	return n != 0, nil
}

// ToBoolSlice converts i to []bool
// i is an array or slice with elements convertiable to bool
func ToBoolSlice(i any) ([]bool, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]bool); ok {
		return l, nil
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %#v of type %T to []bool", i, i)
	}
	num := v.Len()
	res := make([]bool, num)
	var err error
	for j := 0; j < num; j++ {
		res[j], err = ToBool(v.Index(j).Interface())
		if err != nil {
			return nil, fmt.Errorf("convert index %d: %w", i, err)
		}
	}
	return res, nil
}

func ToFloat32(i any) (float32, error) {
	v, err := ToFloat64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to float32", i, i)
	}
	if v > math.MaxFloat32 || v < -math.MaxFloat32 {
		return 0, strconv.ErrRange
	}
	return float32(v), nil
}

func ToFloat64(i any) (float64, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return 0, errorx.NotExist
	}

	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	v := reflect.ValueOf(i)
	if rtx.IsInt(v) {
		return float64(v.Int()), nil
	}

	if rtx.IsUint(v) {
		return float64(v.Uint()), nil
	}

	if rtx.IsFloat(v) {
		return v.Float(), nil
	}

	switch v.Kind() {
	case reflect.String:
		return strconv.ParseFloat(v.String(), 64)
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %#v of type %T to float64", i, i)
	}
}

func ToFloat32Slice(i any) ([]float32, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]float32); ok {
		return l, nil
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %#v of type %T to []float32", i, i)
	}
	num := v.Len()
	res := make([]float32, num)
	var err error
	for j := 0; j < num; j++ {
		res[j], err = ToFloat32(v.Index(j).Interface())
		if err != nil {
			return nil, fmt.Errorf("convert index %d: %w", i, err)
		}
	}
	return res, nil
}

func ToFloat64Slice(i any) ([]float64, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]float64); ok {
		return l, nil
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %#v of type %T to []float64", i, i)
	}
	num := v.Len()
	res := make([]float64, num)
	var err error
	for j := 0; j < num; j++ {
		res[j], err = ToFloat64(v.Index(j).Interface())
		if err != nil {
			return nil, fmt.Errorf("convert index %d: %w", i, err)
		}
	}
	return res, nil
}

const (
	// MaxInt represents maximum int
	MaxInt = 1<<(8*unsafe.Sizeof(int(0))-1) - 1
	// MinInt represents minimum int
	MinInt = -1 << (8*unsafe.Sizeof(int(0)) - 1)
	// MaxUint represents maximum uint
	MaxUint = 1<<(8*unsafe.Sizeof(uint(0))) - 1
)

// ToInt converts i to int
func ToInt(i any) (int, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int64", i, i)
	}
	if n > MaxInt || n < MinInt {
		return 0, strconv.ErrRange
	}
	return int(n), nil
}

// ToInt8 converts i to int8
func ToInt8(i any) (int8, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int8", i, i)
	}
	if n > math.MaxInt8 || n < math.MinInt8 {
		return 0, strconv.ErrRange
	}
	return int8(n), nil
}

// ToInt16 converts i to int16
func ToInt16(i any) (int16, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int16", i, i)
	}
	if n > math.MaxInt16 || n < math.MinInt16 {
		return 0, strconv.ErrRange
	}
	return int16(n), nil
}

func ToInt32(i any) (int32, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int32", i, i)
	}
	if n > math.MaxInt32 || n < math.MinInt32 {
		return 0, strconv.ErrRange
	}
	return int32(n), nil
}

func MustInt32(i any) int32 {
	v, err := ToInt32(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToInt64(i any) (int64, error) {
	return parseInt64(i)
}

func ToUint(i any) (uint, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint", i, i)
	}
	if n > MaxUint {
		return 0, strconv.ErrRange
	}
	return uint(n), nil
}

func ToUint8(i any) (uint8, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint8", i, i)
	}
	if n > math.MaxUint8 {
		return 0, strconv.ErrRange
	}
	return uint8(n), nil
}

func ToUint16(i any) (uint16, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint16", i, i)
	}
	if n > math.MaxUint16 {
		return 0, strconv.ErrRange
	}
	return uint16(n), nil
}

func ToUint32(i any) (uint32, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint32", i, i)
	}
	if n > math.MaxUint32 {
		return 0, strconv.ErrRange
	}
	return uint32(n), nil
}

func ToUint64(i any) (uint64, error) {
	return parseUint64(i)
}

func ToIntSlice(i any) ([]int, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]int); ok {
		return l, nil
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]int, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = ToInt(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert index %d: %w", j, err)
			}
		}
		return res, nil
	default:
		if k, err := ToInt(i); err == nil {
			return []int{k}, nil
		}
		return nil, fmt.Errorf("cannot convert %v to slice", v.Kind())
	}
}

func ToInt64Slice(i any) ([]int64, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]int64); ok {
		return l, nil
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]int64, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = parseInt64(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert element at index %d: %w", i, err)
			}
		}
		return res, nil
	default:
		if k, err := ToInt64(i); err == nil {
			return []int64{k}, nil
		}
		return nil, fmt.Errorf("cannot convert %#v of type %T to []int64", i, i)
	}
}

func ToUintSlice(i any) ([]uint, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]uint); ok {
		return l, nil
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]uint, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = ToUint(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert element at index %d: %w", i, err)
			}
		}
		return res, nil
	default:
		if ui, err := ToUint(i); err == nil {
			return []uint{ui}, nil
		}
		return nil, fmt.Errorf("cannot convert %#v of type %T to []uint", i, i)
	}
}

func ToUint64Slice(i any) ([]uint64, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]uint64); ok {
		return l, nil
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]uint64, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = parseUint64(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert element at index %d: %w", i, err)
			}
		}
		return res, nil
	default:
		if ui, err := ToUint64(i); err == nil {
			return []uint64{ui}, nil
		}
		return nil, fmt.Errorf("cannot convert %#v of type %T to []uint64", i, i)
	}
}

func parseInt64(i any) (int64, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return 0, errorx.NotExist
	}
	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	v := reflect.ValueOf(i)
	if rtx.IsInt(v) {
		return v.Int(), nil
	}

	if rtx.IsUint(v) {
		n := v.Uint()
		if n > math.MaxInt64 {
			return 0, strconv.ErrRange
		}
		return int64(n), nil
	}

	if rtx.IsFloat(v) {
		return int64(v.Float()), nil
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.String:
		n, err := strconv.ParseInt(v.String(), 0, 64)
		if err == nil {
			return n, nil
		}
		if errors.Is(err, strconv.ErrRange) {
			return 0, err
		}
		if f, fErr := strconv.ParseFloat(v.String(), 64); fErr == nil {
			return int64(f), nil
		}
		return 0, err
	default:
		return 0, strconv.ErrSyntax
	}
}

func parseUint64(i any) (uint64, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return 0, errorx.NotExist
	}
	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	v := reflect.ValueOf(i)
	if rtx.IsInt(v) {
		n := v.Int()
		if n < 0 {
			return 0, strconv.ErrRange
		}
		return uint64(n), nil
	}

	if rtx.IsUint(v) {
		return v.Uint(), nil
	}

	if rtx.IsFloat(v) {
		f := v.Float()
		if f < 0 {
			return 0, strconv.ErrRange
		}
		return uint64(f), nil
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.String:
		n, err := strconv.ParseInt(v.String(), 0, 64)
		if err == nil {
			if n < 0 {
				return 0, strconv.ErrRange
			}
			return uint64(n), nil
		}
		if errors.Is(err, strconv.ErrRange) {
			return 0, err
		}
		if f, fErr := strconv.ParseFloat(v.String(), 64); fErr == nil {
			if f < 0 {
				return 0, strconv.ErrRange
			}
			return uint64(f), nil
		}
		return 0, err
	default:
		return 0, strconv.ErrSyntax
	}
}

func ToUniqueIntSlice(a []int) []int {
	m := make(map[int]struct{}, len(a))
	l := make([]int, 0, len(a))
	for _, v := range a {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		l = append(l, v)
	}
	return l
}

func ToUniqueInt64Slice(a []int64) []int64 {
	m := make(map[int64]struct{}, len(a))
	l := make([]int64, 0, len(a))
	for _, v := range a {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		l = append(l, v)
	}
	return l
}

func ToBytes(i any) ([]byte, error) {
	i = rtx.Indirect(i)
	switch v := i.(type) {
	case []byte:
		return v, nil
	case nil:
		return nil, errorx.NotExist
	case string:
		return []byte(v), nil
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
		return v.Bytes(), nil
	}
	return nil, fmt.Errorf("cannot convert %#v of type %T to []byte", i, i)
}

func ToByteArray8[T []byte | string](v T) [8]byte {
	if len(v) > 8 {
		panic("cannot convert into [8]byte")
	}
	var a [8]byte
	copy(a[:], v[:])
	return a
}

func ToByteArray16[T []byte | string](v T) [16]byte {
	if len(v) > 16 {
		panic("cannot convert into [16]byte")
	}
	var a [16]byte
	copy(a[:], v[:])
	return a
}

func ToByteArray32[T []byte | string](v T) [32]byte {
	if len(v) > 32 {
		panic("cannot convert into [32]byte")
	}
	var a [32]byte
	copy(a[:], v[:])
	return a
}

func ToByteArray64[T []byte | string](v T) [64]byte {
	if len(v) > 64 {
		panic("cannot convert into [64]byte")
	}
	var a [64]byte
	copy(a[:], v[:])
	return a
}

// ToString converts i to string
// i can be string, integer types, bool, []byte or any types which implement fmt.Stringer
func ToString(i any) (string, error) {
	i = rtx.IndirectToStringerOrError(i)
	if i == nil {
		return "", errorx.NotExist
	}
	switch v := i.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case fmt.Stringer:
		return v.String(), nil
	case error:
		return v.Error(), nil
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Float32, reflect.Float64:
		return fmt.Sprint(v.Interface()), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes()), nil
		}
	}
	return "", fmt.Errorf("cannot convert %#v of type %T to string", i, i)
}

func ToStringSlice(i any) ([]string, error) {
	i = rtx.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]string); ok {
		return l, nil
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]string, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = ToString(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert element at index %d: %w", i, err)
			}
		}
		return res, nil
	default:
		if s, err := ToString(i); err == nil {
			return strings.Fields(s), nil
		}
		return nil, fmt.Errorf("cannot convert %#v of type %T to []string", i, i)
	}
}

func ToEmailAddress(s string) (string, error) {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

const (
	kilobyte = 1 << (10 * (1 + iota)) // 1 << (10*1)
	megabyte                          // 1 << (10*2)
	gigabyte                          // 1 << (10*3)
	terabyte                          // 1 << (10*4)
	petabyte                          // 1 << (10*5)
)

func SizeToHumanReadable(size int64) string {
	if size < kilobyte {
		return fmt.Sprintf("%d B", size)
	} else if size < megabyte {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(kilobyte))
	} else if size < gigabyte {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(megabyte))
	} else if size < terabyte {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(gigabyte))
	} else if size < petabyte {
		return fmt.Sprintf("%.2f TB", float64(size)/float64(terabyte))
	} else {
		return fmt.Sprintf("%.2f PB", float64(size)/float64(petabyte))
	}
}

// ToList creates list.List
// i can be nil, *list.List, or array/slice
func ToList(i any) *list.List {
	if i == nil {
		return list.New()
	}

	if l, ok := i.(*list.List); ok {
		return l
	}

	lt := reflect.TypeOf((*list.List)(nil))
	if it := reflect.TypeOf(i); it.ConvertibleTo(lt) {
		return reflect.ValueOf(i).Convert(lt).Interface().(*list.List)
	}

	l := list.New()
	v := reflect.ValueOf(rtx.Indirect(i))
	if v.IsValid() && (v.Kind() == reflect.Slice || v.Kind() == reflect.Array) && !v.IsNil() {
		for j := 0; j < v.Len(); j++ {
			l.PushBack(v.Index(j).Interface())
		}
	} else {
		l.PushBack(i)
	}
	return l
}

// SliceToList converts slice to list.List
func SliceToList[E any](a []E) *list.List {
	l := list.New()
	for _, e := range a {
		l.PushBack(e)
	}
	return l
}
