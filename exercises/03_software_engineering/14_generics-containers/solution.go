package genericscontainers

// Stack is a generic last-in-first-out container.
type Stack[T any] struct {
	items []T
}

// NewStack creates a new empty Stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item.
func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, StackEmpty
	}
	n := len(s.items)
	item := s.items[n-1]
	s.items = s.items[:n-1]
	return item, nil
}

// Peek returns the top item without removing it.
func (s *Stack[T]) Peek() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, StackEmpty
	}
	return s.items[len(s.items)-1], nil
}

// Size returns the number of items in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// IsEmpty reports whether the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// StackEmpty is returned when Pop or Peek is called on an empty stack.
var StackEmpty = error(emptyErr{})

type emptyErr struct{}

func (e emptyErr) Error() string { return "stack is empty" }

// Map is a generic wrapper around map[K]V with nil-key protection.
type Map[K comparable, V any] struct {
	data map[K]V
}

// NewMap creates a new empty Map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{data: make(map[K]V)}
}

// Set stores a value at the given key.
func (m *Map[K, V]) Set(key K, value V) {
	m.data[key] = value
}

// Get returns the value for the given key.
func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.data[key]
	return v, ok
}

// Len returns the number of entries in the map.
func (m *Map[K, V]) Len() int {
	return len(m.data)
}
