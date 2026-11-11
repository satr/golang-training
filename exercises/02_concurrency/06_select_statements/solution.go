package selectstatements

import (
	"context"
	"sync"
	"time"
)

// TickerPulse is emitted by TickerMerge for every observed ticker event.
// A tagged result type would avoid sentinel-value ambiguity in production code,
// but this exercise intentionally keeps the output type as int.
const TickerPulse = -1

// Multiplexer merges all non-nil sources into one output channel. The output is
// closed after every source closes or after ctx is cancelled.
func Multiplexer(ctx context.Context, sources []<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	forward := func(source <-chan int) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-source:
				if !ok {
					return
				}

				select {
				case out <- value:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, source := range sources {
		if source == nil {
			continue
		}
		wg.Add(1)
		go forward(source)
	}

	// Only the coordinator closes out, after no forwarder can still send.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// TickerMerge forwards values from in and emits TickerPulse periodically.
// It closes the output when ctx is cancelled or in is closed. A nil in channel
// enables ticker-only operation until cancellation.
func TickerMerge(ctx context.Context, in <-chan int, interval time.Duration) <-chan int {
	if interval <= 0 {
		panic("TickerMerge: interval must be positive")
	}

	out := make(chan int)

	go func() {
		defer close(out)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		send := func(value int) bool {
			select {
			case out <- value:
				return true
			case <-ctx.Done():
				return false
			}
		}

		for {
			select {
			case <-ctx.Done():
				return

			case value, ok := <-in:
				if !ok {
					return
				}
				if !send(value) {
					return
				}

			case <-ticker.C:
				if !send(TickerPulse) {
					return
				}
			}
		}
	}()

	return out
}
