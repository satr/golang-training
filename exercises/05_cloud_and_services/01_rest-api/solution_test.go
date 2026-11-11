package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// memStore is an in-memory Store used for testing.
type memStore struct {
	items map[string]Item
	order []string
}

func (ms *memStore) Get(_ context.Context, id string) (Item, error) {
	if item, ok := ms.items[id]; ok {
		return item, nil
	}
	return Item{}, &NotFoundError{ID: id}
}

func (ms *memStore) List(_ context.Context) ([]Item, error) {
	res := make([]Item, 0, len(ms.order))
	for _, id := range ms.order {
		res = append(res, ms.items[id])
	}
	return res, nil
}

func (ms *memStore) Create(_ context.Context, item Item) (Item, error) {
	if ms.items == nil {
		ms.items = make(map[string]Item)
	}
	ms.items[item.ID] = item
	ms.order = append(ms.order, item.ID)
	return item, nil
}

func newTestStore(items []Item) *memStore {
	ms := &memStore{items: make(map[string]Item)}
	for _, it := range items {
		ms.items[it.ID] = it
		ms.order = append(ms.order, it.ID)
	}
	return ms
}

func TestNewHandler(t *testing.T) {
	h := NewHandler(nil)
	if h == nil {
		t.Fatal("NewHandler returned nil")
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		name     string
		store    *memStore
		method   string
		path     string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "GET /items returns empty list",
			store:    newTestStore(nil),
			method:   http.MethodGet,
			path:     "/items",
			wantCode: http.StatusOK,
			wantBody: "null\n",
		},
		{
			name:     "GET /items returns items",
			store:    newTestStore([]Item{{ID: "1", Name: "a"}, {ID: "2", Name: "b"}}),
			method:   http.MethodGet,
			path:     "/items",
			wantCode: http.StatusOK,
		},
		{
			name:     "GET /items/1 returns item",
			store:    newTestStore([]Item{{ID: "1", Name: "a"}}),
			method:   http.MethodGet,
			path:     "/items/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "GET /items/99 returns 404",
			store:    newTestStore([]Item{}),
			method:   http.MethodGet,
			path:     "/items/99",
			wantCode: http.StatusNotFound,
			wantBody: "not found: 99\n",
		},
		{
			name:     "POST /items creates item",
			store:    newTestStore(nil),
			method:   http.MethodPost,
			path:     "/items",
			body:     `{"name":"widget"}`,
			wantCode: http.StatusCreated,
		},
		{
			name:     "DELETE returns 405",
			store:    newTestStore(nil),
			method:   http.MethodDelete,
			path:     "/items",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "GET /other returns 405",
			store:    newTestStore(nil),
			method:   http.MethodGet,
			path:     "/other",
			wantCode: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHandler(test.store)
			req := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != test.wantCode {
				t.Fatalf("code = %d, want %d", rec.Code, test.wantCode)
			}
			if test.wantBody != "" && rec.Body.String() != test.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), test.wantBody)
			}
		})
	}
}

func TestNotFoundErrorMessage(t *testing.T) {
	err := &NotFoundError{ID: "42"}
	if got := err.Error(); got != "not found: 42" {
		t.Fatalf("Error() = %q, want %q", got, "not found: 42")
	}
}

func TestCreateSetsID(t *testing.T) {
	store := newTestStore(nil)
	h := NewHandler(store)
	for _, name := range []string{"a", "b"} {
		req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"name":"`+name+`"}`))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("POST code = %d, want %d", rec.Code, http.StatusCreated)
		}
		var item Item
		json.NewDecoder(rec.Body).Decode(&item)
		if item.ID == "" {
			t.Error("created item has empty ID")
		}
	}
}
