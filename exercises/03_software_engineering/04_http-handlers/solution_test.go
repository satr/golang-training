package httphandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetStore() { store = NewStore() }

func TestCreateItem(t *testing.T) {
	t.Cleanup(resetStore)
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"valid item", http.MethodPost, `{"name":"widget","quantity":5}`, http.StatusCreated},
		{"missing name", http.MethodPost, `{"quantity":5}`, http.StatusBadRequest},
		{"bad JSON", http.MethodPost, `{invalid`, http.StatusBadRequest},
		{"wrong method", http.MethodGet, "", http.StatusMethodNotAllowed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var bodyReader *strings.Reader
			if tc.body != "" {
				bodyReader = strings.NewReader(tc.body)
			}
			req := httptest.NewRequest(tc.method, "/items", bodyReader)
			rec := httptest.NewRecorder()
			CreateItem(rec, req)
			if rec.Code != tc.wantStatus {
				t.Errorf("got %d, want %d; body: %s", rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestGetItems(t *testing.T) {
	t.Cleanup(resetStore)
	// Test empty
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	GetItems(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("got %d, want 200", rec.Code)
	}
	var items []Item
	json.NewDecoder(rec.Body).Decode(&items)
	if len(items) != 0 {
		t.Errorf("expected empty, got %d items", len(items))
	}

	// Add item then get again
	resetStore()
	store.Add(Item{Name: "a", Quantity: 1})
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/items", nil)
	GetItems(rec, req)
	json.NewDecoder(rec.Body).Decode(&items)
	if len(items) != 1 || items[0].Name != "a" {
		t.Errorf("got %v, want [{a 1}]", items)
	}
}

func TestCreateItemJSON(t *testing.T) {
	t.Cleanup(resetStore)
	payload := []byte(`{"name":"bolts","quantity":10}`)
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	CreateItem(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("got %d, want 201", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var item Item
	if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	if item.Name != "bolts" || item.Quantity != 10 {
		t.Errorf("decoded %v, want {bolts 10}", item)
	}
}
