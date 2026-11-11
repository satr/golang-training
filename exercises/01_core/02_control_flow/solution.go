package controlflow

import "errors"

// ErrInvalidScore reports a score outside the inclusive range 0 through 100.
var ErrInvalidScore = errors.New("score must be between 0 and 100")

// Grade converts a numeric score to a letter grade.
func Grade(score int) (string, error) {
	if score < 0 || score > 100 {
		return "", ErrInvalidScore
	}
	if score >= 90 {
		return "A", nil
	}
	if score >= 80 {
		return "B", nil
	}
	if score >= 70 {
		return "C", nil
	}
	if score >= 60 {
		return "D", nil
	}
	return "F", nil
}

// QualifiesForFreeShipping reports whether an order reaches its customer threshold.
func QualifiesForFreeShipping(orderTotal float64, member bool) bool {
	if orderTotal < 0 {
		return false
	}
	threshold := 50.0
	if member {
		threshold = 25
	}
	return orderTotal >= threshold
}
