package maps_and_deterministic_ordering

import "sort"

// Index builds a map-backed index from a slice of strings.
// The returned slice should be deterministically ordered (sorted by key).
func Index(items []string) []KeyValue {
	// Build the map with counts
	counts := make(map[string]int)
	for _, item := range items {
		counts[item]++
	}

	// Create a slice of keys to sort them
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build the result slice with sorted keys
	result := make([]KeyValue, 0, len(keys))
	for _, key := range keys {
		result = append(result, KeyValue{Key: key, Value: counts[key]})
	}

	return result
}

// KeyValue represents a key-value pair.
type KeyValue struct {
	Key   string
	Value int
}
