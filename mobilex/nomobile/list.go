package nomobile

import (
	"encoding/json"
	"reflect"
)

type equaler[T any] interface {
	Equals(T) bool
}

type List[E comparable] struct {
	Elements []E
}

func (l *List[E]) Len() int {
	return len(l.Elements)
}

func (l *List[E]) Get(index int) E {
	return l.Elements[index]
}

func (l *List[E]) Add(e E) {
	l.Elements = append(l.Elements, e)
}

func (l *List[E]) RemoveAt(i int) {
	l.Elements = append(l.Elements[:i], l.Elements[i+1:]...)
}

func (l *List[E]) Remove(e E) {
	if i := l.IndexOf(e); i >= 0 {
		l.Elements = append(l.Elements[:i], l.Elements[i+1:]...)
	}
}

func (l *List[E]) Insert(i int, e E) {
	if len(l.Elements) <= i {
		l.Elements = append(l.Elements, e)
	} else {
		n := len(l.Elements)
		l.Elements = append(l.Elements, e)
		copy(l.Elements[i+1:], l.Elements[i:n])
		l.Elements[i] = e
	}
}

func (l *List[E]) IndexOf(e E) int {
	for i, v := range l.Elements {
		if v == e {
			return i
		}

		if eq, ok := any(v).(equaler[E]); ok && eq.Equals(e) {
			return i
		}
	}
	return -1
}

func (l *List[E]) First() E {
	var e E
	if l.Len() == 0 {
		if reflect.TypeOf(e).Kind() != reflect.Pointer {
			panic("list is empty")
		}
		return e
	}
	return l.Get(0)
}

func (l *List[E]) Last() E {
	var e E
	if l.Len() == 0 {
		if reflect.TypeOf(e).Kind() != reflect.Pointer {
			panic("list is empty")
		}
		return e
	}
	return l.Get(l.Len() - 1)
}

func (l *List[E]) IsEmpty() bool {
	return len(l.Elements) == 0
}

func (l *List[E]) Clear() {
	l.Elements = l.Elements[0:0]
}

func (l *List[E]) Reverse() {
	for i, j := 0, l.Len()-1; i < j; i, j = i+1, j-1 {
		l.Elements[i], l.Elements[j] = l.Elements[j], l.Elements[i]
	}
}

func (l *List[E]) Clone() *List[E] {
	res := new(List[E])
	res.Elements = make([]E, l.Len())
	copy(res.Elements, l.Elements)
	return res
}

func (l *List[E]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.Elements)
}

func (l *List[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Elements)
}
