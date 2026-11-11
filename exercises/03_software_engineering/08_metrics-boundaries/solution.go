package metricsboundaries

import (
	"sync"
	"time"
)

type Counter struct {
	mu     sync.Mutex
	total  int64
	errors int64
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.total++
	c.mu.Unlock()
}

func (c *Counter) Err() {
	c.mu.Lock()
	c.total++
	c.errors++
	c.mu.Unlock()
}

func (c *Counter) Rate() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.total == 0 {
		return 0
	}
	return float64(c.errors) / float64(c.total)
}

func (c *Counter) Snapshot() (total, errors int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.total, c.errors
}

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

func (s State) String() string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failures         int
	failureThreshold int
	recoveryTimeout  time.Duration
	openedAt         time.Time
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            Closed,
		failureThreshold: threshold,
		recoveryTimeout:  timeout,
	}
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == Open && time.Since(cb.openedAt) > cb.recoveryTimeout {
		cb.state = HalfOpen
	}
	return cb.state
}

func (cb *CircuitBreaker) Allow() bool {
	switch cb.State() {
	case Closed:
		return true
	case HalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == Open && time.Since(cb.openedAt) > cb.recoveryTimeout {
		cb.state = Closed
		cb.failures = 0
		return
	}
	if cb.state == HalfOpen {
		cb.state = Closed
		cb.failures = 0
	}
}

func (cb *CircuitBreaker) Fail() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.state == HalfOpen {
		cb.state = Open
		cb.openedAt = time.Now()
	} else if cb.failures >= cb.failureThreshold {
		cb.state = Open
		cb.openedAt = time.Now()
	}
}
