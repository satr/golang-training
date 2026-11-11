package errorpropagation

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrInvalidWorkerCount reports a non-positive worker count.
	ErrInvalidWorkerCount = errors.New("worker count must be positive")
	// ErrNilProcessor reports a missing Processor function.
	ErrNilProcessor = errors.New("processor must not be nil")
	// ErrNegativeValue reports that Square received a negative value.
	ErrNegativeValue = errors.New("negative value")
)

// Processor transforms one value or returns an error. Implementations should
// observe ctx and return promptly after cancellation.
type Processor func(ctx context.Context, value int) (int, error)

type indexedTask struct {
	index int
	value int
}

type indexedResult struct {
	index int
	value int
	err   error
}

// Square returns value squared. Negative values are rejected with an error
// wrapping ErrNegativeValue.
func Square(ctx context.Context, value int) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if value < 0 {
		return 0, fmt.Errorf("%w: %d", ErrNegativeValue, value)
	}
	return value * value, nil
}

// ProcessConcurrent processes inputs with a bounded number of workers.
// Results preserve input order. On the first observed processing error, it
// cancels sibling work, waits for internal goroutines to stop, and returns a
// wrapped error with no partial result slice.
func ProcessConcurrent(
	ctx context.Context,
	inputs []int,
	workers int,
	process Processor,
) ([]int, error) {
	if workers <= 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidWorkerCount, workers)
	}
	if process == nil {
		return nil, ErrNilProcessor
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(inputs) == 0 {
		return nil, nil
	}
	if workers > len(inputs) {
		workers = len(inputs)
	}

	workCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	tasks := make(chan indexedTask)
	results := make(chan indexedResult, workers)
	feederDone := make(chan struct{})

	go func() {
		defer close(feederDone)
		defer close(tasks)
		for index, value := range inputs {
			select {
			case tasks <- indexedTask{index: index, value: value}:
			case <-workCtx.Done():
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-workCtx.Done():
					return
				case task, ok := <-tasks:
					if !ok {
						return
					}

					value, err := process(workCtx, task.value)
					result := indexedResult{
						index: task.index,
						value: value,
						err:   err,
					}

					select {
					case results <- result:
					case <-workCtx.Done():
						return
					}

					if err != nil {
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	values := make([]int, len(inputs))
	var firstErr error
	for result := range results {
		if result.err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("process input %d: %w", result.index, result.err)
				cancel()
			}
			continue
		}
		values[result.index] = result.value
	}
	<-feederDone

	if firstErr != nil {
		return nil, firstErr
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return values, nil
}
