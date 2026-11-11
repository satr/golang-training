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
	if divisor == 0 {
		return 0, ErrDivisionByZero
	}
	return dividend / divisor, nil
}

// MinMax returns the smallest and largest input values.
func MinMax(values []int) (int, int, error) {
	if len(values) == 0 {
		return 0, 0, ErrEmptyInput
	}
	minimum, maximum := values[0], values[0]
	for _, value := range values[1:] {
		if value < minimum {
			minimum = value
		}
		if value > maximum {
			maximum = value
		}
	}
	return minimum, maximum, nil
}
