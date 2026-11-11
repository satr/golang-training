package contextcancellation

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestLongRunningWork(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{name: "negative count", n: -1, want: nil},
		{name: "zero items", n: 0, want: nil},
		{name: "one item", n: 1, want: []int{0}},
		{name: "multiple items", n: 5, want: []int{0, 1, 2, 3, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := LongRunningWork(context.Background(), test.n)
			if !slices.Equal(got, test.want) {
				t.Fatalf("LongRunningWork(ctx, %d) = %v; want %v", test.n, got, test.want)
			}
		})
	}
}

func TestLongRunningWorkWithAlreadyCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got := LongRunningWork(ctx, 10)
	if len(got) != 0 {
		t.Fatalf("LongRunningWork() = %v; want no results", got)
	}
}

func TestLongRunningWorkStopsDuringProcessing(t *testing.T) {
	// This test context closes Done deterministically after three checks. It
	// avoids relying on sleeps or scheduler timing to simulate mid-work
	// cancellation.
	ctx := newCancelAfterChecksContext(3)

	got := LongRunningWork(ctx, 10)
	want := []int{0, 1, 2}
	if !slices.Equal(got, want) {
		t.Fatalf("LongRunningWork() = %v; want %v", got, want)
	}
	if ctx.Err() != context.Canceled {
		t.Fatalf("ctx.Err() = %v; want context.Canceled", ctx.Err())
	}
}

func TestFanInWithContext(t *testing.T) {
	tests := []struct {
		name string
		in1  []int
		in2  []int
		nil1 bool
		nil2 bool
	}{
		{name: "both closed and empty"},
		{name: "both nil", nil1: true, nil2: true},
		{name: "first only", in1: []int{1, 2}, nil2: true},
		{name: "second only", in2: []int{3, 4}, nil1: true},
		{name: "both contain values", in1: []int{1, 2}, in2: []int{3, 4}},
		{name: "duplicates", in1: []int{2, 2}, in2: []int{2, 3}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var in1, in2 <-chan int
			if !test.nil1 {
				in1 = channelOf(test.in1...)
			}
			if !test.nil2 {
				in2 = channelOf(test.in2...)
			}

			got := collectUntilClosed(t, FanInWithContext(context.Background(), in1, in2))
			want := append(append([]int(nil), test.in1...), test.in2...)
			slices.Sort(got)
			slices.Sort(want)

			if !slices.Equal(got, want) {
				t.Fatalf("FanInWithContext() = %v; want %v", got, want)
			}
		})
	}
}

func TestFanInClosesWhenCancelledWhileWaitingForInput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in1 := make(chan int)
	in2 := make(chan int)
	out := FanInWithContext(ctx, in1, in2)

	cancel()

	got := collectUntilClosed(t, out)
	if len(got) != 0 {
		t.Fatalf("FanInWithContext() = %v after cancellation; want no results", got)
	}
}

func TestFanInCancellationUnblocksPendingOutput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan int)
	sent := make(chan struct{})

	out := FanInWithContext(ctx, in, nil)
	go func() {
		in <- 7
		close(sent)
		close(in)
	}()

	// The fan-in has received the value. Since out has no receiver yet, it is
	// now at (or immediately approaching) the cancellation-aware output send.
	<-sent
	cancel()

	got := collectUntilClosed(t, out)
	if len(got) > 1 || len(got) == 1 && got[0] != 7 {
		t.Fatalf("FanInWithContext() = %v; want either no value or [7]", got)
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
			t.Fatal("timed out waiting for output channel to close")
		}
	}
}

// cancelAfterChecksContext is a deterministic test context. Done remains open
// for openChecks calls and closes on the next call.
type cancelAfterChecksContext struct {
	context.Context

	done       chan struct{}
	openChecks int

	mu     sync.Mutex
	checks int
	closed bool
}

func newCancelAfterChecksContext(openChecks int) *cancelAfterChecksContext {
	return &cancelAfterChecksContext{
		Context:    context.Background(),
		done:       make(chan struct{}),
		openChecks: openChecks,
	}
}

func (ctx *cancelAfterChecksContext) Done() <-chan struct{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.checks++
	if ctx.checks > ctx.openChecks && !ctx.closed {
		close(ctx.done)
		ctx.closed = true
	}
	return ctx.done
}

func (ctx *cancelAfterChecksContext) Err() error {
	select {
	case <-ctx.done:
		return context.Canceled
	default:
		return nil
	}
}
