//go:build todo

package deadlinesandtimeouts

import (
	"context"
	"time"
)

// BoundedCompute squares each input in order. If ctx is cancelled, it returns
// the completed prefix together with ctx.Err().
func BoundedCompute(ctx context.Context, inputs []int) ([]int, error) {
	// TODO: Loop over inputs. Before each squaring, check ctx.Done()
	// with select. If done, return partial results and ctx.Err()
	// wrapped with a descriptive message.
	return nil, nil
}

// SafeSleep waits for d or returns early when ctx is cancelled. Cancellation
// that is already observable on entry takes precedence over a non-positive d.
func SafeSleep(ctx context.Context, d time.Duration) error {
	// TODO: Use select on ctx.Done() and time.After(d).
	// Return ctx.Err() if the context expires first.
	return nil
}
