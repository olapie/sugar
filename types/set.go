package types

import (
	"encoding/json"
)

type void struct{}

type Set[K comparable] struct {
	entries map[K]void
}

func NewSet[K comparable](capacity int) *Set[K] {
	s := &Set[K]{}
	s.entries = make(map[K]void, capacity)
	return s
}

func (s *Set[K]) Add(item K) {
	s.entries[item] = void{}
}

func (s *Set[K]) Contains(item K) bool {
	_, found := s.entries[item]
	return found
}

func (s *Set[K]) Remove(item K) {
	delete(s.entries, item)
}

func (s *Set[K]) Slice() []K {
	l := make([]K, 0, len(s.entries))
	for k := range s.entries {
		l = append(l, k)
	}
	return l
}

func (s *Set[K]) Range(f func(v K) bool) {
	for k := range s.entries {
		if !f(k) {
			break
		}
	}
}

func (s *Set[K]) Len() int {
	return len(s.entries)
}

var (
	_ json.Unmarshaler = (*Set[int64])(nil)
	_ json.Marshaler   = (*Set[int64])(nil)
)

func (s *Set[K]) UnmarshalJSON(data []byte) error {
	var a []K
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	s.entries = make(map[K]void)
	for _, v := range a {
		s.Add(v)
	}
	return nil
}

func (s *Set[K]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Slice())
}
