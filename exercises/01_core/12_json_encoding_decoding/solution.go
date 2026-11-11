package jsonencodingdecoding

import (
	"encoding/json"
	"fmt"
)

// Person represents a person record with name, age, and email fields.
type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// ParsePerson decodes a JSON string into a Person struct.
// Returns an error if the JSON is invalid or cannot be decoded.
func ParsePerson(data string) (Person, error) {
	var p Person
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		return Person{}, fmt.Errorf("parse person: %w", err)
	}
	return p, nil
}

// SerializePerson encodes a Person struct into a JSON string.
// Returns an error if the Person cannot be marshaled.
func SerializePerson(p Person) (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("serialize person: %w", err)
	}
	return string(data), nil
}
