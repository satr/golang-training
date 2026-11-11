//go:build todo

package jsonencodingdecoding

// ParsePerson decodes a JSON string into a Person struct.
// Returns an error if the JSON is invalid or cannot be decoded.
func ParsePerson(data string) (Person, error) {
	// TODO: Use json.Unmarshal to decode the JSON string into a Person.
	// TODO: Return the Person and any error.
	return Person{}, nil
}

// SerializePerson encodes a Person struct into a JSON string.
// Returns an error if the Person cannot be marshaled.
func SerializePerson(p Person) (string, error) {
	// TODO: Use json.Marshal to encode the Person.
	// TODO: Return the JSON string and any error.
	return "", nil
}
