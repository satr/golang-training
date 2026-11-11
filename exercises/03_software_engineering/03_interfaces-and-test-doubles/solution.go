package interfacestestdoubles

import (
	"context"
	"errors"
	"sync"
)

// User represents a domain entity.
type User struct {
	ID    string
	Name  string
	Email string
}

// ErrNotFound is returned when a user is not present in the store.
var ErrNotFound = errors.New("user not found")

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
	return &InMemoryRepo{users: make(map[string]User)}
}

// Find returns the user with the given id.
func (r *InMemoryRepo) Find(id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

// Save persists the given user.
func (r *InMemoryRepo) Save(u User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

// UserHandler delegates user operations to its repository.
type UserHandler struct {
	repo UserRepo
}

// NewUserHandler creates a UserHandler backed by repo.
func NewUserHandler(repo UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUser returns the user with the given id. Propagates context cancellation.
func (h *UserHandler) GetUser(ctx context.Context, id string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}
	return h.repo.Find(id)
}

// UpdateUser persists the given user. Propagates context cancellation.
func (h *UserHandler) UpdateUser(ctx context.Context, u User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return h.repo.Save(u)
}
