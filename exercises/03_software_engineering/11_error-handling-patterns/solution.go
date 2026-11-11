package errorhandling

import (
	"errors"
	"fmt"
)

var (
	ErrEmpty    = errors.New("empty value")
	ErrTooShort = errors.New("value too short")
	ErrTooLong  = errors.New("value too long")
	ErrInvalid  = errors.New("invalid value")
)

type ValidationError struct {
	Field string
	Value string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %q: %v", e.Field, e.Err)
}

func (e *ValidationError) Unwrap() error { return e.Err }

func Validate(username string) error {
	return validateField("username", username)
}

func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Value: email, Err: ErrEmpty}
	}
	if len(email) < 5 || len(email) > 254 {
		return &ValidationError{Field: "email", Value: email, Err: ErrTooShort}
	}
	return nil
}

func validateField(field, value string) error {
	if value == "" {
		return &ValidationError{Field: field, Value: value, Err: ErrEmpty}
	}
	if len(value) < 3 {
		return &ValidationError{Field: field, Value: value, Err: ErrTooShort}
	}
	if len(value) > 128 {
		return &ValidationError{Field: field, Value: value, Err: ErrTooLong}
	}
	return nil
}

func wrapValidate(field string, fn func(string) error) func(string) error {
	return func(v string) error {
		if err := fn(v); err != nil {
			return fmt.Errorf("validating %s: %w", field, err)
		}
		return nil
	}
}
