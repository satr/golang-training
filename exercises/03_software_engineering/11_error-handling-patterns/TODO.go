//go:build todo

package errorhandling

import "errors"

var (
	ErrEmpty    = errors.New("empty value")
	ErrTooShort = errors.New("value too short")
	ErrTooLong  = errors.New("value too long")
	ErrInvalid  = errors.New("invalid value")
)

type ValidationError struct{}

func (e *ValidationError) Error() string { return "" }
func (e *ValidationError) Unwrap() error { return nil }

func Validate(username string) error          { return nil }
func ValidateEmail(email string) error        { return nil }
func validateField(field, value string) error { return nil }
func wrapValidate(field string, fn func(string) error) func(string) error {
	return func(string) error { return nil }
}
