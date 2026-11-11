package httpclientstimeouts

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSafeClient_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := NewSafeClient(5 * time.Second)
	ctx := context.Background()
	resp, err := client.Get(ctx, srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: %d", resp.StatusCode)
	}
}

func TestSafeClient_Get_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))
	defer srv.Close()

	client := NewSafeClient(50 * time.Millisecond)
	ctx := context.Background()
	_, err := client.Get(ctx, srv.URL)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "context deadline") && !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Logf("error was: %v", err)
	}
}

func TestSafeClient_Post(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := NewSafeClient(5 * time.Second)
	ctx := context.Background()
	resp, err := client.Post(ctx, srv.URL, "application/json", strings.NewReader(`{"test":true}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if gotBody != `{"test":true}` {
		t.Errorf("body: %q", gotBody)
	}
}

func TestSafeClient_Post_NilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := NewSafeClient(5 * time.Second)
	ctx := context.Background()
	resp, err := client.Post(ctx, srv.URL, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: %d", resp.StatusCode)
	}
}

func TestSafeClient_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer srv.Close()

	client := NewSafeClient(30 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.Get(ctx, srv.URL)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}
