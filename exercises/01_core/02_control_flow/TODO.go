//go:build todo

package controlflow

import "errors"

// ErrInvalidScore reports a score outside the inclusive range 0 through 100.
var ErrInvalidScore = errors.New("score must be between 0 and 100")

// Grade converts a numeric score to a letter grade.
func Grade(score int) (string, error) {
	// TODO: Validate score, then return the letter for its inclusive boundary.
	return "", nil
}

// QualifiesForFreeShipping reports whether an order reaches its customer threshold.
func QualifiesForFreeShipping(orderTotal float64, member bool) bool {
	// TODO: Reject negative totals and apply the appropriate threshold.
	return false
}
