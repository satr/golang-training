package valuesandtypes

import "testing"

func TestProfileSummary(t *testing.T) {
	tests := []struct {
		name         string
		person       string
		age          int
		heightMeters float64
		active       bool
		want         string
	}{
		{name: "active profile", person: "Ari", age: 31, heightMeters: 1.725, active: true, want: "Ari is 31 years old, 1.73 m tall, and active."},
		{name: "zero values", person: "", age: 0, heightMeters: 0, active: false, want: " is 0 years old, 0.00 m tall, and inactive."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ProfileSummary(test.person, test.age, test.heightMeters, test.active)
			if got != test.want {
				t.Fatalf("ProfileSummary() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		want    float64
	}{
		{name: "freezing", celsius: 0, want: 32},
		{name: "boiling", celsius: 100, want: 212},
		{name: "same reading", celsius: -40, want: -40},
		{name: "fraction", celsius: 37.5, want: 99.5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CelsiusToFahrenheit(test.celsius); got != test.want {
				t.Fatalf("CelsiusToFahrenheit(%v) = %v, want %v", test.celsius, got, test.want)
			}
		})
	}
}
