//go:build todo

package pointers_and_mutation

// Person represents a person with a name and age
type Person struct {
	name string
	age  int
}

// NewPerson creates a new person with the given name and age
func NewPerson(name string, age int) Person {
	return Person{name: name, age: age}
}

// Name returns the person's name
func (p Person) Name() string {
	return p.name
}

// Age returns the person's age
func (p Person) Age() int {
	return p.age
}

// TODO: Implement a method that modifies the person's name using a value receiver.
// This method should not modify the original person when called.
func (p Person) ChangeName(newName string) Person {
	// TODO: Implement this method to return a new person with updated name
	return p
}

// TODO: Implement a method that modifies the person's name using a pointer receiver.
// This method should modify the original person when called.
func (p *Person) ChangeNamePtr(newName string) {
	// TODO: Implement this method to change the person's name in place
}

// TODO: Implement a method that modifies the person's age using a value receiver.
// This method should not modify the original person when called.
func (p Person) SetAge(age int) Person {
	// TODO: Implement this method to return a new person with updated age
	return p
}

// TODO: Implement a method that modifies the person's age using a pointer receiver.
// This method should modify the original person when called.
func (p *Person) SetAgePtr(age int) {
	// TODO: Implement this method to change the person's age in place
}
