//go:build todo

package maps_and_deterministic_ordering

// Index builds a map-backed index from a slice of strings.
// The returned slice should be deterministically ordered (sorted by key).
// TODO: Implement this function to build a deterministic index.
func Index(items []string) []KeyValue {
	// TODO: Implement the function body
	return nil
}

// KeyValue represents a key-value pair.
type KeyValue struct {
	Key   string
	Value int
}
