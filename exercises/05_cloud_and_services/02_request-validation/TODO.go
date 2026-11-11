//go:build todo

package requestvalidation

// FieldError describes a single field validation issue.
type FieldError struct {
	Field string
	Issue string
}

func (e *FieldError) Error() string {
	// TODO: Return "Field: Issue".
	return ""
}

// ValidationError is an aggregate of field errors.
type ValidationError struct {
	Errors []FieldError
}

func (e *ValidationError) Error() string {
	// TODO: Return a joined error message.
	return ""
}

func (e *ValidationError) Is(target error) bool {
	// TODO: Return true if target is *ValidationError.
	return false
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
	// TODO: For each rule, resolve req field by name, run Check,
	// collect FieldErrors, return nil when no errors.
	return nil
}
