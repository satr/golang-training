//go:build todo

package boundedworkerpools

import "context"

// WorkerPool squares inputs using at most workers concurrent workers.
// Results are returned in completion order, which may differ from input order.
func WorkerPool(inputs []int, workers int) []int {
	// TODO: Create a tasks channel and a results channel. Launch N
	// workers that range over tasks and send v*v to results. Fan-in
	// all results back into a slice.
	return nil
}

// WorkerPoolWithTimeout squares inputs using at most workers concurrent workers.
// It stops scheduling and processing work when ctx is cancelled and returns the
// results delivered before the pipeline stops. Cancellation races with channel
// operations, so a value already in flight may still be included.
func WorkerPoolWithTimeout(ctx context.Context, inputs []int, workers int) []int {
	// TODO: Like WorkerPool, but use select with ctx.Done() to
	// allow early return when the context is cancelled.
	return nil
}
