package jsonserialization

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	u := User{ID: "1", Name: "Alice", Email: "alice@example.com", Role: "admin"}
	data, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	out := string(data)
	if !strings.Contains(out, "***") {
		t.Errorf("expected masked email in %s", out)
	}
}

func TestUnmarshal(t *testing.T) {
	data := []byte(`{"id":"1","name":"Bob","email":"bob@example.com"}`)
	var u User
	if err := json.Unmarshal(data, &u); err != nil {
		t.Fatal(err)
	}
	if u.ID != "1" || u.Name != "Bob" || u.Role != "user" {
		t.Errorf("got %+v", u)
	}
}

func TestUnmarshal_Error_NoID(t *testing.T) {
	data := []byte(`{"name":"No ID"}`)
	var u User
	err := json.Unmarshal(data, &u)
	if err == nil {
		t.Error("expected error for missing id")
	}
}

func TestUnmarshal_Error_InvalidJSON(t *testing.T) {
	data := []byte(`{bad`)
	var u User
	err := json.Unmarshal(data, &u)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestMarshalEmail_Masked(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{"alice@example.com", "al***@example.com"},
		{"ab@example.com", "xx@example.com"},
		{"nobody", "nobody"},
	}
	for _, tt := range tests {
		if got := maskEmail(tt.email); got != tt.want {
			t.Errorf("email %q: got %q, want %q", tt.email, got, tt.want)
		}
	}
}
