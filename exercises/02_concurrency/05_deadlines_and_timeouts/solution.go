package deadlinesandtimeouts

import (
	"context"
	"time"
)

// BoundedCompute squares each input in order. If ctx is cancelled, it returns
// the completed prefix together with ctx.Err().
func BoundedCompute(ctx context.Context, inputs []int) ([]int, error) {
	// Check before allocating or starting work. This also defines the behavior
	// for an already-cancelled context with empty input.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Grow on demand instead of reserving all input space before knowing how
	// much work the context will allow us to complete.
	var results []int
	for _, value := range inputs {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
			results = append(results, value*value)
		}
	}
	return results, nil
}

// SafeSleep waits for d or returns early when ctx is cancelled. Cancellation
// that is already observable on entry takes precedence over a non-positive d.
func SafeSleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d <= 0 {
		return nil
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
