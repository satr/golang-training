package multiplereturns

import (
	"errors"
	"testing"
)

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		dividend float64
		divisor  float64
		want     float64
		wantErr  error
	}{
		{name: "whole quotient", dividend: 12, divisor: 3, want: 4},
		{name: "fractional quotient", dividend: 1, divisor: 4, want: 0.25},
		{name: "negative divisor", dividend: 9, divisor: -3, want: -3},
		{name: "zero divisor", dividend: 9, divisor: 0, wantErr: ErrDivisionByZero},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Divide(test.dividend, test.divisor)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Divide() error = %v, want %v", err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("Divide() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		name    string
		values  []int
		wantMin int
		wantMax int
		wantErr error
	}{
		{name: "nil", values: nil, wantErr: ErrEmptyInput},
		{name: "empty", values: []int{}, wantErr: ErrEmptyInput},
		{name: "one", values: []int{7}, wantMin: 7, wantMax: 7},
		{name: "mixed", values: []int{8, -3, 12, 4, -1}, wantMin: -3, wantMax: 12},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotMin, gotMax, err := MinMax(test.values)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("MinMax() error = %v, want %v", err, test.wantErr)
			}
			if gotMin != test.wantMin || gotMax != test.wantMax {
				t.Fatalf("MinMax() = (%d, %d), want (%d, %d)", gotMin, gotMax, test.wantMin, test.wantMax)
			}
		})
	}
}
