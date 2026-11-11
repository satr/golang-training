package ratelimiting

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

// Limiter is a token-bucket rate limiter.
type Limiter struct {
	maxTokens int
	tokens    int
	mu        sync.Mutex
}

// NewLimiter creates a limiter with the given maximum token count.
func NewLimiter(max int) *Limiter {
	return &Limiter{
		maxTokens: max,
		tokens:    max,
	}
}

// Allow reports whether a single token is available. If so, it consumes the
// token and returns true.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tokens > 0 {
		l.tokens--
		return true
	}
	return false
}

// AllowN reports whether n tokens are available atomically. If so, it
// consumes all n tokens and returns true.
func (l *Limiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tokens >= n {
		l.tokens -= n
		return true
	}
	return false
}

// Middleware wraps an http.Handler and enforces rate limiting.
type Middleware struct {
	limiter *Limiter
}

// ServeHTTP returns an http.Handler that checks the rate limiter before
// passing the request to next. Returns 429 when the limiter is exhausted.
func (m *Middleware) ServeHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Test helper: new middleware with a simple handler.
func newMiddlewareForTest(max int, next http.Handler) *Middleware {
	return &Middleware{limiter: NewLimiter(max)}
}

// AllowAll records whether each request was allowed (for testing).
func newRecordingMiddleware(max int) (*Middleware, *[]bool) {
	var allowed []bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed = append(allowed, true)
	})
	m := newMiddlewareForTest(max, next)
	return m, &allowed
}

func init() {
	_ = newRecordingMiddleware
}

func testNewRecordingMiddleware(max int) (*Middleware, *httptest.ResponseRecorder, chan bool) {
	allowed := make(chan bool, 100)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed <- true
	})
	m := &Middleware{limiter: NewLimiter(max)}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	_ = rec
	_ = req
	_ = m
	_ = next
	_ = allowed
	return m, rec, allowed
}
