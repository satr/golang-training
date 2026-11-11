package pointers_and_mutation

import (
	"testing"
)

func TestPersonValueReceiver(t *testing.T) {
	person := NewPerson("Alice", 30)

	// Test that value receiver methods don't modify the original
	newPerson := person.ChangeName("Bob")
	if newPerson.Name() != "Bob" {
		t.Errorf("Expected name 'Bob', got %s", newPerson.Name())
	}

	if person.Name() != "Alice" {
		t.Errorf("Expected original person's name to remain 'Alice', got %s", person.Name())
	}

	// Test age modification with value receiver
	newPerson2 := person.SetAge(35)
	if newPerson2.Age() != 35 {
		t.Errorf("Expected age 35, got %d", newPerson2.Age())
	}

	if person.Age() != 30 {
		t.Errorf("Expected original person's age to remain 30, got %d", person.Age())
	}
}

func TestPersonPointerReceiver(t *testing.T) {
	person := NewPerson("Alice", 30)

	// Test that pointer receiver methods modify the original
	person.ChangeNamePtr("Bob")
	if person.Name() != "Bob" {
		t.Errorf("Expected name 'Bob', got %s", person.Name())
	}

	// Test age modification with pointer receiver
	person.SetAgePtr(35)
	if person.Age() != 35 {
		t.Errorf("Expected age 35, got %d", person.Age())
	}
}

func TestPersonCombined(t *testing.T) {
	person := NewPerson("Charlie", 25)

	// Mix value and pointer receiver methods
	newPerson := person.ChangeName("David") // Value receiver - doesn't change original
	if person.Name() != "Charlie" {
		t.Errorf("Expected original person's name to remain 'Charlie', got %s", person.Name())
	}
	if newPerson.Name() != "David" {
		t.Errorf("Expected new person's name to be 'David', got %s", newPerson.Name())
	}

	person.ChangeNamePtr("Eve") // Pointer receiver - does change original
	if person.Name() != "Eve" {
		t.Errorf("Expected name 'Eve', got %s", person.Name())
	}

	person.SetAgePtr(40) // Pointer receiver - does change original
	if person.Age() != 40 {
		t.Errorf("Expected age 40, got %d", person.Age())
	}
}
