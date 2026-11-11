package errors_and_wrapped_errors

import (
	"errors"
	"fmt"
)

// ProcessData takes a string and returns an error with context
// if the data is invalid.
func ProcessData(data string) error {
	if data == "" {
		return fmt.Errorf("invalid data: %w", errors.New("data cannot be empty"))
	}
	return nil
}

// ReadFile attempts to read a file and returns an error with context
// if it fails.
func ReadFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("failed to read file: %w", errors.New("filename cannot be empty"))
	}

	// Simulate a file reading error
	return fmt.Errorf("failed to open file %q: %w", filename, errors.New("file not found"))
}

// SimulateErrorSimulates an error that should be wrapped
func SimulateError() error {
	return errors.New("simulated error")
}

// CreateWrappedError creates a wrapped error with additional context
func CreateWrappedError(err error) error {
	return fmt.Errorf("operation failed: %w", err)
}
