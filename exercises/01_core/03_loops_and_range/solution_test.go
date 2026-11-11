package loopsandrange

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{name: "nil", values: nil, want: 0},
		{name: "empty", values: []int{}, want: 0},
		{name: "mixed signs", values: []int{5, -2, 8, -11}, want: 0},
		{name: "positive", values: []int{2, 4, 8}, want: 14},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Sum(test.values); got != test.want {
				t.Fatalf("Sum(%v) = %d, want %d", test.values, got, test.want)
			}
		})
	}
}

func TestRuneCounts(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[rune]int
	}{
		{name: "empty", text: "", want: map[rune]int{}},
		{name: "ASCII", text: "gopher", want: map[rune]int{'g': 1, 'o': 1, 'p': 1, 'h': 1, 'e': 1, 'r': 1}},
		{name: "repeated Unicode", text: "Go語語🙂", want: map[rune]int{'G': 1, 'o': 1, '語': 2, '🙂': 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := RuneCounts(test.text)
			if got == nil {
				t.Fatal("RuneCounts() returned nil map")
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("RuneCounts(%q) = %#v, want %#v", test.text, got, test.want)
			}
		})
	}
}
