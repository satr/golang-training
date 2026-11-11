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

// ChangeName returns a new person with the updated name
// Using value receiver (doesn't modify original)
func (p Person) ChangeName(newName string) Person {
	p.name = newName
	return p
}

// ChangeNamePtr changes the person's name in place using pointer receiver
// This modifies the original person
func (p *Person) ChangeNamePtr(newName string) {
	p.name = newName
}

// SetAge returns a new person with the updated age
// Using value receiver (doesn't modify original)
func (p Person) SetAge(age int) Person {
	p.age = age
	return p
}

// SetAgePtr changes the person's age in place using pointer receiver
// This modifies the original person
func (p *Person) SetAgePtr(age int) {
	p.age = age
}
