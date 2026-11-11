package jsonencodingdecoding

import (
	"reflect"
	"testing"
)

func TestParsePerson(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Person
		wantErr bool
	}{
		{
			name: "valid JSON with all fields",
			data: `{"name":"Alice","age":30,"email":"alice@example.com"}`,
			want: Person{
				Name:  "Alice",
				Age:   30,
				Email: "alice@example.com",
			},
			wantErr: false,
		},
		{
			name: "valid JSON missing email field",
			data: `{"name":"Bob","age":25}`,
			want: Person{
				Name:  "Bob",
				Age:   25,
				Email: "",
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON string",
			data:    `{"name":"Charlie","age":35`,
			want:    Person{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePerson(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePerson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializePerson(t *testing.T) {
	tests := []struct {
		name    string
		person  Person
		want    string
		wantErr bool
	}{
		{
			name: "valid Person structure",
			person: Person{
				Name:  "Alice",
				Age:   30,
				Email: "alice@example.com",
			},
			want:    `{"name":"Alice","age":30,"email":"alice@example.com"}`,
			wantErr: false,
		},
		{
			name:    "empty Person",
			person:  Person{},
			want:    `{"name":"","age":0,"email":""}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializePerson(tt.person)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializePerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SerializePerson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	original := Person{
		Name:  "David",
		Age:   40,
		Email: "david@example.com",
	}

	jsonStr, err := SerializePerson(original)
	if err != nil {
		t.Fatalf("SerializePerson failed: %v", err)
	}

	parsed, err := ParsePerson(jsonStr)
	if err != nil {
		t.Fatalf("ParsePerson failed: %v", err)
	}

	if !reflect.DeepEqual(original, parsed) {
		t.Errorf("Round trip failed. Original: %v, Parsed: %v", original, parsed)
	}
}

func TestSerializePersonFieldNames(t *testing.T) {
	person := Person{
		Name:  "Eve",
		Age:   35,
		Email: "eve@example.com",
	}

	jsonStr, err := SerializePerson(person)
	if err != nil {
		t.Fatalf("SerializePerson failed: %v", err)
	}

	// Check that JSON field names match struct tags
	expected := `{"name":"Eve","age":35,"email":"eve@example.com"}`
	if jsonStr != expected {
		t.Errorf("SerializePerson field names incorrect. Got: %s, want: %s", jsonStr, expected)
	}
}
