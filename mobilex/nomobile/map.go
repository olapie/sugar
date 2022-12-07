package nomobile

import (
	"encoding/json"
	"fmt"

	"code.olapie.com/sugar/mapping"
)

type Map[K comparable, V any] struct {
	m map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Get(k K) V {
	return m.m[k]
}

func (m *Map[K, V]) Contains(k K) bool {
	_, ok := m.m[k]
	return ok
}

func (m *Map[K, V]) Set(k K, v V) {
	m.m[k] = v
}

func (m *Map[K, V]) Remove(k K) {
	delete(m.m, k)
}

func (m *Map[K, V]) Count() int {
	return len(m.m)
}

////////////////////////////////
// Methods below are not supported by gomobile

func (m *Map[K, V]) Keys() *List[K] {
	return &List[K]{
		elements: mapping.GetKeys(m.m),
	}
}

func (m *Map[K, V]) Clone() *Map[K, V] {
	return &Map[K, V]{
		m: mapping.Clone(m.m),
	}
}

func (m *Map[K, V]) InsertMap(m2 *Map[K, V]) {
	for k, v := range m2.m {
		m.m[k] = v
	}
}

func (m *Map[K, V]) JSONString() string {
	data, err := json.Marshal(m.m)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.m)
}

func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.m)
}
