package ratelimiting

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLimiterAllowOne(t *testing.T) {
	l := NewLimiter(1)
	if !l.Allow() {
		t.Error("expected first Allow() to succeed")
	}
	if l.Allow() {
		t.Error("expected second Allow() to fail")
	}
}

func TestLimiterAllowMultiple(t *testing.T) {
	l := NewLimiter(3)
	results := []bool{l.Allow(), l.Allow(), l.Allow(), l.Allow()}
	expected := []bool{true, true, true, false}
	for i := range expected {
		if results[i] != expected[i] {
			t.Errorf("Allow() #%d = %v; want %v", i+1, results[i], expected[i])
		}
	}
}

func TestLimiterAllowNSuccess(t *testing.T) {
	l := NewLimiter(5)
	if !l.AllowN(3) {
		t.Error("expected AllowN(3) to succeed with 5 tokens")
	}
	if l.tokens != 2 {
		t.Errorf("tokens = %d; want 2", l.tokens)
	}
}

func TestLimiterAllowNFail(t *testing.T) {
	l := NewLimiter(3)
	if l.AllowN(5) {
		t.Error("expected AllowN(5) to fail with 3 tokens")
	}
	if l.tokens != 3 {
		t.Errorf("tokens = %d; want 3 (unchanged on failure)", l.tokens)
	}
}

func TestMiddlewareAllowsWithinLimit(t *testing.T) {
	callCount := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Write([]byte("ok"))
	})

	m := &Middleware{limiter: NewLimiter(2)}
	handler := m.ServeHTTP(next)

	for i := 0; i < 3; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(rec, req)

		if i < 2 {
			if rec.Code != http.StatusOK {
				t.Errorf("request %d: status = %d; want %d", i+1, rec.Code, http.StatusOK)
			}
		} else {
			if rec.Code != http.StatusTooManyRequests {
				t.Errorf("request %d: status = %d; want %d", i+1, rec.Code, http.StatusTooManyRequests)
			}
			if !strings.Contains(rec.Body.String(), "rate limited") {
				t.Errorf("body = %q; want to contain %q", rec.Body.String(), "rate limited")
			}
		}
	}

	if callCount != 2 {
		t.Errorf("callCount = %d; want 2", callCount)
	}
}

func TestMiddlewareStatus429(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	m := &Middleware{limiter: NewLimiter(0)}
	handler := m.ServeHTTP(next)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("status = %d; want %d", rec.Code, http.StatusTooManyRequests)
	}
}
