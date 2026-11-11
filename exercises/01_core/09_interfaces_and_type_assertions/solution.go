package interfaces_and_type_assertions

import "fmt"

// Stringer is an interface that has a single method String() string.
type Stringer interface {
	String() string
}

// MyString implements the Stringer interface
type MyString struct {
	value string
}

func (m MyString) String() string {
	return m.value
}

// PrintValue prints the value of x depending on its type.
// If x is a string, it prints the string directly.
// If x is an int, it prints the integer with a prefix "int:".
// For any other type, it prints "unknown type".
func PrintValue(x interface{}) {
	switch v := x.(type) {
	case string:
		fmt.Println(v)
	case int:
		fmt.Println("int:", v)
	default:
		fmt.Println("unknown type")
	}
}

// ProcessValue takes an interface{} and returns a string representation.
// It uses type assertions to determine the underlying type of x.
func ProcessValue(x interface{}) string {
	switch v := x.(type) {
	case string:
		return "string: " + v
	case int:
		return fmt.Sprintf("int: %d", v)
	case bool:
		return fmt.Sprintf("bool: %t", v)
	default:
		return "unknown type"
	}
}

// CanPrint checks if value x can be printed using the Stringer interface.
// It returns true if the value implements the Stringer interface, false otherwise.
func CanPrint(x interface{}) bool {
	_, ok := x.(Stringer)
	return ok
}
