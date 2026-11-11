//go:build todo

package backpressure

import "context"

// ProcessBatch sends all values in vals to ch, blocking when the channel
// buffer is full. Returns when every value has been sent.
func ProcessBatch(ch chan<- int, vals []int) {
	// TODO: Range over vals and send each value to ch.
	// The send blocks naturally when the channel buffer is full.
}

// MergeBuffered merges ch1 and ch2 into an output channel with capacity
// bufSize. Nil inputs are treated as absent. The output closes after both
// inputs close or ctx is cancelled.
func MergeBuffered(ctx context.Context, ch1, ch2 <-chan int, bufSize int) <-chan int {
	// TODO: Create a buffered output channel. Launch a goroutine that
	// uses select to read from ch1 and ch2, forwarding values to out.
	// Close out when both inputs are closed.
	return nil
}

// Drain collects values until ch closes or ctx is cancelled. Cancellation
// already observable on entry takes precedence over buffered channel values.
func Drain(ctx context.Context, ch <-chan int) []int {
	// TODO: Use select to receive from ch and ctx.Done(). Append each
	// value to a slice. Return the slice when the channel closes or
	// the context is cancelled.
	return nil
}
