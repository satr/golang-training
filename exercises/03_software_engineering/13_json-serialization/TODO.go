//go:build todo

package jsonserialization

// User represents a user that can be serialized to JSON with a masked email.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// MarshalJSON returns a JSON encoding of User with the email address masked.
func (u User) MarshalJSON() ([]byte, error) {
	// TODO: Mask email in output (first 2 chars + ***)
	return nil, nil
}

// UnmarshalJSON parses JSON into User, validating required fields and setting a default role.
func (u *User) UnmarshalJSON(data []byte) error {
	// TODO: Validate required fields, set default role
	return nil
}
