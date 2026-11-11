package contextcancellation

import "context"

// LongRunningWork produces the integers 0 through n-1. It returns the prefix
// completed before ctx is cancelled. For n <= 0, it returns nil.
func LongRunningWork(ctx context.Context, n int) []int {
	if n <= 0 {
		return nil
	}

	// Do not eagerly allocate capacity n: n may be very large while ctx is
	// canceled early.
	var results []int
	for value := range n {
		select {
		case <-ctx.Done():
			return results
		default:
			results = append(results, value)
		}
	}
	return results
}

// FanInWithContext merges in1 and in2 into one channel. It closes the output
// after both inputs close or when ctx is cancelled. A nil input is treated as
// absent rather than as a channel that must be waited on forever.
func FanInWithContext(ctx context.Context, in1, in2 <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		send := func(value int) bool {
			select {
			case out <- value:
				return true
			case <-ctx.Done():
				return false
			}
		}

		for in1 != nil || in2 != nil {
			select {
			case <-ctx.Done():
				return

			case value, ok := <-in1:
				if !ok {
					// A nil channel disables this select case.
					in1 = nil
					continue
				}
				if !send(value) {
					return
				}

			case value, ok := <-in2:
				if !ok {
					in2 = nil
					continue
				}
				if !send(value) {
					return
				}
			}
		}
	}()

	return out
}
