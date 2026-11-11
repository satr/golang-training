package structs_and_methods

import (
	"testing"
)

func TestPersonGreet(t *testing.T) {
	person := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	greeting := person.Greet()
	expected := "Hello, my name is Alice"

	if greeting != expected {
		t.Errorf("Expected greeting %q, got %q", expected, greeting)
	}
}

func TestPersonIsAdult(t *testing.T) {
	adult := Person{
		Name:  "Bob",
		Age:   25,
		Email: "bob@example.com",
	}

	if !adult.IsAdult() {
		t.Errorf("Expected adult to be adult, but got false")
	}

	minor := Person{
		Name:  "Charlie",
		Age:   17,
		Email: "charlie@example.com",
	}

	if minor.IsAdult() {
		t.Errorf("Expected minor to not be adult, but got true")
	}
}

func TestPersonUpdateEmail(t *testing.T) {
	person := Person{
		Name:  "David",
		Age:   35,
		Email: "david@example.com",
	}

	person.UpdateEmail("newdavid@example.com")

	if person.Email != "newdavid@example.com" {
		t.Errorf("Expected email to be %q, got %q", "newdavid@example.com", person.Email)
	}
}
