package dependencyinjection

import (
	"errors"
	"strings"
	"sync"
)

// Storage abstracts key–value persistence.
type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

// MapStore is an in-memory implementation of Storage.
type MapStore struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewMapStore creates a ready-to-use MapStore.
func NewMapStore() *MapStore {
	return &MapStore{data: make(map[string]string)}
}

// Get returns the value for key. Returns ErrNotFound when key is absent.
func (m *MapStore) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return val, nil
}

// Set stores value under key.
func (m *MapStore) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

// Delete removes key. Returns ErrNotFound when key is absent.
func (m *MapStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[key]; !ok {
		return ErrNotFound
	}
	delete(m.data, key)
	return nil
}

// ErrNotFound is returned when a key does not exist.
var ErrNotFound = errors.New("not found")

// Service uses a Storage backend for key–value operations.
type Service struct {
	storage Storage
}

// NewService returns a Service backed by s. Returns nil when s is nil.
func NewService(s Storage) *Service {
	if s == nil {
		return nil
	}
	return &Service{storage: s}
}

// Process reads the current value for key, uppercases it, and writes it back.
func (s *Service) Process(key, val string) error {
	if err := s.storage.Set(key, strings.ToUpper(val)); err != nil {
		return err
	}
	return nil
}

// Get delegates to the underlying storage.
func (s *Service) Get(key string) (string, error) {
	return s.storage.Get(key)
}

// Delete delegates to the underlying storage.
func (s *Service) Delete(key string) error {
	return s.storage.Delete(key)
}
