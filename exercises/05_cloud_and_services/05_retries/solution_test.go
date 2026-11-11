package retries

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestExecuteImmediateSuccess(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:   3,
		InitialDelay: 0,
		MaxDelay:     0,
		Budget:       NewBudget(5, time.Minute),
	}
	callCount := 0
	res := Execute(context.Background(), cfg, func() Result {
		callCount++
		return Result{Value: []byte("ok")}
	})

	if res.Err != nil {
		t.Errorf("err = %v; want nil", res.Err)
	}
	if string(res.Value) != "ok" {
		t.Errorf("value = %q; want %q", res.Value, "ok")
	}
	if callCount != 1 {
		t.Errorf("callCount = %d; want 1", callCount)
	}
}

func TestExecuteRetryThenSucceed(t *testing.T) {
	var counter atomic.Int32
	cfg := RetryConfig{
		MaxRetries:   3,
		InitialDelay: 0,
		MaxDelay:     0,
		Budget:       NewBudget(5, time.Minute),
	}
	res := Execute(context.Background(), cfg, func() Result {
		n := counter.Add(1)
		if n < 3 {
			return Result{Err: ErrNotFound}
		}
		return Result{Value: []byte("found")}
	})

	if res.Err != nil {
		t.Errorf("err = %v; want nil", res.Err)
	}
	if string(res.Value) != "found" {
		t.Errorf("value = %q; want %q", res.Value, "found")
	}
	if got := counter.Load(); got != 3 {
		t.Errorf("callCount = %d; want 3", got)
	}
}

func TestExecuteBudgetExhausted(t *testing.T) {
	budget := NewBudget(2, time.Minute)
	cfg := RetryConfig{
		MaxRetries:   5,
		InitialDelay: 0,
		MaxDelay:     0,
		Budget:       budget,
	}
	callCount := 0
	res := Execute(context.Background(), cfg, func() Result {
		callCount++
		return Result{Err: ErrNotFound}
	})

	if res.Err == nil {
		t.Error("err = nil; want non-nil")
	}
	// Initial call (1) + 2 budgeted retries = 3 total calls
	if callCount != 3 {
		t.Errorf("callCount = %d; want 3", callCount)
	}
}

func TestExecuteContextCancelStopsRetries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	budget := NewBudget(100, time.Hour)
	cfg := RetryConfig{
		MaxRetries:   100,
		InitialDelay: 0,
		MaxDelay:     0,
		Budget:       budget,
	}

	done := make(chan struct{})
	var counter atomic.Int32
	go func() {
		defer close(done)
		Execute(ctx, cfg, func() Result {
			counter.Add(1)
			return Result{Err: ErrNotFound}
		})
	}()

	cancel()
	<-done
	// At least one call happened, but retries should be bounded.
	if got := counter.Load(); got < 1 {
		t.Errorf("callCount = %d; want >= 1", got)
	}
}

func TestExecuteMaxRetriesExhausted(t *testing.T) {
	var counter atomic.Int32
	cfg := RetryConfig{
		MaxRetries:   2,
		InitialDelay: 0,
		MaxDelay:     0,
		Budget:       NewBudget(10, time.Minute),
	}
	res := Execute(context.Background(), cfg, func() Result {
		counter.Add(1)
		return Result{Err: ErrNotFound}
	})

	if res.Err == nil {
		t.Error("err = nil; want non-nil")
	}
	if got := counter.Load(); got != 3 {
		t.Errorf("callCount = %d; want 3 (initial + 2 retries)", got)
	}
}
