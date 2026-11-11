package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Echo-ID", r.Header.Get("X-Request-ID"))
	})
	h := RequestID(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Header().Get("X-Echo-ID") == "" {
		t.Error("expected request ID")
	}

	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "my-id")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Header().Get("X-Echo-ID") != "my-id" {
		t.Errorf("expected my-id, got %s", rec.Header().Get("X-Echo-ID"))
	}
}

func TestRecovery(t *testing.T) {
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("x") })
	h := Recovery(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want 500", rec.Code)
	}
}

func TestChain(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Echo-ID", r.Header.Get("X-Request-ID"))
	})
	h := Chain(next, RequestID)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Header().Get("X-Echo-ID") == "" {
		t.Error("chain did not apply RequestID")
	}
}
