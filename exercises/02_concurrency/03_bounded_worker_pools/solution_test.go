package boundedworkerpools

import (
	"context"
	"sort"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []int
		workers int
		want    []int
	}{
		{name: "empty", inputs: nil, workers: 2, want: nil},
		{name: "single", inputs: []int{3}, workers: 1, want: []int{9}},
		{name: "multiple", inputs: []int{1, 2, 3, 4, 5}, workers: 2, want: []int{1, 4, 9, 16, 25}},
		{name: "duplicates and negatives", inputs: []int{-3, 2, -3, 0}, workers: 3, want: []int{9, 4, 9, 0}},
		{name: "more workers than inputs", inputs: []int{2, 4}, workers: 20, want: []int{4, 16}},
		{name: "zero workers defaults to one", inputs: []int{2, 3}, workers: 0, want: []int{4, 9}},
		{name: "negative workers defaults to one", inputs: []int{2, 3}, workers: -2, want: []int{4, 9}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := WorkerPool(test.inputs, test.workers)
			assertSameInts(t, got, test.want)
		})
	}
}

func TestWorkerPoolWithTimeoutCompletes(t *testing.T) {
	inputs := []int{-3, 0, 2, 5, 2}
	want := []int{9, 0, 4, 25, 4}

	got := WorkerPoolWithTimeout(context.Background(), inputs, 3)

	assertSameInts(t, got, want)
}

func TestWorkerPoolWithTimeoutAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got := WorkerPoolWithTimeout(ctx, []int{1, 2, 3}, 2)
	if len(got) != 0 {
		t.Fatalf("got %v; want no results for an already-cancelled context", got)
	}
}

func TestWorkerPoolWithTimeoutCancellationReturnsPromptly(t *testing.T) {
	// Repetition makes send/receive/close races visible under go test -race.
	for iteration := 0; iteration < 50; iteration++ {
		ctx, cancel := context.WithCancel(context.Background())
		inputs := make([]int, 100_000)
		for i := range inputs {
			inputs[i] = i
		}

		done := make(chan []int, 1)
		go func() {
			done <- WorkerPoolWithTimeout(ctx, inputs, 4)
		}()

		cancel()

		select {
		case got := <-done:
			if len(got) > len(inputs) {
				t.Fatalf("iteration %d returned %d results for %d inputs", iteration, len(got), len(inputs))
			}
			assertResultsBelongToInputs(t, got, inputs)
		case <-time.After(time.Second):
			t.Fatalf("iteration %d did not stop after cancellation", iteration)
		}
	}
}

func assertSameInts(t *testing.T, got, want []int) {
	t.Helper()

	gotCopy := append([]int(nil), got...)
	wantCopy := append([]int(nil), want...)
	sort.Ints(gotCopy)
	sort.Ints(wantCopy)

	if len(gotCopy) != len(wantCopy) {
		t.Fatalf("got %v; want %v", gotCopy, wantCopy)
	}
	for i := range wantCopy {
		if gotCopy[i] != wantCopy[i] {
			t.Fatalf("got %v; want %v", gotCopy, wantCopy)
		}
	}
}

func assertResultsBelongToInputs(t *testing.T, results, inputs []int) {
	t.Helper()

	available := make(map[int]int, len(inputs))
	for _, value := range inputs {
		available[value*value]++
	}
	for _, result := range results {
		if available[result] == 0 {
			t.Fatalf("result %d was not produced by an unconsumed input", result)
		}
		available[result]--
	}
}
