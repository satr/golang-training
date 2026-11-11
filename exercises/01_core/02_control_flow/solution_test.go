package controlflow

import (
	"errors"
	"testing"
)

func TestGrade(t *testing.T) {
	tests := []struct {
		score   int
		want    string
		wantErr bool
	}{
		{score: -1, wantErr: true},
		{score: 0, want: "F"},
		{score: 59, want: "F"},
		{score: 60, want: "D"},
		{score: 70, want: "C"},
		{score: 80, want: "B"},
		{score: 90, want: "A"},
		{score: 100, want: "A"},
		{score: 101, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.want+"score", func(t *testing.T) {
			got, err := Grade(test.score)
			if test.wantErr {
				if !errors.Is(err, ErrInvalidScore) {
					t.Fatalf("Grade(%d) error = %v, want ErrInvalidScore", test.score, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Grade(%d) unexpected error: %v", test.score, err)
			}
			if got != test.want {
				t.Fatalf("Grade(%d) = %q, want %q", test.score, got, test.want)
			}
		})
	}
}

func TestQualifiesForFreeShipping(t *testing.T) {
	tests := []struct {
		name   string
		total  float64
		member bool
		want   bool
	}{
		{name: "negative member order", total: -1, member: true, want: false},
		{name: "member below threshold", total: 24.99, member: true, want: false},
		{name: "member at threshold", total: 25, member: true, want: true},
		{name: "guest below threshold", total: 49.99, member: false, want: false},
		{name: "guest at threshold", total: 50, member: false, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := QualifiesForFreeShipping(test.total, test.member); got != test.want {
				t.Fatalf("QualifiesForFreeShipping(%v, %t) = %t, want %t", test.total, test.member, got, test.want)
			}
		})
	}
}
