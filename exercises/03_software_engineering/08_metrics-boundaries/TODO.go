//go:build todo

package metricsboundaries

import "time"

type Counter struct{}

func (c *Counter) Inc()                     {}
func (c *Counter) Err()                     {}
func (c *Counter) Rate() float64            { return 0 }
func (c *Counter) Snapshot() (int64, int64) { return 0, 0 }

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct{}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{}
}

func (cb *CircuitBreaker) State() State { return Closed }
func (cb *CircuitBreaker) Allow() bool  { return true }
func (cb *CircuitBreaker) Success()     {}
func (cb *CircuitBreaker) Fail()        {}
