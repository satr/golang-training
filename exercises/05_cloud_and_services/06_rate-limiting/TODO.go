//go:build todo

package ratelimiting

import (
	"net/http"
	"sync"
)

// Limiter is a token-bucket rate limiter.
type Limiter struct {
	// ...
}

// NewLimiter creates a limiter with the given maximum token count.
func NewLimiter(max int) *Limiter {
	// TODO: Return a *Limiter with tokens set to max.
	return nil
}

// Allow reports whether a single token is available. If so, it consumes the
// token and returns true.
func (l *Limiter) Allow() bool {
	// TODO: Acquire the mutex. If tokens > 0, decrement and return true.
	// Otherwise return false.
	return false
}

// AllowN reports whether n tokens are available atomically. If so, it
// consumes all n tokens and returns true.
func (l *Limiter) AllowN(n int) bool {
	// TODO: Acquire the mutex. If tokens >= n, decrement by n and return true.
	// Otherwise return false without modifying tokens.
	return false
}

// Middleware wraps an http.Handler and enforces rate limiting.
type Middleware struct {
	limiter *Limiter
}

// ServeHTTP returns an http.Handler that checks the rate limiter before
// passing the request to next. Returns 429 when the limiter is exhausted.
func (m *Middleware) ServeHTTP(next http.Handler) http.Handler {
	// TODO: Return an http.Handler that checks m.limiter.Allow() before
	// calling next.ServeHTTP. Write a 429 response body "rate limited\n"
	// when the token is not available.
	panic("not implemented")
}
