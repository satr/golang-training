package structs_and_methods

// Person represents a person with basic information
type Person struct {
	Name  string
	Age   int
	Email string
}

// Greet returns a greeting message for the person
func (p Person) Greet() string {
	return "Hello, my name is " + p.Name
}

// IsAdult returns true if the person is 18 years or older
func (p Person) IsAdult() bool {
	return p.Age >= 18
}

// UpdateEmail updates the email field of the person
func (p *Person) UpdateEmail(email string) {
	p.Email = email
}
