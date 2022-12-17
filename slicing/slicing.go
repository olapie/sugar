package slicing

import (
	"fmt"
	"reflect"

	"code.olapie.com/sugar/rtx"
)

func Clone[T any](a []T) []T {
	res := make([]T, len(a))
	copy(res, a)
	return res
}

func Transform[A any, B any](a []A, f func(A) (B, error)) ([]B, error) {
	b := make([]B, len(a))
	var err error
	for i := range a {
		b[i], err = f(a[i])
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
	}
	return b, nil
}

func ToSet[A any, B comparable](a []A, f func(A) (B, error)) (map[B]bool, error) {
	m := make(map[B]bool, len(a))
	for i, v := range a {
		if f == nil {
			m[any(a).(B)] = true
		} else {
			b, err := f(v)
			if err != nil {
				return nil, fmt.Errorf("index %d: %w", i, err)
			}
			m[b] = true
		}
	}
	return m, nil
}

func MustTransform[A any, B any](a []A, f func(A) B) []B {
	b := make([]B, len(a))
	for i := range a {
		b[i] = f(a[i])
	}
	return b
}

func MustToSet[A any, B comparable](a []A, f func(A) B) map[B]bool {
	m := make(map[B]bool, len(a))
	for _, v := range a {
		if f == nil {
			m[any(a).(B)] = true
		} else {
			m[f(v)] = true
		}
	}
	return m
}

func Unique[E comparable](a []E) []E {
	m := make(map[E]struct{}, len(a))
	l := make([]E, 0, len(a))
	for _, v := range a {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		l = append(l, v)
	}
	return l
}

func Reverse[E comparable](a []E) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func ReverseArray(a any) bool {
	a = rtx.Indirect(a)
	if a == nil {
		return false
	}
	v := reflect.ValueOf(a)
	if v.IsNil() || !v.IsValid() || v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return false
	}

	for i, j := 0, v.Len()-1; i < j; i, j = i+1, j-1 {
		vi, vj := v.Index(i), v.Index(j)
		tmp := vi.Interface()
		if !vi.CanSet() {
			return false
		}
		vi.Set(vj)
		vj.Set(reflect.ValueOf(tmp))
	}
	return true
}

func Remove[E comparable](a []E, v E) []E {
	for i, e := range a {
		if e == v {
			a = append(a[:i], a[i+1:]...)
			break
		}
	}
	return a
}

func Contains[E comparable](a []E, v E) bool {
	for _, e := range a {
		if e == v {
			return true
		}
	}
	return false
}

func IndexOf[E comparable](a []E, v E) int {
	for i, e := range a {
		if e == v {
			return i
		}
	}
	return -1
}

func Filter[E comparable](a []E, filter func(e E) bool) []E {
	res := make([]E, 0, len(a)/2)
	for _, v := range a {
		if filter(v) {
			res = append(res, v)
		}
	}
	return res
}
