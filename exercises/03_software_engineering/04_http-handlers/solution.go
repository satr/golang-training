package httphandlers

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Item represents a store item.
type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// Store holds items in memory.
type Store struct {
	items []Item
	mu    sync.RWMutex
}

// NewStore creates a new Store.
func NewStore() *Store {
	return &Store{}
}

// Add adds an item to the store.
func (s *Store) Add(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Items returns a copy of all items.
func (s *Store) Items() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return []Item{}
	}
	result := make([]Item, len(s.items))
	copy(result, s.items)
	return result
}

// store is the package-level item store.
var store = NewStore()

// CreateItem handles POST requests to create items.
func CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if item.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	store.Add(item)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// GetItems handles GET requests to list all items.
func GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.Items())
}
