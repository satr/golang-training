//go:build todo

package dependencyinjection

import (
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
	// TODO: Initialize and return a MapStore with a pre-allocated map.
	return &MapStore{}
}

// Get returns the value for key. Returns an error when key is absent.
func (m *MapStore) Get(key string) (string, error) {
	// TODO: Implement thread-safe lookup. Return ErrNotFound when key is absent.
	return "", nil
}

// Set stores value under key.
func (m *MapStore) Set(key, value string) error {
	// TODO: Implement thread-safe write.
	return nil
}

// Delete removes key. Returns an error when key is absent.
func (m *MapStore) Delete(key string) error {
	// TODO: Implement thread-safe delete. Return ErrNotFound when key is absent.
	return nil
}

// ErrNotFound is returned when a key does not exist.
var ErrNotFound error = nil // TODO: define a real sentinel error

// Service uses a Storage backend for key–value operations.
type Service struct {
	storage Storage
}

// NewService returns a Service backed by s. Returns nil when s is nil.
func NewService(s Storage) *Service {
	// TODO: Reject nil storage. Return a Service holding s.
	return nil
}

// Process reads the current value for key, uppercases it, and writes it back.
func (s *Service) Process(key, val string) error {
	// TODO: Set the uppercased value for key using s.storage.
	return nil
}

// Get delegates to the underlying storage.
func (s *Service) Get(key string) (string, error) {
	// TODO: Delegate to s.storage.Get.
	return "", nil
}

// Delete delegates to the underlying storage.
func (s *Service) Delete(key string) error {
	// TODO: Delegate to s.storage.Delete.
	return nil
}
