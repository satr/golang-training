package contextpropagation

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestRequestIDGeneratesWhenMissing(t *testing.T) {
	idGen.Store(0)
	ctx := context.Background()
	id := RequestID(ctx)

	got, err := fmtParseID(id)
	if err != nil {
		t.Fatalf("unexpected ID format %q: %v", id, err)
	}
	if got != 1 {
		t.Errorf("RequestID(bg) = %q; want seq 1", id)
	}
}

func TestRequestIDExtractsFromContext(t *testing.T) {
	ctx := WithRequestID(context.Background(), "abc-123")
	got := RequestID(ctx)
	if got != "abc-123" {
		t.Errorf("RequestID(ctx) = %q; want %q", got, "abc-123")
	}
}

func TestRequestIDSequentialGeneration(t *testing.T) {
	idGen.Store(0)

	ctx := context.Background()
	ids := make([]string, 3)
	for i := range ids {
		ids[i] = RequestID(ctx)
	}

	if ids[0] != "req-1" || ids[1] != "req-2" || ids[2] != "req-3" {
		t.Errorf("sequential IDs = %v; want [req-1 req-2 req-3]", ids)
	}
}

func TestNewPipeline(t *testing.T) {
	p := NewPipeline(
		func(ctx context.Context, in []byte) ([]byte, error) {
			return in, nil
		},
	)
	if len(p.steps) != 1 {
		t.Errorf("steps = %d; want 1", len(p.steps))
	}
}

func TestPipelineThreeSteps(t *testing.T) {
	var callOrder []string

	steps := []Step{
		func(ctx context.Context, in []byte) ([]byte, error) {
			callOrder = append(callOrder, "a")
			return append([]byte("a-"), in...), nil
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			callOrder = append(callOrder, "b")
			return append([]byte("b-"), in...), nil
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			callOrder = append(callOrder, "c")
			return append([]byte("c-"), in...), nil
		},
	}

	p := NewPipeline(steps...)
	got, err := p.Run(context.Background(), []byte("payload"))
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got := string(got); got != "a-b-c-payload" {
		t.Errorf("Run() = %q; want %q", got, "a-b-c-payload")
	}
	expectedOrder := []string{"a", "b", "c"}
	if callOrder[0] != expectedOrder[0] ||
		callOrder[1] != expectedOrder[1] ||
		callOrder[2] != expectedOrder[2] {
		t.Errorf("call order = %v; want %v", callOrder, expectedOrder)
	}
}

func TestPipelineStopsOnCancellation(t *testing.T) {
	var counter atomic.Int32

	steps := []Step{
		func(ctx context.Context, in []byte) ([]byte, error) {
			counter.Add(1)
			return in, nil
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			cancelCh.(chan struct{}) <- struct{}{}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-counter.Chan():
			}
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			counter.Add(1)
			return in, nil
		},
	}

	// This test uses a simpler approach: pre-cancel the context.
	ctx, cancel := context.WithCancel(context.Background())

	// Build a pipeline where the second step triggers cancellation.
	sideways := []Step{
		func(ctx context.Context, in []byte) ([]byte, error) {
			return in, nil
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			cancel()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			return in, nil
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			t.Error("third step should not run after cancellation")
			return in, nil
		},
	}

	p := NewPipeline(sideways...)
	_, err := p.Run(ctx, []byte("data"))
	if err != nil {
		// The second step calls cancel() then synchronously returns; Run may
		// either succeed (if step returned before the next ctx.Done() check)
		// or see the cancellation on the next iteration.
		// Either outcome is acceptable; the key point is that the third step
		// does not report an error from its t.Error call above.
		if err != context.Canceled {
			t.Fatalf("Run() error = %v; want context.Canceled", err)
		}
	}
}

func TestPipelineDeadlineExceeded(t *testing.T) {
	var blocked chan struct{}

	steps := []Step{
		func(_ context.Context, in []byte) ([]byte, error) {
			<-blocked // blocks forever
			return in, nil
		},
	}

	p := NewPipeline(steps...)

	// Use a channel to synchronize: cancel the context after a tick so that
	// Run detects the deadline before the step returns.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		cancel()
	}()

	// Since the step blocks, rely on the next iteration's ctx.Done() check.
	// But there is no next iteration because there is only one step. So we
	// use a different approach: a context with deadline and a step that
	// cooperates.
	ctx2, cancel2 := context.WithCancelCause(context.Background())
	defer cancel2(nil)

	steps2 := []Step{
		func(ctx context.Context, in []byte) ([]byte, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			cancel2(context.DeadlineExceeded)
			return in, nil
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			return in, nil
		},
	}
	blocked = make(chan struct{})
	_ = blocked

	p2 := NewPipeline(steps2...)
	_, err := p2.Run(ctx2, []byte("x"))
	if err != context.DeadlineExceeded {
		t.Errorf("Run() error = %v; want context.DeadlineExceeded", err)
	}
}

func TestContextValuesPropagate(t *testing.T) {
	idGen.Store(0)
	ctx := WithRequestID(context.Background(), "fixed-id")

	var captured []string
	steps := []Step{
		func(ctx context.Context, in []byte) ([]byte, error) {
			captured = append(captured, RequestID(ctx))
			return in, nil
		},
		func(ctx context.Context, in []byte) ([]byte, error) {
			captured = append(captured, RequestID(ctx))
			return in, nil
		},
	}

	_, err := NewPipeline(steps...).Run(ctx, nil)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	for i, id := range captured {
		if id != "fixed-id" {
			t.Errorf("captured[%d] = %q; want %q", i, id, "fixed-id")
		}
	}
}

// Helper to parse "req-N" into N.
var fmtParseID = func(id string) (int, error) {
	var n int
	_, err := fmt.Sscanf(id, "req-%d", &n)
	return n, err
}

var cancelCh chan struct{}

func init() {
	cancelCh = make(chan struct{}, 1)
}
