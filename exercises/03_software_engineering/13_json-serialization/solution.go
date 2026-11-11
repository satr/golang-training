package jsonserialization

import (
	"encoding/json"
	"fmt"
	"strings"
)

// User represents a user that can be serialized to JSON with a masked email.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// MarshalJSON returns a JSON encoding of User with the email address masked.
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Alias
		Email string `json:"email"`
	}{
		Alias: (Alias)(u),
		Email: maskEmail(u.Email),
	})
}

// UnmarshalJSON parses JSON into User, validating required fields and setting a default role.
func (u *User) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("decode user: %w", err)
	}
	if raw.ID == "" {
		return fmt.Errorf("user id is required")
	}
	if raw.Name == "" {
		return fmt.Errorf("user name is required")
	}
	u.ID = raw.ID
	u.Name = raw.Name
	u.Email = raw.Email
	if raw.Role == "" {
		u.Role = "user"
	} else {
		u.Role = raw.Role
	}
	return nil
}

// maskEmail redacts the local part of an email, keeping only the first two characters.
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	if len(parts[0]) <= 2 {
		return "xx@" + parts[1]
	}
	return parts[0][:2] + "***@" + parts[1]
}
