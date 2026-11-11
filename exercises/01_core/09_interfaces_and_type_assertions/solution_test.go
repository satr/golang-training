package interfaces_and_type_assertions

import (
	"testing"
)

func TestPrintValue(t *testing.T) {
	// Test string
	// We can't easily test the output, but we can check that it doesn't panic
	PrintValue("hello")
	PrintValue(42)
	PrintValue(3.14)
}

func TestProcessValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string input",
			input:    "hello",
			expected: "string: hello",
		},
		{
			name:     "int input",
			input:    42,
			expected: "int: 42",
		},
		{
			name:     "bool true input",
			input:    true,
			expected: "bool: true",
		},
		{
			name:     "bool false input",
			input:    false,
			expected: "bool: false",
		},
		{
			name:     "float input",
			input:    3.14,
			expected: "unknown type",
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "unknown type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCanPrint(t *testing.T) {
	// Test with a type that doesn't implement Stringer (string)
	s := "hello"
	if CanPrint(s) {
		t.Error("string should not be able to print (string does not implement Stringer)")
	}

	// Test with a type that implements Stringer
	myStr := MyString{value: "test"}
	if !CanPrint(myStr) {
		t.Error("MyString should be able to print")
	}

	// Test with a type that doesn't implement Stringer
	i := 42
	if CanPrint(i) {
		t.Error("int should not be able to print")
	}
}
