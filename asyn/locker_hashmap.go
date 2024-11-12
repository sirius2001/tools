package asyn

import (
	"fmt"
	"sync"
)

type Map[K comparable, V any] struct {
	m     map[K]V
	mlock sync.Mutex
}

// NewMap instantiates a hash map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{m: make(map[K]V)}
}

// Put inserts element into the map.
func (m *Map[K, V]) Put(key K, value V) {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	m.m[key] = value
}

// Get searches the element in the map by key and returns its value or nil if key is not found in map.
// Second return parameter is true if key was found, otherwise false.
func (m *Map[K, V]) Get(key K) (value V, found bool) {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	value, found = m.m[key]
	return
}

// Remove removes the element from the map by key.
func (m *Map[K, V]) Remove(key K) {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	delete(m.m, key)
}

// Empty returns true if map does not contain any elements
func (m *Map[K, V]) Empty() bool {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	return m.Size() == 0
}

// Size returns number of elements in the map.
func (m *Map[K, V]) Size() int {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	return len(m.m)
}

// Keys returns all keys (random order).
func (m *Map[K, V]) Keys() []K {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	keys := make([]K, m.Size())
	count := 0
	for key := range m.m {
		keys[count] = key
		count++
	}
	return keys
}

// Values returns all values (random order).
func (m *Map[K, V]) Values() []V {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	values := make([]V, m.Size())
	count := 0
	for _, value := range m.m {
		values[count] = value
		count++
	}
	return values
}

// Clear removes all elements from the map.
func (m *Map[K, V]) Clear() {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	clear(m.m)
}

// String returns a string representation of container
func (m *Map[K, V]) String() string {
	m.mlock.Lock()
	defer m.mlock.Unlock()
	str := "HashMap\n"
	str += fmt.Sprintf("%v", m.m)
	return str
}
