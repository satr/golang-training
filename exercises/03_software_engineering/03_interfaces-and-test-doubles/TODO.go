//go:build todo

package interfacestestdoubles

import (
	"context"
	"sync"
)

// User represents a domain entity.
type User struct {
	ID    string
	Name  string
	Email string
}

// ErrNotFound is returned when a user is not present in the store.
var ErrNotFound error = nil // TODO: define a real sentinel error

// UserRepo abstracts user persistence.
type UserRepo interface {
	Find(id string) (User, error)
	Save(u User) error
}

// InMemoryRepo is an in-memory implementation of UserRepo.
type InMemoryRepo struct {
	users map[string]User
	mu    sync.RWMutex
}

// NewInMemoryRepo creates a ready-to-use InMemoryRepo.
func NewInMemoryRepo() *InMemoryRepo {
	// TODO: Initialize and return a repo with an allocated map.
	return &InMemoryRepo{}
}

// Find returns the user with the given id.
func (r *InMemoryRepo) Find(id string) (User, error) {
	// TODO: Implement thread-safe lookup; return ErrNotFound for missing ids.
	return User{}, nil
}

// Save persists the given user.
func (r *InMemoryRepo) Save(u User) error {
	// TODO: Implement thread-safe save.
	return nil
}

// UserHandler delegates user operations to its repository.
type UserHandler struct {
	repo UserRepo
}

// NewUserHandler creates a UserHandler backed by repo.
func NewUserHandler(repo UserRepo) *UserHandler {
	// TODO: Return a configured UserHandler.
	return &UserHandler{}
}

// GetUser returns the user with the given id. Propagates context cancellation.
func (h *UserHandler) GetUser(ctx context.Context, id string) (User, error) {
	// TODO: Check ctx cancellation, then delegate to h.repo.Find.
	return User{}, nil
}

// UpdateUser persists the given user. Propagates context cancellation.
func (h *UserHandler) UpdateUser(ctx context.Context, u User) error {
	// TODO: Check ctx cancellation, then delegate to h.repo.Save.
	return nil
}
