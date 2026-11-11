package testingpatterns

import (
	"errors"
	"strconv"
)

// CalcResult holds the result of a calculation.
type CalcResult struct {
	Value float64
	Err   error
}

// Calculate performs basic arithmetic operations.
func Calculate(op string, a, b float64) CalcResult {
	switch op {
	case "add":
		return CalcResult{Value: a + b}
	case "sub":
		return CalcResult{Value: a - b}
	case "mul":
		return CalcResult{Value: a * b}
	case "div":
		if b == 0 {
			return CalcResult{Err: errors.New("divide by zero")}
		}
		return CalcResult{Value: a / b}
	default:
		return CalcResult{Err: errors.New("unknown operation")}
	}
}

// FormatResult formats a CalcResult for display.
func FormatResult(op string, a, b float64, r CalcResult) string {
	if r.Err != nil {
		return "error: " + r.Err.Error()
	}
	if r.Value == float64(int(r.Value)) {
		return "result: " + strconv.FormatInt(int64(r.Value), 10)
	}
	return "result: " + strconv.FormatFloat(r.Value, 'f', -1, 64)
}
