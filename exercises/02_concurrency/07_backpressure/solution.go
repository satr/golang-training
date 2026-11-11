package backpressure

import "context"

// ProcessBatch sends every value to ch in order. It blocks whenever the
// receiver and channel capacity cannot accept the next value. It does not close
// ch because the caller may have other producers.
func ProcessBatch(ch chan<- int, values []int) {
	for _, value := range values {
		ch <- value
	}
}

// MergeBuffered merges ch1 and ch2 into an output channel with capacity
// bufSize. Nil inputs are treated as absent. The output closes after both
// inputs close or ctx is cancelled.
func MergeBuffered(ctx context.Context, ch1, ch2 <-chan int, bufSize int) <-chan int {
	if bufSize < 0 {
		panic("MergeBuffered: negative buffer size")
	}

	out := make(chan int, bufSize)

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

		for ch1 != nil || ch2 != nil {
			select {
			case <-ctx.Done():
				return

			case value, ok := <-ch1:
				if !ok {
					// A nil channel disables this case in future selections.
					ch1 = nil
					continue
				}
				if !send(value) {
					return
				}

			case value, ok := <-ch2:
				if !ok {
					ch2 = nil
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

// Drain collects values until ch closes or ctx is cancelled. Cancellation
// already observable on entry takes precedence over buffered channel values.
func Drain(ctx context.Context, ch <-chan int) []int {
	if ctx.Err() != nil {
		return nil
	}

	var values []int
	for {
		select {
		case <-ctx.Done():
			return values
		case value, ok := <-ch:
			if !ok {
				return values
			}
			values = append(values, value)
		}
	}
}
