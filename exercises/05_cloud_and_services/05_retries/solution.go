package retries

import (
	"context"
	"errors"
	"sync"
	"time"
)

// RetryConfig holds backoff and budget parameters for a retry attempt.
type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Budget       *Budget
}

// Result is what the retried function returns.
type Result struct {
	Value []byte
	Err   error
}

// Budget enforces a system-wide cap on retry attempts within a sliding window.
type Budget struct {
	maxRetries int
	window     time.Duration
	attempts   []time.Time
	mu         sync.Mutex
}

// NewBudget creates a budget allowing up to max retries within the given window.
func NewBudget(max int, window time.Duration) *Budget {
	return &Budget{
		maxRetries: max,
		window:     window,
	}
}

// Allow reports whether a retry is permitted under the current budget.
func (b *Budget) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-b.window)

	valid := make([]time.Time, 0, len(b.attempts))
	for _, t := range b.attempts {
		if !t.Before(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) < b.maxRetries {
		b.attempts = append(valid, now)
		return true
	}
	b.attempts = valid
	return false
}

// Execute calls fn and, on error, retries up to MaxRetries times with
// exponential backoff (InitialDelay * 2^attempt, capped at MaxDelay) while
// observing ctx for cancellation and the budget for allowances.
func Execute(ctx context.Context, cfg RetryConfig, fn func() Result) Result {
	res := fn()
	if res.Err == nil {
		return res
	}

	for attempt := 0; attempt < cfg.MaxRetries; attempt++ {
		if cfg.Budget != nil && !cfg.Budget.Allow() {
			return res
		}

		delay := cfg.InitialDelay * (1 << uint(attempt))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		select {
		case <-ctx.Done():
			return res
		case <-time.After(delay):
		}

		res = fn()
		if res.Err == nil {
			return res
		}
	}
	return res
}

// ErrNotFound is a sentinel used in tests.
var ErrNotFound = errors.New("not found")
