package errorpropagation

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSquare(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		want    int
		wantErr error
	}{
		{name: "zero", value: 0, want: 0},
		{name: "positive", value: 7, want: 49},
		{name: "negative", value: -3, wantErr: ErrNegativeValue},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Square(context.Background(), test.value)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Square(%d) error = %v; want %v", test.value, err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("Square(%d) = %d; want %d", test.value, got, test.want)
			}
		})
	}
}

func TestSquareAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Square(ctx, -1)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Square() error = %v; want context.Canceled", err)
	}
}

func TestProcessConcurrent(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []int
		workers int
		want    []int
	}{
		{name: "empty", inputs: nil, workers: 1, want: nil},
		{name: "one", inputs: []int{3}, workers: 1, want: []int{9}},
		{name: "multiple", inputs: []int{1, 2, 3, 4}, workers: 2, want: []int{1, 4, 9, 16}},
		{name: "more workers than inputs", inputs: []int{2, 5}, workers: 10, want: []int{4, 25}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ProcessConcurrent(context.Background(), test.inputs, test.workers, Square)
			if err != nil {
				t.Fatalf("ProcessConcurrent() error = %v; want nil", err)
			}
			if !slices.Equal(got, test.want) {
				t.Fatalf("ProcessConcurrent() = %v; want %v", got, test.want)
			}
		})
	}
}

func TestProcessConcurrentPreservesInputOrder(t *testing.T) {
	secondFinished := make(chan struct{})
	var once sync.Once

	process := func(ctx context.Context, value int) (int, error) {
		switch value {
		case 1:
			select {
			case <-secondFinished:
			case <-ctx.Done():
				return 0, ctx.Err()
			}
		case 2:
			once.Do(func() { close(secondFinished) })
		}
		return value * 10, nil
	}

	got, err := ProcessConcurrent(context.Background(), []int{1, 2, 3}, 2, process)
	if err != nil {
		t.Fatalf("ProcessConcurrent() error = %v; want nil", err)
	}
	want := []int{10, 20, 30}
	if !slices.Equal(got, want) {
		t.Fatalf("ProcessConcurrent() = %v; want ordered results %v", got, want)
	}
}

func TestProcessConcurrentValidatesArguments(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		process Processor
		wantErr error
	}{
		{name: "zero workers", workers: 0, process: Square, wantErr: ErrInvalidWorkerCount},
		{name: "negative workers", workers: -1, process: Square, wantErr: ErrInvalidWorkerCount},
		{name: "nil processor", workers: 1, process: nil, wantErr: ErrNilProcessor},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ProcessConcurrent(context.Background(), []int{1}, test.workers, test.process)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("ProcessConcurrent() error = %v; want %v", err, test.wantErr)
			}
			if got != nil {
				t.Fatalf("ProcessConcurrent() = %v; want nil results", got)
			}
		})
	}
}

func TestProcessConcurrentWrapsProcessingError(t *testing.T) {
	got, err := ProcessConcurrent(context.Background(), []int{2, -3, 4}, 1, Square)
	if !errors.Is(err, ErrNegativeValue) {
		t.Fatalf("ProcessConcurrent() error = %v; want ErrNegativeValue", err)
	}
	if !strings.Contains(err.Error(), "input 1") {
		t.Fatalf("ProcessConcurrent() error = %q; want failing input index", err)
	}
	if got != nil {
		t.Fatalf("ProcessConcurrent() = %v; want nil partial results", got)
	}
}

func TestProcessConcurrentCancelsSiblingWork(t *testing.T) {
	workErr := errors.New("worker failed")
	blockerStarted := make(chan struct{})
	blockerStopped := make(chan struct{})
	var startOnce sync.Once
	var stopOnce sync.Once

	process := func(ctx context.Context, value int) (int, error) {
		if value == 1 {
			startOnce.Do(func() { close(blockerStarted) })
			<-ctx.Done()
			stopOnce.Do(func() { close(blockerStopped) })
			return 0, ctx.Err()
		}

		<-blockerStarted
		return 0, workErr
	}

	type outcome struct {
		values []int
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		values, err := ProcessConcurrent(context.Background(), []int{1, 2}, 2, process)
		done <- outcome{values: values, err: err}
	}()

	select {
	case result := <-done:
		if !errors.Is(result.err, workErr) {
			t.Fatalf("ProcessConcurrent() error = %v; want worker error", result.err)
		}
		if result.values != nil {
			t.Fatalf("ProcessConcurrent() = %v; want nil results", result.values)
		}
		select {
		case <-blockerStopped:
		default:
			t.Fatal("ProcessConcurrent returned before cancelled sibling stopped")
		}
	case <-time.After(time.Second):
		t.Fatal("ProcessConcurrent did not return after a worker error")
	}
}

func TestProcessConcurrentPropagatesParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	stopped := make(chan struct{})

	process := func(ctx context.Context, value int) (int, error) {
		close(started)
		<-ctx.Done()
		close(stopped)
		return 0, ctx.Err()
	}

	done := make(chan error, 1)
	go func() {
		_, err := ProcessConcurrent(ctx, []int{1}, 1, process)
		done <- err
	}()

	<-started
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("ProcessConcurrent() error = %v; want context.Canceled", err)
		}
		select {
		case <-stopped:
		default:
			t.Fatal("ProcessConcurrent returned before processor stopped")
		}
	case <-time.After(time.Second):
		t.Fatal("ProcessConcurrent did not return after parent cancellation")
	}
}

func TestProcessConcurrentAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := ProcessConcurrent(ctx, []int{1}, 1, Square)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ProcessConcurrent() error = %v; want context.Canceled", err)
	}
	if got != nil {
		t.Fatalf("ProcessConcurrent() = %v; want nil results", got)
	}
}
