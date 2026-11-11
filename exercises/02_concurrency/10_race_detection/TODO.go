//go:build todo

package racedetection

import (
	"errors"
	"sync"
)

var (
	ErrEmptyKey  = errors.New("cache key must not be empty")
	ErrNilLoader = errors.New("loader must not be nil")
)

// Loader computes a value for a cache key.
type Loader func(key string) (string, error)

type entry struct {
	ready chan struct{}
	value string
	err   error
}

// Cache memoizes successfully loaded values and coalesces concurrent loads for
// the same key. Its zero value is ready for use.
type Cache struct {
	mu      sync.Mutex
	entries map[string]*entry
}

// Get returns a cached value or loads it once for concurrent callers.
func (c *Cache) Get(key string, load Loader) (string, error) {
	panic("TODO: implement Get")
}

// Delete removes a completed cached value without interrupting an in-flight
// load.
func (c *Cache) Delete(key string) bool {
	panic("TODO: implement Delete")
}

// Len returns the number of completed successful entries.
func (c *Cache) Len() int {
	panic("TODO: implement Len")
}
