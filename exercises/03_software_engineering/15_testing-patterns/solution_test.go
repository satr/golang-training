package testingpatterns

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		a, b    float64
		want    float64
		wantErr string
	}{
		{"add", "add", 2, 3, 5, ""},
		{"sub", "sub", 10, 4, 6, ""},
		{"mul", "mul", 3, 7, 21, ""},
		{"div", "div", 10, 2, 5, ""},
		{"div_zero", "div", 1, 0, 0, "divide by zero"},
		{"unknown", "pow", 2, 3, 0, "unknown operation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.op, tt.a, tt.b)
			if tt.wantErr != "" {
				if got.Err == nil {
					t.Errorf("expected error %q", tt.wantErr)
				} else if !strings.Contains(got.Err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want %q", got.Err.Error(), tt.wantErr)
				}
			} else {
				if got.Err != nil {
					t.Errorf("unexpected error: %v", got.Err)
				} else if math.Abs(got.Value-tt.want) > 1e-9 {
					t.Errorf("got %v, want %v", got.Value, tt.want)
				}
			}
		})
	}
}

func TestFormatResult(t *testing.T) {
	tests := []struct {
		input CalcResult
		want  string
	}{
		{CalcResult{Value: 5}, "result: 5"},
		{CalcResult{Err: errors.New("something went wrong")}, "error: something went wrong"},
	}
	for _, tt := range tests {
		got := FormatResult("add", 2, 3, tt.input)
		if got != tt.want {
			t.Errorf("got %q, want %q", got, tt.want)
		}
	}
}
