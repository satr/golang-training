package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

// Item represents a domain entity stored and returned via JSON.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Store abstracts persistence for items.
type Store interface {
	Get(ctx context.Context, id string) (Item, error)
	List(ctx context.Context) ([]Item, error)
	Create(ctx context.Context, item Item) (Item, error)
}

// NotFoundError indicates that a requested item could not be found.
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.ID)
}

// Handler serves HTTP requests against a Store.
type Handler struct {
	store  Store
	nextID uint64
	mu     sync.Mutex
}

// NewHandler returns a Handler that uses the given store.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// ServeHTTP dispatches GET and POST requests to /items.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/items":
		h.handleList(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/items/"):
		h.handleGet(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/items":
		h.handleCreate(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	item, err := h.store.Get(r.Context(), id)
	if err != nil {
		var nfe *NotFoundError
		if errors.As(err, &nfe) {
			http.Error(w, nfe.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.nextID++
	nid := atomic.LoadUint64(&h.nextID)
	item.ID = fmt.Sprintf("item-%d", nid)
	h.mu.Unlock()

	created, err := h.store.Create(r.Context(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
