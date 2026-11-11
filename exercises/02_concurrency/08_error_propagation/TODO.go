//go:build todo

package errorpropagation

import (
	"context"
	"errors"
)

var (
	ErrInvalidWorkerCount = errors.New("worker count must be positive")
	ErrNilProcessor       = errors.New("processor must not be nil")
	ErrNegativeValue      = errors.New("negative value")
)

// Processor transforms one value or returns an error. Implementations should
// observe ctx and return promptly after cancellation.
type Processor func(ctx context.Context, value int) (int, error)

// Square returns value squared.
//
// Requirements:
//   - Return ctx.Err() when cancellation is already observable.
//   - Reject negative values with an error wrapping ErrNegativeValue.
//   - Return the square and nil for non-negative values.
func Square(ctx context.Context, value int) (int, error) {
	panic("TODO: implement Square")
}

// ProcessConcurrent processes inputs with a bounded worker pool.
//
// Requirements:
//   - Validate workers, process, and ctx before starting goroutines.
//   - Preserve input order on success.
//   - Wrap processing errors with the failing input index.
//   - Cancel sibling work after the first observed processing error.
//   - Make task sends, task receives, and result sends cancellation-aware.
//   - Wait for all internal workers and close owned channels safely.
//   - Return no partial result slice when an error occurs.
func ProcessConcurrent(
	ctx context.Context,
	inputs []int,
	workers int,
	process Processor,
) ([]int, error) {
	panic("TODO: implement ProcessConcurrent")
}
