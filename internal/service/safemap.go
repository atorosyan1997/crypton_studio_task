// Package service contains business logic.
package service

import (
	"sync/atomic"

	"crypton_studio_task/internal/model"
)

// SafeMapStorage is an interface for a thread-safe map with integer keys.
type SafeMapStorage interface {
	// Update atomically reads the value by key (creates it if missing) and writes the result of fn.
	// Returns true if the key was created.
	Update(key int, fn func(val int) int) (created bool)
	// GetAll returns a copy of all data.
	GetAll() map[int]int
}

// SafeMap is a service for working with a thread-safe map.
type SafeMap struct {
	storage     SafeMapStorage
	accessCount atomic.Int64
	insertCount atomic.Int64
}

// NewSafeMap creates a new service instance.
func NewSafeMap(storage SafeMapStorage) *SafeMap {
	return &SafeMap{
		storage: storage,
	}
}

// GetAndIncrement retrieves the value by key (creates it if missing) and increments it by 1.
func (s *SafeMap) GetAndIncrement(key int) {
	s.accessCount.Add(1)

	created := s.storage.Update(key, s.increment)
	if created {
		s.insertCount.Add(1)
	}
}

// increment adds 1 to the value.
func (s *SafeMap) increment(val int) int {
	return val + 1
}

// Stats returns access and insert counters.
func (s *SafeMap) Stats() model.SafeMapStats {
	return model.SafeMapStats{
		AccessCount: s.accessCount.Load(),
		InsertCount: s.insertCount.Load(),
	}
}
