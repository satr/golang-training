//go:build todo

package genericscontainers

// TODO: Implement generic Stack[T] and Map[K, V] containers.
type Stack[T any] struct{}
type Map[K comparable, V any] struct{}

func NewStack[T any]() *Stack[T]              { return nil }
func (s *Stack[T]) Push(item T)               {}
func (s *Stack[T]) Pop() (T, error)           { var zero T; return zero, nil }
func (s *Stack[T]) Peek() (T, error)          { var zero T; return zero, nil }
func (s *Stack[T]) Size() int                 { return 0 }
func (s *Stack[T]) IsEmpty() bool             { return true }
func NewMap[K comparable, V any]() *Map[K, V] { return nil }
func (m *Map[K, V]) Set(key K, value V)       {}
func (m *Map[K, V]) Get(key K) (V, bool)      { var zero V; return zero, false }
func (m *Map[K, V]) Len() int                 { return 0 }
