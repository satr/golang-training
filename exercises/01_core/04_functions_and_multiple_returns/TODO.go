//go:build todo

package multiplereturns

import "errors"

var (
	// ErrDivisionByZero reports an invalid zero divisor.
	ErrDivisionByZero = errors.New("division by zero")
	// ErrEmptyInput reports that a calculation requires at least one value.
	ErrEmptyInput = errors.New("input is empty")
)

// Divide divides dividend by divisor.
func Divide(dividend, divisor float64) (float64, error) {
	// TODO: Reject a zero divisor and otherwise return the quotient.
	return 0, nil
}

// MinMax returns the smallest and largest input values.
func MinMax(values []int) (int, int, error) {
	// TODO: Reject empty input and find both results in one traversal.
	return 0, 0, nil
}
