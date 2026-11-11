package errors_and_wrapped_errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestProcessData(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    "hello",
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessData(tt.data)
			if tt.wantErr && err == nil {
				t.Errorf("ProcessData() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ProcessData() expected no error, got %v", err)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid filename",
			filename: "test.txt",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			filename: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadFile(tt.filename)
			if tt.wantErr && err == nil {
				t.Errorf("ReadFile() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ReadFile() expected no error, got %v", err)
			}
		})
	}
}

func TestCreateWrappedError(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := CreateWrappedError(baseErr)

	if wrappedErr == nil {
		t.Fatalf("CreateWrappedError() returned nil error")
	}

	expected := "operation failed: base error"
	actual := wrappedErr.Error()
	if actual != expected {
		t.Errorf("CreateWrappedError() = %q, want %q", actual, expected)
	}

	// Test that we can unwrap the error
	var unwrappedErr error
	if !errors.As(wrappedErr, &unwrappedErr) {
		t.Errorf("CreateWrappedError() should be able to unwrap")
	}
	if unwrappedErr.Error() != "base error" {
		t.Errorf("Unwrapped error = %q, want %q", unwrappedErr.Error(), "base error")
	}
}

func TestSimulateError(t *testing.T) {
	err := SimulateError()
	if err == nil {
		t.Fatalf("SimulateError() returned nil error")
	}

	expected := "simulated error"
	actual := err.Error()
	if actual != expected {
		t.Errorf("SimulateError() = %q, want %q", actual, expected)
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test that wrapping preserves error information
	baseErr := SimulateError()
	wrappedErr := CreateWrappedError(baseErr)

	// The wrapped error should contain the original error message
	if !errors.Is(wrappedErr, baseErr) {
		t.Errorf("CreateWrappedError() should be able to check for original error")
	}

	// Test error formatting
	expectedFormat := "operation failed: simulated error"
	actualFormat := fmt.Sprintf("%v", wrappedErr)
	if actualFormat != expectedFormat {
		t.Errorf("Error formatted as %q, want %q", actualFormat, expectedFormat)
	}
}
