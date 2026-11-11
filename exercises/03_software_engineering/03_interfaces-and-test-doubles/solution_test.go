package interfacestestdoubles

import (
	"context"
	"sync"
	"testing"
)

// testDouble is a lightweight test double that satisfies UserRepo.
type testDouble struct {
	mu    sync.RWMutex
	users map[string]User
}

func (d *testDouble) Find(id string) (User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	u, ok := d.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (d *testDouble) Save(u User) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.users[u.ID] = u
	return nil
}

func newTestDouble(users []User) *testDouble {
	m := make(map[string]User, len(users))
	for _, u := range users {
		m[u.ID] = u
	}
	return &testDouble{users: m}
}

func TestGetUserExisting(t *testing.T) {
	td := newTestDouble([]User{{ID: "1", Name: "Ada", Email: "ada@example.com"}})
	h := NewUserHandler(td)

	u, err := h.GetUser(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if u.Name != "Ada" {
		t.Fatalf("Name = %q, want %q", u.Name, "Ada")
	}
	if u.Email != "ada@example.com" {
		t.Fatalf("Email = %q, want %q", u.Email, "ada@example.com")
	}
}

func TestGetUserNotFound(t *testing.T) {
	td := newTestDouble(nil)
	h := NewUserHandler(td)

	_, err := h.GetUser(context.Background(), "99")
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}

func TestUpdateUserExisting(t *testing.T) {
	td := newTestDouble([]User{{ID: "1", Name: "Ada", Email: "ada@example.com"}})
	h := NewUserHandler(td)

	err := h.UpdateUser(context.Background(), User{ID: "1", Name: "Ada Updated", Email: "ada@example.com"})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}

	got, err := td.Find("1")
	if err != nil {
		t.Fatalf("Find after update: %v", err)
	}
	if got.Name != "Ada Updated" {
		t.Fatalf("Name = %q, want %q", got.Name, "Ada Updated")
	}
}

func TestUpdateUserSavesNew(t *testing.T) {
	td := newTestDouble(nil)
	h := NewUserHandler(td)

	newUser := User{ID: "2", Name: "Grace", Email: "grace@example.com"}
	err := h.UpdateUser(context.Background(), newUser)
	if err != nil {
		t.Fatalf("UpdateUser (new): %v", err)
	}

	got, err := td.Find("2")
	if err != nil {
		t.Fatalf("Find new user: %v", err)
	}
	if got.Name != "Grace" {
		t.Fatalf("Name = %q, want %q", got.Name, "Grace")
	}
}

func TestGetUserContextCancelled(t *testing.T) {
	td := newTestDouble(nil)
	h := NewUserHandler(td)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before calling

	_, err := h.GetUser(ctx, "1")
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestUpdateUserContextCancelled(t *testing.T) {
	td := newTestDouble(nil)
	h := NewUserHandler(td)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := h.UpdateUser(ctx, User{ID: "1", Name: "X", Email: "x@x"})
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}
