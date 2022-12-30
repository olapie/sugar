package xmap

import "sync"

// Map is a thread safe map
// must not be copied after first use
type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Get(key K) (actual V, got bool) {
	v, got := m.m.Load(key)
	if !got {
		return
	}
	actual, got = v.(V)
	return
}

func (m *Map[K, V]) GetOrStore(key K, value V) (actual V, got bool) {
	v, got := m.m.LoadOrStore(key, value)
	if !got {
		return actual, false
	}
	actual, got = v.(V)
	if got {
		return actual, true
	}
	m.m.Store(key, value)
	return value, true
}

func (m *Map[K, V]) Set(k K, v V) {
	m.m.Store(k, v)
}

func (m *Map[K, V]) Delete(k K) {
	m.m.Delete(k)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
