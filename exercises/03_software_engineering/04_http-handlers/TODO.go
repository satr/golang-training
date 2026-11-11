//go:build todo

package httphandlers

import (
	"net/http"
)

// Item represents a store item.
type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// Store holds items in memory.
type Store struct{}

// NewStore creates a new Store.
func NewStore() *Store {
	// TODO: Implement thread-safe in-memory store.
	return &Store{}
}

// Add adds an item to the store.
func (s *Store) Add(item Item) {
	// TODO: Append item to internal slice with mutex protection.
}

// Items returns a copy of all items.
func (s *Store) Items() []Item {
	// TODO: Return a slice copy with read lock.
	return []Item{}
}

// store is the package-level item store.
var store *Store = NewStore()

// CreateItem handles POST requests to create items.
func CreateItem(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode JSON, validate name is non-empty, store item, return 201.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// GetItems handles GET requests to list all items.
func GetItems(w http.ResponseWriter, r *http.Request) {
	// TODO: Return all items as JSON array.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
