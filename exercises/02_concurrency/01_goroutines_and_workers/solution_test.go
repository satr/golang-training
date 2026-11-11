package goroutinesworkers

import (
	"slices"
	"testing"
)

func TestProcessItems(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  []string
	}{
		{name: "nil", items: nil, want: nil},
		{name: "empty", items: []string{}, want: nil},
		{name: "one item", items: []string{"hello"}, want: []string{"HELLO"}},
		{name: "preserves order", items: []string{"third", "a", "second"}, want: []string{"THIRD", "A", "SECOND"}},
		{name: "already uppercase", items: []string{"GO", "123"}, want: []string{"GO", "123"}},
		{name: "unicode", items: []string{"café", "blåbær", "κόσμος"}, want: []string{"CAFÉ", "BLÅBÆR", "ΚΌΣΜΟΣ"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ProcessItems(test.items)
			if !slices.Equal(got, test.want) {
				t.Fatalf("ProcessItems(%v) = %v; want %v", test.items, got, test.want)
			}
		})
	}
}

func TestProcessItemsDoesNotModifyInput(t *testing.T) {
	items := []string{"go", "gopher"}
	wantInput := slices.Clone(items)

	_ = ProcessItems(items)

	if !slices.Equal(items, wantInput) {
		t.Fatalf("ProcessItems modified input: got %v; want %v", items, wantInput)
	}
}

func TestProcessItemsManyValues(t *testing.T) {
	items := make([]string, 1_000)
	want := make([]string, len(items))
	for index := range items {
		items[index] = "go"
		want[index] = "GO"
	}

	got := ProcessItems(items)
	if !slices.Equal(got, want) {
		t.Fatal("ProcessItems returned incorrect results for a large input")
	}
}

func TestGoroutineCount(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "negative", n: -10, want: 0},
		{name: "zero", n: 0, want: 0},
		{name: "one", n: 1, want: 1},
		{name: "ten", n: 10, want: 10},
		{name: "one thousand", n: 1_000, want: 1_000},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GoroutineCount(test.n)
			if got != test.want {
				t.Fatalf("GoroutineCount(%d) = %d; want %d", test.n, got, test.want)
			}
		})
	}
}
