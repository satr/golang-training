package packagearchitecture

import (
	"fmt"
	"sync"
	"testing"
)

func TestStore_SetAndGet(t *testing.T) {
	store := NewInMemoryStore()
	err := store.Set("name", "go")
	if err != nil {
		t.Fatal(err)
	}
	v, err := store.Get("name")
	if err != nil {
		t.Fatal(err)
	}
	if v != "go" {
		t.Errorf("got %q, want go", v)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	store := NewInMemoryStore()
	_, err := store.Get("nope")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestStore_Delete(t *testing.T) {
	store := NewInMemoryStore()
	store.Set("key", "val")
	if err := store.Delete("key"); err != nil {
		t.Fatal(err)
	}
	_, err := store.Get("key")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestStore_DeleteNotFound(t *testing.T) {
	store := NewInMemoryStore()
	if err := store.Delete("nope"); err == nil {
		t.Error("expected error deleting missing key")
	}
}

func TestStore_Concurrent(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		key := fmt.Sprintf("k%d", i)
		go func() { defer wg.Done(); store.Set(key, "v") }()
		go func() { defer wg.Done(); store.Get(key) }()
	}
	wg.Wait()
}
