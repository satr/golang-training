//go:build todo

package structs_and_methods

// TODO: Define a struct called Person with fields Name (string), Age (int), and Email (string)
// TODO: Add a method called Greet that takes no parameters and returns a string
// TODO: Add a method called IsAdult that returns a boolean indicating if the person is 18 or older
// TODO: Add a method called UpdateEmail that takes a string parameter and updates the email field

// Person represents a person with basic information
type Person struct {
	// TODO: Add fields here
}

// Greet returns a greeting message for the person
func (p Person) Greet() string {
	// TODO: Implement this method
	return ""
}

// IsAdult returns true if the person is 18 years or older
func (p Person) IsAdult() bool {
	// TODO: Implement this method
	return false
}

// UpdateEmail updates the email field of the person
func (p *Person) UpdateEmail(email string) {
	// TODO: Implement this method
}
