// Package memory provides in-memory repository implementations.
package memory

import "sync"

// SafeMap is a thread-safe map implementation protected by a mutex.
type SafeMap struct {
	mu   sync.Mutex
	data map[int]int
}

// NewSafeMap creates a new SafeMap instance.
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[int]int),
	}
}

// Update atomically reads the value by key (creates it if missing) and writes the result of fn.
// Returns true if the key was created.
func (m *SafeMap) Update(key int, fn func(val int) int) (created bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		m.data[key] = 0
		created = true
	}

	m.data[key] = fn(m.data[key])
	return
}

// GetAll returns a copy of all data.
func (m *SafeMap) GetAll() map[int]int {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make(map[int]int, len(m.data))
	for k, v := range m.data {
		result[k] = v
	}

	return result
}
