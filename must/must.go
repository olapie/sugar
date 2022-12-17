package must

import (
	"fmt"
	"log"

	"code.olapie.com/sugar/conv"
	"code.olapie.com/sugar/rtx"
)

// Get eliminates nil err and panics if err isn't nil
func Get[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// GetTwo eliminates nil err and panics if err isn't nil
func GetTwo[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// True panics if b is not true
func True(b bool, msgAndArgs ...any) {
	if !b {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// False panics if b is not true
func False(b bool, msgAndArgs ...any) {
	if b {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// Error panics if b is not nil
func Error(err error, msgAndArgs ...any) {
	if err == nil {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// NoError panics if b is not nil
func NoError(err error, msgAndArgs ...any) {
	if err != nil {
		if len(msgAndArgs) == 0 {
			rtx.PanicWithMessages(err)
		} else {
			msgAndArgs[0] = err.Error() + " " + fmt.Sprint(msgAndArgs[0])
			rtx.PanicWithMessages(msgAndArgs...)
		}
	}
}

// Nil panics if v is not nil
func Nil[T any](v *T, msgAndArgs ...any) {
	if v != nil {
		if len(msgAndArgs) == 0 {
			rtx.PanicWithMessages(fmt.Sprintf("%#v", v))
		}
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// NotNil panics if v is nil
func NotNil[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func EmptySlice[T any](a []T, msgAndArgs ...any) {
	if len(a) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func NotEmptySlice[T any](a []T, msgAndArgs ...any) {
	if len(a) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func EmptyString[S ~string](s S, msgAndArgs ...any) {
	if len(s) != 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func NotEmptyString[S ~string](s S, msgAndArgs ...any) {
	if len(s) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func EmptyMap[K comparable, V any](m map[K]V, msgAndArgs ...any) {
	if len(m) != 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

func NotEmptyMap[K comparable, V any](m map[K]V, msgAndArgs ...any) {
	if len(m) == 0 {
		rtx.PanicWithMessages(msgAndArgs...)
	}
}

// ToBoolSlice converts i to []bool, will panic if failed
func ToBoolSlice(i any) []bool {
	v, err := conv.ToBoolSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

// ToBool converts i to bool, will panic if failed
func ToBool(i any) bool {
	v, err := conv.ToBool(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToFloat32(i any) float32 {
	v, err := conv.ToFloat32(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToFloat64(i any) float64 {
	v, err := conv.ToFloat64(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToFloat32Slice(i any) []float32 {
	v, err := conv.ToFloat32Slice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToFloat64Slice(i any) []float64 {
	v, err := conv.ToFloat64Slice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

// ToInt panics if ToInt(i) failed
func ToInt(i any) int {
	v, err := conv.ToInt(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

// ToInt8 panics if ToInt8(i) failed
func ToInt8(i any) int8 {
	v, err := conv.ToInt8(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

// ToInt16 panics if ToInt16(i) failed
func ToInt16(i any) int16 {
	v, err := conv.ToInt16(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToInt64(i any) int64 {
	v, err := conv.ToInt64(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
func ToUint(i any) uint {
	v, err := conv.ToUint(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
func ToUint8(i any) uint8 {
	v, err := conv.ToUint8(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToUint16(i any) uint16 {
	v, err := conv.ToUint16(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToUint32(i any) uint32 {
	v, err := conv.ToUint32(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToUint64(i any) uint64 {
	v, err := conv.ToUint64(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToIntSlice(i any) []int {
	v, err := conv.ToIntSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToInt64Slice(i any) []int64 {
	v, err := conv.ToInt64Slice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToUintSlice(i any) []uint {
	v, err := conv.ToUintSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
func ToUint64Slice(i any) []uint64 {
	v, err := conv.ToUint64Slice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToString(i any) string {
	v, err := conv.ToString(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
func ToStringSlice(i any) []string {
	v, err := conv.ToStringSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
