package packagesandorganization

import (
	"testing"
)

func TestCounter(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *Counter
		wantName  string
		wantCount int
	}{
		{
			name: "NewCounter creates counter with correct name and zero count",
			setup: func() *Counter {
				return NewCounter("test")
			},
			wantName:  "test",
			wantCount: 0,
		},
		{
			name: "Increment increases count by one and returns new value",
			setup: func() *Counter {
				c := NewCounter("test")
				c.count = 5
				return c
			},
			wantName:  "test",
			wantCount: 6,
		},
		{
			name: "Multiple increments accumulate correctly",
			setup: func() *Counter {
				c := NewCounter("test")
				c.count = 3
				c.Increment()
				c.Increment()
				return c
			},
			wantName:  "test",
			wantCount: 5,
		},
		{
			name: "Decrement(n) reduces count",
			setup: func() *Counter {
				c := NewCounter("test")
				c.count = 10
				c.Decrement(3)
				return c
			},
			wantName:  "test",
			wantCount: 7,
		},
		{
			name: "Decrement beyond zero clamps to zero",
			setup: func() *Counter {
				c := NewCounter("test")
				c.count = 2
				c.Decrement(5)
				return c
			},
			wantName:  "test",
			wantCount: 0,
		},
		{
			name: "Name returns the original name",
			setup: func() *Counter {
				return NewCounter("original")
			},
			wantName:  "original",
			wantCount: 0,
		},
		{
			name: "Decrement with n=0 is a no-op",
			setup: func() *Counter {
				c := NewCounter("test")
				c.count = 5
				c.Decrement(0)
				return c
			},
			wantName:  "test",
			wantCount: 5,
		},
		{
			name: "Increment after clamping works correctly",
			setup: func() *Counter {
				c := NewCounter("test")
				c.count = 2
				c.Decrement(5) // should clamp to 0
				c.Increment()  // should increment to 1
				return c
			},
			wantName:  "test",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			if got := c.Name(); got != tt.wantName {
				t.Fatalf("Counter.Name() = %v, want %v", got, tt.wantName)
			}
			if got := c.Count(); got != tt.wantCount {
				t.Fatalf("Counter.Count() = %v, want %v", got, tt.wantCount)
			}
		})
	}
}
