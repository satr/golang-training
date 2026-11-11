package racedetection

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const testTimeout = 5 * time.Second

func TestCacheZeroValueAndSuccessfulHit(t *testing.T) {
	var cache Cache
	var calls atomic.Int32
	load := func(key string) (string, error) {
		calls.Add(1)
		return "value:" + key, nil
	}

	for range 2 {
		got, err := cache.Get("alpha", load)
		if err != nil {
			t.Fatalf("Get() error = %v; want nil", err)
		}
		if got != "value:alpha" {
			t.Fatalf("Get() = %q; want %q", got, "value:alpha")
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("loader calls = %d; want 1", got)
	}
	if got := cache.Len(); got != 1 {
		t.Fatalf("Len() = %d; want 1", got)
	}
}

func TestCacheValidatesArguments(t *testing.T) {
	var cache Cache
	loaderCalled := false
	load := func(string) (string, error) {
		loaderCalled = true
		return "", nil
	}

	if _, err := cache.Get("", load); !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("Get(empty key) error = %v; want ErrEmptyKey", err)
	}
	if _, err := cache.Get("alpha", nil); !errors.Is(err, ErrNilLoader) {
		t.Fatalf("Get(nil loader) error = %v; want ErrNilLoader", err)
	}
	if loaderCalled {
		t.Fatal("loader was called for invalid input")
	}
}

func TestConcurrentSameKeyUsesOneLoad(t *testing.T) {
	var cache Cache
	const callers = 32
	var calls atomic.Int32
	loaderStarted := make(chan struct{})
	releaseLoader := make(chan struct{})
	var startedOnce sync.Once

	load := func(key string) (string, error) {
		calls.Add(1)
		startedOnce.Do(func() { close(loaderStarted) })
		<-releaseLoader
		return "shared:" + key, nil
	}

	start := make(chan struct{})
	results := make(chan string, callers)
	errs := make(chan error, callers)
	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			<-start
			value, err := cache.Get("alpha", load)
			results <- value
			errs <- err
		}()
	}

	close(start)
	waitClosed(t, loaderStarted, "loader did not start")
	close(releaseLoader)
	wg.Wait()
	close(results)
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("Get() error = %v; want nil", err)
		}
	}
	for value := range results {
		if value != "shared:alpha" {
			t.Fatalf("Get() = %q; want %q", value, "shared:alpha")
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("loader calls = %d; want 1", got)
	}
}

func TestDifferentKeysLoadIndependently(t *testing.T) {
	var cache Cache
	started := make(chan string, 2)
	release := make(chan struct{})
	load := func(key string) (string, error) {
		started <- key
		<-release
		return key, nil
	}

	type result struct {
		value string
		err   error
	}
	done := make(chan result, 2)
	for _, key := range []string{"alpha", "beta"} {
		go func() {
			value, err := cache.Get(key, load)
			done <- result{value: value, err: err}
		}()
	}

	seen := make(map[string]bool)
	for range 2 {
		select {
		case key := <-started:
			seen[key] = true
		case <-time.After(testTimeout):
			t.Fatal("different-key loader was serialized or did not start")
		}
	}
	close(release)

	for range 2 {
		select {
		case got := <-done:
			if got.err != nil {
				t.Fatalf("Get() error = %v; want nil", got.err)
			}
			if !seen[got.value] {
				t.Fatalf("Get() value = %q; want alpha or beta", got.value)
			}
		case <-time.After(testTimeout):
			t.Fatal("Get() did not finish after loader release")
		}
	}
}

func TestFailedLoadIsNotCached(t *testing.T) {
	var cache Cache
	wantErr := errors.New("backend unavailable")
	var calls atomic.Int32
	load := func(string) (string, error) {
		if calls.Add(1) == 1 {
			return "", wantErr
		}
		return "recovered", nil
	}

	if value, err := cache.Get("alpha", load); !errors.Is(err, wantErr) || value != "" {
		t.Fatalf("first Get() = (%q, %v); want empty value and wrapped loader error", value, err)
	}
	if got := cache.Len(); got != 0 {
		t.Fatalf("Len() after failure = %d; want 0", got)
	}

	value, err := cache.Get("alpha", load)
	if err != nil {
		t.Fatalf("retry Get() error = %v; want nil", err)
	}
	if value != "recovered" {
		t.Fatalf("retry Get() = %q; want %q", value, "recovered")
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("loader calls = %d; want 2", got)
	}
}

func TestDeleteDoesNotRemoveInFlightLoad(t *testing.T) {
	var cache Cache
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		_, err := cache.Get("alpha", func(string) (string, error) {
			close(started)
			<-release
			return "value", nil
		})
		done <- err
	}()

	waitClosed(t, started, "loader did not start")
	if cache.Delete("alpha") {
		t.Fatal("Delete() = true for in-flight load; want false")
	}
	close(release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Get() error = %v; want nil", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("Get() did not finish")
	}
	if !cache.Delete("alpha") {
		t.Fatal("Delete() = false for completed value; want true")
	}
	if got := cache.Len(); got != 0 {
		t.Fatalf("Len() = %d; want 0", got)
	}
}

func TestConcurrentMixedOperations(t *testing.T) {
	var cache Cache
	const goroutines = 24
	const iterations = 200
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for worker := range goroutines {
		go func() {
			defer wg.Done()
			<-start
			for iteration := range iterations {
				key := fmt.Sprintf("key-%d", (worker+iteration)%4)
				value, err := cache.Get(key, func(key string) (string, error) {
					return "value:" + key, nil
				})
				if err != nil || value != "value:"+key {
					t.Errorf("Get(%q) = (%q, %v)", key, value, err)
					return
				}
				if iteration%17 == 0 {
					cache.Delete(key)
				}
				_ = cache.Len()
			}
		}()
	}

	close(start)
	wg.Wait()
}

func waitClosed(t *testing.T, channel <-chan struct{}, message string) {
	t.Helper()
	select {
	case <-channel:
	case <-time.After(testTimeout):
		t.Fatal(message)
	}
}
