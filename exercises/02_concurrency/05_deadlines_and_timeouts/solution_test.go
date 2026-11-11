package deadlinesandtimeouts

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestBoundedComputeCompletes(t *testing.T) {
	tests := []struct {
		name   string
		inputs []int
		want   []int
	}{
		{name: "empty", inputs: nil, want: nil},
		{name: "single", inputs: []int{7}, want: []int{49}},
		{name: "multiple", inputs: []int{2, 3, 4}, want: []int{4, 9, 16}},
		{name: "negative zero and duplicate", inputs: []int{-3, 0, 2, -3}, want: []int{9, 0, 4, 9}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := BoundedCompute(context.Background(), test.inputs)
			if err != nil {
				t.Fatalf("BoundedCompute() error = %v; want nil", err)
			}
			if !slices.Equal(got, test.want) {
				t.Fatalf("BoundedCompute() = %v; want %v", got, test.want)
			}
		})
	}
}

func TestBoundedComputeAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := BoundedCompute(ctx, []int{1, 2, 3})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("BoundedCompute() error = %v; want context.Canceled", err)
	}
	if len(got) != 0 {
		t.Fatalf("BoundedCompute() = %v; want no results", got)
	}
}

func TestBoundedComputeExpiredDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	got, err := BoundedCompute(ctx, []int{1, 2, 3})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("BoundedCompute() error = %v; want context.DeadlineExceeded", err)
	}
	if len(got) != 0 {
		t.Fatalf("BoundedCompute() = %v; want no results", got)
	}
}

func TestBoundedComputeReturnsPartialResults(t *testing.T) {
	// Done closes deterministically after two loop checks, avoiding sleeps and
	// scheduler-dependent assertions.
	ctx := newCancelAfterChecksContext(2)

	got, err := BoundedCompute(ctx, []int{2, 3, 4, 5})
	want := []int{4, 9}
	if !slices.Equal(got, want) {
		t.Fatalf("BoundedCompute() = %v; want partial results %v", got, want)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("BoundedCompute() error = %v; want context.Canceled", err)
	}
}

func TestSafeSleep(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
	}{
		{name: "negative duration", d: -time.Second},
		{name: "zero duration", d: 0},
		{name: "timer completes", d: time.Millisecond},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := SafeSleep(context.Background(), test.d); err != nil {
				t.Fatalf("SafeSleep() error = %v; want nil", err)
			}
		})
	}
}

func TestSafeSleepAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Cancellation on entry has precedence even when no sleep is required.
	err := SafeSleep(ctx, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("SafeSleep() error = %v; want context.Canceled", err)
	}
}

func TestSafeSleepReturnsAtDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- SafeSleep(ctx, time.Second)
	}()

	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("SafeSleep() error = %v; want context.DeadlineExceeded", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("SafeSleep did not return promptly after the deadline")
	}
}

// cancelAfterChecksContext is a deterministic test context. Done remains open
// for openChecks calls and closes on the following call.
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
