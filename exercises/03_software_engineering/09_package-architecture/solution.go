package packagearchitecture

import (
	"fmt"
	"sync"
)

// Store abstracts persistent key-value access.
type Store interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewInMemoryStore() Store {
	return &InMemoryStore{data: make(map[string]string)}
}

func (s *InMemoryStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}
	return v, nil
}

func (s *InMemoryStore) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *InMemoryStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; !ok {
		return fmt.Errorf("key %q not found", key)
	}
	delete(s.data, key)
	return nil
}
