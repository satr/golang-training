package maps_and_deterministic_ordering

import (
	"reflect"
	"testing"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []KeyValue
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []KeyValue{},
		},
		{
			name:  "single item",
			input: []string{"apple"},
			expected: []KeyValue{
				{Key: "apple", Value: 1},
			},
		},
		{
			name:  "multiple items with duplicates",
			input: []string{"apple", "banana", "apple", "cherry", "banana", "apple"},
			expected: []KeyValue{
				{Key: "apple", Value: 3},
				{Key: "banana", Value: 2},
				{Key: "cherry", Value: 1},
			},
		},
		{
			name:  "items in reverse order",
			input: []string{"zebra", "yellow", "xray", "white", "yellow"},
			expected: []KeyValue{
				{Key: "white", Value: 1},
				{Key: "xray", Value: 1},
				{Key: "yellow", Value: 2},
				{Key: "zebra", Value: 1},
			},
		},
		{
			name:  "same items in different order",
			input: []string{"c", "b", "a", "c", "b", "a"},
			expected: []KeyValue{
				{Key: "a", Value: 2},
				{Key: "b", Value: 2},
				{Key: "c", Value: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Index(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Index(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIndexDeterministicOrdering tests that the output is always in the same order
func TestIndexDeterministicOrdering(t *testing.T) {
	input := []string{"zebra", "apple", "banana", "apple", "zebra", "cherry"}

	// Run multiple times to ensure consistent ordering
	var results [][]KeyValue
	for i := 0; i < 5; i++ {
		result := Index(input)
		results = append(results, result)
	}

	// Check that all results are identical
	for i := 1; i < len(results); i++ {
		if !reflect.DeepEqual(results[0], results[i]) {
			t.Errorf("Index ordering is not deterministic: %v vs %v", results[0], results[i])
		}
	}
}
