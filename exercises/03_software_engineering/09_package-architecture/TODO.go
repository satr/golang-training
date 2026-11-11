//go:build todo

package packagearchitecture

// Store abstracts persistent key-value access.
type Store interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

// InMemoryStore is an in-memory Store implementation.
type InMemoryStore struct{}

func NewInMemoryStore() Store                           { return nil }
func (s *InMemoryStore) Get(key string) (string, error) { return "", nil }
func (s *InMemoryStore) Set(key, value string) error    { return nil }
func (s *InMemoryStore) Delete(key string) error        { return nil }
