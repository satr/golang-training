package backpressure

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestProcessBatch(t *testing.T) {
	tests := []struct {
		name    string
		values  []int
		bufSize int
	}{
		{name: "empty", values: nil, bufSize: 0},
		{name: "one value", values: []int{42}, bufSize: 1},
		{name: "fits in buffer", values: []int{1, 2, 3}, bufSize: 3},
		{name: "requires consumer progress", values: []int{1, 2, 3, 4, 5}, bufSize: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := make(chan int, test.bufSize)
			producerDone := make(chan struct{})

			go func() {
				ProcessBatch(ch, test.values)
				close(ch)
				close(producerDone)
			}()

			got := collectUntilClosed(t, ch)
			<-producerDone
			if !slices.Equal(got, test.values) {
				t.Fatalf("ProcessBatch() produced %v; want %v", got, test.values)
			}
		})
	}
}

func TestProcessBatchDoesNotCloseChannel(t *testing.T) {
	ch := make(chan int, 2)
	ProcessBatch(ch, []int{1})

	// A second producer can still use the channel.
	ch <- 2
	close(ch)

	got := collectUntilClosed(t, ch)
	want := []int{1, 2}
	if !slices.Equal(got, want) {
		t.Fatalf("channel contains %v; want %v", got, want)
	}
}

func TestMergeBuffered(t *testing.T) {
	tests := []struct {
		name    string
		ch1     []int
		ch2     []int
		bufSize int
	}{
		{name: "both empty", bufSize: 0},
		{name: "first only", ch1: []int{1, 2}, bufSize: 1},
		{name: "second only", ch2: []int{3, 4}, bufSize: 3},
		{name: "both inputs", ch1: []int{1, 3}, ch2: []int{2, 4}, bufSize: 4},
		{name: "larger than buffer", ch1: []int{1, 2, 3, 4}, ch2: []int{5, 6, 7, 8}, bufSize: 2},
		{name: "duplicates", ch1: []int{2, 2}, ch2: []int{2, 3}, bufSize: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			merged := MergeBuffered(
				context.Background(),
				channelOf(test.ch1...),
				channelOf(test.ch2...),
				test.bufSize,
			)
			if got := cap(merged); got != test.bufSize {
				t.Fatalf("output capacity = %d; want %d", got, test.bufSize)
			}

			got := collectUntilClosed(t, merged)
			want := append(append([]int(nil), test.ch1...), test.ch2...)
			slices.Sort(got)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Fatalf("MergeBuffered() = %v; want %v", got, want)
			}
		})
	}
}

func TestMergeBufferedTreatsNilInputsAsAbsent(t *testing.T) {
	t.Run("one nil", func(t *testing.T) {
		merged := MergeBuffered(context.Background(), nil, channelOf(1, 2), 1)
		if got, want := collectUntilClosed(t, merged), []int{1, 2}; !slices.Equal(got, want) {
			t.Fatalf("MergeBuffered() = %v; want %v", got, want)
		}
	})

	t.Run("both nil", func(t *testing.T) {
		merged := MergeBuffered(context.Background(), nil, nil, 1)
		if got := collectUntilClosed(t, merged); len(got) != 0 {
			t.Fatalf("MergeBuffered() = %v; want no values", got)
		}
	})
}

func TestMergeBufferedClosesOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	merged := MergeBuffered(ctx, make(chan int), make(chan int), 2)

	cancel()

	if got := collectUntilClosed(t, merged); len(got) != 0 {
		t.Fatalf("MergeBuffered() = %v after cancellation; want no values", got)
	}
}

func TestMergeBufferedRejectsNegativeBuffer(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("MergeBuffered did not panic for a negative buffer size")
		}
	}()

	MergeBuffered(context.Background(), nil, nil, -1)
}

func TestDrain(t *testing.T) {
	tests := []struct {
		name   string
		values []int
	}{
		{name: "empty"},
		{name: "single", values: []int{9}},
		{name: "multiple", values: []int{1, 2, 3}},
		{name: "duplicates", values: []int{2, 2, 3}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Drain(context.Background(), channelOf(test.values...))
			if !slices.Equal(got, test.values) {
				t.Fatalf("Drain() = %v; want %v", got, test.values)
			}
		})
	}
}

func TestDrainAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got := Drain(ctx, channelOf(1, 2, 3))
	if len(got) != 0 {
		t.Fatalf("Drain() = %v; want no values for cancellation on entry", got)
	}
}

func TestDrainCancellationUnblocksReceive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan []int, 1)

	go func() {
		done <- Drain(ctx, make(chan int))
	}()
	cancel()

	select {
	case got := <-done:
		if len(got) != 0 {
			t.Fatalf("Drain() = %v; want no values", got)
		}
	case <-time.After(time.Second):
		t.Fatal("Drain did not return after cancellation")
	}
}

func channelOf(values ...int) <-chan int {
	ch := make(chan int, len(values))
	for _, value := range values {
		ch <- value
	}
	close(ch)
	return ch
}

func collectUntilClosed(t *testing.T, ch <-chan int) []int {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	var values []int
	for {
		select {
		case value, ok := <-ch:
			if !ok {
				return values
			}
			values = append(values, value)
		case <-timer.C:
			t.Fatal("timed out waiting for channel to close")
		}
	}
}
