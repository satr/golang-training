//go:build todo

package retries

import (
	"context"
	"sync"
	"time"
)

// RetryConfig holds backoff and budget parameters for a retry attempt.
type RetryConfig struct {
	MaxRetries int
	// ...
}

// Result is what the retried function returns.
type Result struct {
	// ...
}

// Budget enforces a system-wide cap on retry attempts within a sliding window.
type Budget struct {
	// ...
}

// NewBudget creates a budget allowing up to max retries within the given window.
func NewBudget(max int, window time.Duration) *Budget {
	// TODO: Return a *Budget configured with max and window.
	return nil
}

// Allow reports whether a retry is permitted under the current budget.
func (b *Budget) Allow() bool {
	// TODO: Lock the budget mutex, remove attempts older than the window,
	// check if the attempt count is below the max, and if so record a new
	// attempt and return true. Otherwise return false.
	return false
}

// Execute calls fn and, on error, retries up to MaxRetries times with
// exponential backoff (InitialDelay * 2^attempt, capped at MaxDelay) while
// observing ctx for cancellation and the budget for allowances.
func Execute(ctx context.Context, cfg RetryConfig, fn func() Result) Result {
	// TODO: Call fn once. On error, loop attempting retries up to cfg.MaxRetries.
	// For each retry, verify the budget via b.Allow() and check ctx.Done().
	// Sleep for the computed backoff delay before each retry.
	// Return the last Result obtained.
	panic("not implemented")
}
