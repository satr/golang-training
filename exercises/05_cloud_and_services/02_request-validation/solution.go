package requestvalidation

import (
	"strings"
)

// FieldError describes a single field validation issue.
type FieldError struct {
	Field string
	Issue string
}

func (e *FieldError) Error() string {
	return e.Field + ": " + e.Issue
}

// ValidationError is an aggregate of field errors.
type ValidationError struct {
	Errors []FieldError
}

func (e *ValidationError) Error() string {
	parts := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		parts = append(parts, err.Error())
	}
	return strings.Join(parts, "; ")
}

func (e *ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}

// Rule describes a check for a single request field.
type Rule struct {
	Field string
	Check func(string) string // "" = valid, nonzero = issue message
}

// Request is the input record to validate.
type Request struct {
	Email    string
	Name     string
	Quantity string
}

// Validate runs all rules against the given request.
func Validate(req Request, rules []Rule) *ValidationError {
	var errs []FieldError
	for _, rule := range rules {
		val, ok := fieldByName(rule.Field, req)
		if !ok {
			continue
		}
		msg := rule.Check(val)
		if msg != "" {
			errs = append(errs, FieldError{Field: rule.Field, Issue: msg})
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return &ValidationError{Errors: errs}
}

func fieldByName(name string, req Request) (string, bool) {
	switch name {
	case "Email":
		return req.Email, true
	case "Name":
		return req.Name, true
	case "Quantity":
		return req.Quantity, true
	default:
		return "", false
	}
}
