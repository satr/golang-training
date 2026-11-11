package racedetection

import (
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrEmptyKey reports an empty cache key.
	ErrEmptyKey = errors.New("cache key must not be empty")
	// ErrNilLoader reports a missing Loader function.
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
// the same key. The zero value is ready for use. A Cache must not be copied
// after first use.
type Cache struct {
	mu      sync.Mutex
	entries map[string]*entry
}

// Get returns the cached value for key. On a cache miss, one caller runs load;
// concurrent callers for the same key wait for and share that result.
// Failed loads are returned to current waiters but are not cached.
func (c *Cache) Get(key string, load Loader) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if load == nil {
		return "", ErrNilLoader
	}

	c.mu.Lock()
	if c.entries == nil {
		c.entries = make(map[string]*entry)
	}
	if existing, ok := c.entries[key]; ok {
		c.mu.Unlock()
		<-existing.ready
		return existing.value, existing.err
	}

	current := &entry{ready: make(chan struct{})}
	c.entries[key] = current
	c.mu.Unlock()

	value, err := load(key)
	if err != nil {
		value = ""
		err = fmt.Errorf("load %q: %w", key, err)
	}

	c.mu.Lock()
	current.value = value
	current.err = err
	if err != nil {
		delete(c.entries, key)
	}
	close(current.ready)
	c.mu.Unlock()

	return value, err
}

// Delete removes a successfully cached value. It reports whether a completed
// value was removed. An in-flight load is not interrupted or removed.
func (c *Cache) Delete(key string) bool {
	if key == "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	current, ok := c.entries[key]
	if !ok {
		return false
	}
	select {
	case <-current.ready:
		delete(c.entries, key)
		return current.err == nil
	default:
		return false
	}
}

// Len returns the number of successfully cached values. In-flight loads are
// not counted.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for _, current := range c.entries {
		select {
		case <-current.ready:
			if current.err == nil {
				count++
			}
		default:
		}
	}
	return count
}
