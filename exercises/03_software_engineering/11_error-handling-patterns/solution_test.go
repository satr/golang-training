package errorhandling

import (
	"errors"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		input string
		want  error
	}{
		{"", ErrEmpty},
		{"ab", ErrTooShort},
		{"alice", nil},
	}
	for _, tt := range tests {
		err := Validate(tt.input)
		if tt.want != nil {
			if !errors.Is(err, tt.want) {
				t.Errorf("input %q: got %v, want Is(%v)", tt.input, err, tt.want)
			}
		} else {
			if err != nil {
				t.Errorf("input %q: unexpected error %v", tt.input, err)
			}
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input string
		want  error
	}{
		{"", ErrEmpty},
		{"a@b", ErrTooShort},
		{"user@example.com", nil},
	}
	for _, tt := range tests {
		err := ValidateEmail(tt.input)
		if tt.want != nil {
			if !errors.Is(err, tt.want) {
				t.Errorf("input %q: got %v, want Is(%v)", tt.input, err, tt.want)
			}
		} else {
			if err != nil {
				t.Errorf("input %q: unexpected error %v", tt.input, err)
			}
		}
	}
}

func TestValidationError_As(t *testing.T) {
	err := Validate("")
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatal("expected ValidationError via errors.As")
	}
	if ve.Field != "username" {
		t.Errorf("field = %q, want username", ve.Field)
	}
}

func TestValidationError_ErrorString(t *testing.T) {
	err := Validate("")
	if !strings.Contains(err.Error(), "username") {
		t.Errorf("error string %q should contain field name", err.Error())
	}
}

func TestWrapValidate(t *testing.T) {
	wrapped := wrapValidate("name", Validate)
	err := wrapped("")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("wrapped error %q should contain field name", err.Error())
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Error("wrapped error should unwrap to ValidationError")
	}
}
