package nomobile

import (
	"encoding/json"
	"log"
	"reflect"
)

type Equaler[T any] interface {
	Equals(T) bool
}

type List[E comparable] struct {
	elements []E
}

func NewList[E comparable](l []E) *List[E] {
	return &List[E]{elements: l}
}

func (l *List[E]) Elements() []E {
	return l.elements
}

func (l *List[E]) Len() int {
	return len(l.elements)
}

func (l *List[E]) Get(index int) E {
	return l.elements[index]
}

func (l *List[E]) Add(e E) {
	l.elements = append(l.elements, e)
}

func (l *List[E]) AddList(v *List[E]) {
	l.elements = append(l.elements, v.elements...)
}

func (l *List[E]) RemoveAt(i int) {
	l.elements = append(l.elements[:i], l.elements[i+1:]...)
}

func (l *List[E]) Remove(e E) {
	if i := l.IndexOf(e); i >= 0 {
		l.elements = append(l.elements[:i], l.elements[i+1:]...)
	}
}

func (l *List[E]) Insert(i int, e E) {
	if len(l.elements) <= i {
		l.elements = append(l.elements, e)
	} else {
		n := len(l.elements)
		l.elements = append(l.elements, e)
		copy(l.elements[i+1:], l.elements[i:n])
		l.elements[i] = e
	}
}

func (l *List[E]) IndexOf(e E) int {
	for i, v := range l.elements {
		if v == e {
			return i
		}

		if eq, ok := any(v).(Equaler[E]); ok && eq.Equals(e) {
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
	return len(l.elements) == 0
}

func (l *List[E]) Clear() {
	l.elements = l.elements[0:0]
}

func (l *List[E]) Reverse() {
	for i, j := 0, l.Len()-1; i < j; i, j = i+1, j-1 {
		l.elements[i], l.elements[j] = l.elements[j], l.elements[i]
	}
}

func (l *List[E]) Clone() *List[E] {
	res := new(List[E])
	res.elements = make([]E, l.Len())
	copy(res.elements, l.elements)
	return res
}

func (l *List[E]) JSONString() string {
	data, err := json.Marshal(l.elements)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(data)
}

func (l *List[E]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.elements)
}

func (l *List[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.elements)
}
