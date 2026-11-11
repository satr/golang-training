package dependencyinjection

import (
	"errors"
	"strings"
	"testing"
)

func TestNewServiceNil(t *testing.T) {
	svc := NewService(nil)
	if svc != nil {
		t.Fatal("NewService(nil) should return nil")
	}
}

func TestStoreGetSetDelete(t *testing.T) {
	store := NewMapStore()

	if err := store.Set("a", "hello"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err := store.Get("a")
	if err != nil {
		t.Fatalf("Get existing: %v", err)
	}
	if got != "hello" {
		t.Fatalf("Get existing: got %q, want %q", got, "hello")
	}

	_, err = store.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get missing: expected ErrNotFound, got %v", err)
	}

	if err := store.Delete("a"); err != nil {
		t.Fatalf("Delete existing: %v", err)
	}

	_, err = store.Get("a")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get deleted: expected ErrNotFound, got %v", err)
	}

	if err := store.Delete("a"); err == nil {
		t.Fatal("Delete missing: expected error, got nil")
	}
}

func TestServiceProcess(t *testing.T) {
	store := NewMapStore()
	svc := NewService(store)
	if svc == nil {
		t.Fatal("NewService with valid store returned nil")
	}

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{name: "uppercase", key: "greeting", value: "hello", wantErr: false},
		{name: "already upper", key: "code", value: "GO123", wantErr: false},
		{name: "new key", key: "newkey", value: "world", wantErr: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.Process(tc.key, tc.value)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, err := svc.Get(tc.key)
			if err != nil {
				t.Fatalf("Get after Process: %v", err)
			}
			want := strings.ToUpper(tc.value)
			if got != want {
				t.Fatalf("Get: got %q, want %q", got, want)
			}
		})
	}
}

func TestServiceGetNotFound(t *testing.T) {
	store := NewMapStore()
	svc := NewService(store)
	_, err := svc.Get("nope")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceDelete(t *testing.T) {
	store := NewMapStore()
	svc := NewService(store)

	if err := store.Set("k", "v"); err != nil {
		t.Fatalf("setup Set: %v", err)
	}

	if err := svc.Delete("k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := svc.Get("k")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get after delete: expected ErrNotFound, got %v", err)
	}
}
