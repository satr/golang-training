package slicesandarrays

import (
	"reflect"
	"testing"
)

func TestReverseFive(t *testing.T) {
	tests := []struct {
		name   string
		values [5]int
		want   [5]int
	}{
		{name: "zero value", values: [5]int{}, want: [5]int{}},
		{name: "sequence", values: [5]int{1, 2, 3, 4, 5}, want: [5]int{5, 4, 3, 2, 1}},
		{name: "mixed", values: [5]int{-2, 7, 0, 7, 9}, want: [5]int{9, 7, 0, 7, -2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.values
			if got := ReverseFive(test.values); got != test.want {
				t.Fatalf("ReverseFive(%v) = %v, want %v", test.values, got, test.want)
			}
			if test.values != original {
				t.Fatalf("ReverseFive mutated input: got %v, want %v", test.values, original)
			}
		})
	}
}

func TestAppendCopy(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		extra  []int
		want   []int
		nilOut bool
	}{
		{name: "both nil", nilOut: true},
		{name: "empty allocated", values: []int{}, extra: []int{}, want: []int{}},
		{name: "combined", values: []int{1, 2}, extra: []int{3, 4}, want: []int{1, 2, 3, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := AppendCopy(test.values, test.extra)
			if (got == nil) != test.nilOut {
				t.Fatalf("AppendCopy() nil = %t, want %t", got == nil, test.nilOut)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("AppendCopy() = %v, want %v", got, test.want)
			}
			if len(got) > 0 && len(test.values) > 0 {
				got[0] = 99
				if test.values[0] == 99 {
					t.Fatal("AppendCopy result shares writable storage with values")
				}
			}
			if len(got) > len(test.values) && len(test.extra) > 0 {
				got[len(test.values)] = 98
				if test.extra[0] == 98 {
					t.Fatal("AppendCopy result shares writable storage with extra")
				}
			}
		})
	}
}
