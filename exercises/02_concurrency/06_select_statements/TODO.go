//go:build todo

package selectstatements

import (
	"context"
	"time"
)

// TickerPulse is emitted by TickerMerge for every observed ticker event.
// A tagged result type would avoid sentinel-value ambiguity in production code,
// but this exercise intentionally keeps the output type as int.
const TickerPulse = -1

// Multiplexer merges all non-nil sources into one output channel. The output is
// closed after every source closes or after ctx is cancelled.
func Multiplexer(ctx context.Context, sources []<-chan int) <-chan int {
	// TODO: Use select to read from all sources. Set closed channels to nil
	// to remove them from the select. Close output when ctx is done or
	// all sources are closed.
	return nil
}

// TickerMerge forwards values from in and emits TickerPulse periodically.
// It closes the output when ctx is cancelled or in is closed. A nil in channel
// enables ticker-only operation until cancellation.
func TickerMerge(ctx context.Context, in <-chan int, interval time.Duration) <-chan int {
	// TODO: Start a time.Ticker at the given interval. Use select to read
	// from both the input channel and the ticker. Close output when ctx is
	// done or input is closed and ticker has no pending pulses.
	return nil
}
