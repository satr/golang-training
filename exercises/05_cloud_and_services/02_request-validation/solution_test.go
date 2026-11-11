package requestvalidation

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateEmptyRules(t *testing.T) {
	req := Request{Email: "a@example.com", Name: "A", Quantity: "1"}
	got := Validate(req, nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestValidateEmptyRulesSlice(t *testing.T) {
	req := Request{}
	got := Validate(req, []Rule{})
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestValidateSingleError(t *testing.T) {
	req := Request{Name: ""}
	rules := []Rule{
		{Field: "Name", Check: func(v string) string {
			if v == "" {
				return "required"
			}
			return ""
		}},
	}
	got := Validate(req, rules)
	if got == nil {
		t.Fatal("expected ValidationError")
	}
	if len(got.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(got.Errors))
	}
	if got.Errors[0].Field != "Name" || got.Errors[0].Issue != "required" {
		t.Fatalf("unexpected error: %v", got.Errors[0])
	}
}

func TestValidateMultipleErrorsAccumulate(t *testing.T) {
	req := Request{Email: "", Name: ""}
	rules := []Rule{
		{Field: "Email", Check: func(v string) string {
			if v == "" {
				return "required"
			}
			return ""
		}},
		{Field: "Name", Check: func(v string) string {
			if v == "" {
				return "required"
			}
			return ""
		}},
	}
	got := Validate(req, rules)
	if got == nil {
		t.Fatal("expected ValidationError")
	}
	if len(got.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(got.Errors))
	}
}

func TestValidateAllPass(t *testing.T) {
	rules := []Rule{
		{Field: "Email", Check: func(v string) string {
			if !strings.Contains(v, "@") {
				return "invalid email"
			}
			return ""
		}},
		{Field: "Quantity", Check: func(v string) string {
			if v == "0" {
				return "must be > 0"
			}
			return ""
		}},
	}
	req := Request{Email: "u@example.com", Quantity: "5"}
	got := Validate(req, rules)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFieldErrorErrorMethod(t *testing.T) {
	fe := FieldError{Field: "email", Issue: "invalid"}
	want := "email: invalid"
	if got := fe.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestValidationErrorErrorMethod(t *testing.T) {
	ve := &ValidationError{Errors: []FieldError{
		{Field: "a", Issue: "bad"},
		{Field: "b", Issue: "missing"},
	}}
	want := "a: bad; b: missing"
	if got := ve.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestValidationErrorIsEmpty(t *testing.T) {
	ve := &ValidationError{}
	tests := []struct {
		name string
		e    error
		want bool
	}{
		{name: "same type", e: &ValidationError{}, want: true},
		{name: "other type", e: errors.New("no"), want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ve.Is(test.e); got != test.want {
				t.Fatalf("Is(%T) = %v, want %v", test.e, got, test.want)
			}
		})
	}
}
