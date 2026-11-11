package valuesandtypes

import "fmt"

// ProfileSummary formats a person's typed profile values.
func ProfileSummary(name string, age int, heightMeters float64, active bool) string {
	return fmt.Sprintf("%s is %d years old, %.2f m tall, and %s.", name, age, heightMeters, statusLabel(active))
}

func statusLabel(active bool) string {
	if active {
		return "active"
	}
	return "inactive"
}

// CelsiusToFahrenheit converts a Celsius temperature to Fahrenheit.
func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}
