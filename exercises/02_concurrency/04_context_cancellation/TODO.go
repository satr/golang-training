//go:build todo

package contextcancellation

import "context"

// LongRunningWork produces the integers 0 through n-1. It returns the prefix
// completed before ctx is cancelled. For n <= 0, it returns nil.
func LongRunningWork(ctx context.Context, n int) []int {
	// TODO: Loop from 0 to n-1. On each iteration, use select to check
	// ctx.Done(). If cancelled, return results collected so far.
	// Otherwise, append the current value and continue.
	return nil
}

// FanInWithContext merges in1 and in2 into one channel. It closes the output
// after both inputs close or when ctx is cancelled. A nil input is treated as
// absent rather than as a channel that must be waited on forever.
func FanInWithContext(ctx context.Context, in1, in2 <-chan int) <-chan int {
	// TODO: Spawn a goroutine that reads from in1 and in2 using select,
	// forwarding values to the output channel. Stop and close the output
	// when ctx is cancelled or both inputs are closed.
	return nil
}
