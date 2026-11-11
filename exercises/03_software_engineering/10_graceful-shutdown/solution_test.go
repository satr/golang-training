package gracefulshutdown

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	s := NewServer("127.0.0.1:0", nil)
	if s == nil {
		t.Fatal("NewServer returned nil")
	}
	if s.shutdownTimeout == 0 {
		t.Error("zero shutdownTimeout")
	}
}

func TestListenAndServeStopsOnCancel(t *testing.T) {
	var hookRan atomic.Bool

	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	s := NewServer("127.0.0.1:0", ok,
		WithShutdownTimeout(100*time.Millisecond),
		WithHook(func(ctx context.Context) error {
			hookRan.Store(true)
			return nil
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ListenAndServe(ctx)
	}()

	<-s.started

	url := fmt.Sprintf("http://%s", s.srv.Addr)
	res, err := http.Get(url)
	if err == nil {
		res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want 200", res.StatusCode)
		}
	}

	cancel()
	err = <-errCh
	if err == nil {
		t.Fatalf("expected non-nil error from ListenAndServe")
	}
	if !hookRan.Load() {
		t.Error("hook did not run")
	}
}

func TestWithShutdownTimeout(t *testing.T) {
	d := 30 * time.Second
	s := NewServer("127.0.0.1:0", nil)
	if s.shutdownTimeout != 5*time.Second {
		t.Error("not default")
	}
	opt := WithShutdownTimeout(d)
	opt(s)
	if s.shutdownTimeout != d {
		t.Errorf("shutdownTimeout = %v, want %v", s.shutdownTimeout, d)
	}
}

func TestHooksRunInOrder(t *testing.T) {
	var order []int
	s := NewServer("127.0.0.1:0", nil,
		WithShutdownTimeout(100*time.Millisecond),
		WithHook(func(ctx context.Context) error {
			order = append(order, 1)
			return nil
		}),
		WithHook(func(ctx context.Context) error {
			order = append(order, 2)
			return nil
		}),
		WithHook(func(ctx context.Context) error {
			order = append(order, 3)
			return nil
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ListenAndServe(ctx)
	}()

	<-s.started
	cancel()
	<-errCh

	want := []int{1, 2, 3}
	if len(order) != len(want) {
		t.Fatalf("hooks order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("hooks order = %v, want %v", order, want)
		}
	}
}
