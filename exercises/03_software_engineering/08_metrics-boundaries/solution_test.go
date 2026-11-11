package metricsboundaries

import (
	"math"
	"sync"
	"testing"
	"time"
)

func TestCounter_IncAndErr(t *testing.T) {
	c := &Counter{}
	c.Inc()
	c.Inc()
	c.Err()
	total, errs := c.Snapshot()
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if errs != 1 {
		t.Errorf("errors = %d, want 1", errs)
	}
}

func TestCounter_Rate(t *testing.T) {
	c := &Counter{}
	c.Inc()
	c.Inc()
	c.Err()
	rate := c.Rate()
	if !almostEqual(rate, 0.333333) {
		t.Errorf("rate = %f, want ~0.333", rate)
	}
}

func TestCounter_Concurrent(t *testing.T) {
	c := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); c.Inc() }()
		go func() { defer wg.Done(); c.Err() }()
	}
	wg.Wait()
	total, errs := c.Snapshot()
	if total != 200 {
		t.Errorf("total = %d, want 200", total)
	}
	if errs != 100 {
		t.Errorf("errors = %d, want 100", errs)
	}
}

func TestCircuitBreaker_ClosedAllows(t *testing.T) {
	cb := NewCircuitBreaker(3, 50*time.Millisecond)
	if !cb.Allow() {
		t.Error("expected Allow on closed breaker")
	}
}

func TestCircuitBreaker_TripsOpen(t *testing.T) {
	cb := NewCircuitBreaker(3, 50*time.Millisecond)
	for i := 0; i < 3; i++ {
		cb.Fail()
	}
	if cb.State() != Open {
		t.Errorf("expected Open, got %s", cb.State())
	}
	if cb.Allow() {
		t.Error("expected no allow when open")
	}
}

func TestCircuitBreaker_HalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)
	cb.Fail()
	cb.Fail()
	time.Sleep(60 * time.Millisecond)
	if cb.State() != HalfOpen {
		t.Errorf("expected HalfOpen, got %s", cb.State())
	}
	if !cb.Allow() {
		t.Error("expected Allow on half-open")
	}
}

func TestCircuitBreaker_Recovery(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)
	cb.Fail()
	cb.Fail()
	time.Sleep(60 * time.Millisecond)
	cb.Success()
	if cb.State() != Closed {
		t.Errorf("expected Closed after recovery, got %s", cb.State())
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}
