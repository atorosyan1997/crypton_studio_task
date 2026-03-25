// Package service contains business logic.
package service

import (
	"math/rand"
	"sync"
	"testing"

	"crypton_studio_task/internal/repository/memory"
)

const (
	// year = 1799 — birth year of Alexander Sergeyevich Pushkin.
	year = 1799
	// goroutines — number of goroutines as required by the task.
	goroutines = 4
	// expectedValue — expected value for each key after all goroutines finish.
	expectedValue = 3
)

func TestSafeMapConcurrentIncrement(t *testing.T) {
	storage := memory.NewSafeMap()
	svc := NewSafeMap(storage)

	// Build a list of keys: each key from 1 to year appears exactly expectedValue times.
	keys := make([]int, 0, year*expectedValue)
	for i := 1; i <= year; i++ {
		for j := 0; j < expectedValue; j++ {
			keys = append(keys, i)
		}
	}

	// Shuffle so that goroutines do not access keys sequentially.
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// Split into chunks for goroutines.
	chunks := splitIntoChunks(keys, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(chunk []int) {
			defer wg.Done()
			for _, key := range chunk {
				svc.GetAndIncrement(key)
			}
		}(chunks[i])
	}

	wg.Wait()

	// Verify counters before reading values.
	stats := svc.Stats()

	// Access count = year * expectedValue.
	var expectedAccess int64 = year * expectedValue
	if stats.AccessCount != expectedAccess {
		t.Errorf("access count: expected %d, got %d", expectedAccess, stats.AccessCount)
	}

	// Insert count = year (each key is created once on first access).
	var expectedInsert int64 = year
	if stats.InsertCount != expectedInsert {
		t.Errorf("insert count: expected %d, got %d", expectedInsert, stats.InsertCount)
	}

	// Verify that each key holds expectedValue (via GetAll to avoid affecting counters).
	data := storage.GetAll()
	for key := 1; key <= year; key++ {
		val, ok := data[key]
		if !ok {
			t.Errorf("key %d: missing from map", key)
			continue
		}
		if val != expectedValue {
			t.Errorf("key %d: expected %d, got %d", key, expectedValue, val)
		}
	}
}

// BenchmarkGetAndIncrement_SingleGoroutine measures single-threaded performance.
func BenchmarkGetAndIncrement_SingleGoroutine(b *testing.B) {
	storage := memory.NewSafeMap()
	svc := NewSafeMap(storage)

	for i := 0; i < b.N; i++ {
		svc.GetAndIncrement(i % year)
	}
}

// BenchmarkGetAndIncrement_Concurrent measures concurrent performance with 4 goroutines.
func BenchmarkGetAndIncrement_Concurrent(b *testing.B) {
	storage := memory.NewSafeMap()
	svc := NewSafeMap(storage)

	b.RunParallel(func(pb *testing.PB) {
		key := 0
		for pb.Next() {
			svc.GetAndIncrement(key % year)
			key++
		}
	})
}

// splitIntoChunks splits a slice into n roughly equal parts.
func splitIntoChunks(items []int, n int) [][]int {
	chunks := make([][]int, n)
	chunkSize := len(items) / n
	remainder := len(items) % n

	offset := 0
	for i := 0; i < n; i++ {
		size := chunkSize
		if i < remainder {
			size++
		}
		chunks[i] = items[offset : offset+size]
		offset += size
	}

	return chunks
}
