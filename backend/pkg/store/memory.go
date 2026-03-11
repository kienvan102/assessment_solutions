package store

import (
	"sync"
)

// Store is a generic thread-safe in-memory key-value store.
// It can be extended or replaced by a real database implementation.
type Store[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// New creates a new generic in-memory store.
func New[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		data: make(map[K]V),
	}
}

// Get retrieves a value by key. Returns the value and a boolean indicating if it was found.
func (s *Store[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Set stores a value by key.
func (s *Store[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Update modifies an existing value using an update function in a thread-safe manner.
func (s *Store[K, V]) Update(key K, updateFn func(V) V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if !ok {
		return false
	}
	s.data[key] = updateFn(val)
	return true
}

// Delete removes a value by key.
func (s *Store[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Len returns the number of items in the store.
func (s *Store[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
