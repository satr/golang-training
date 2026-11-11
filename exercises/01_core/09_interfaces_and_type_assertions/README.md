# Interfaces and Type Assertions

## Theory

Go's interfaces are a powerful feature that allow you to define behavior contracts. An interface in Go is satisfied implicitly - if a type implements all the methods of an interface, it automatically implements that interface without needing to explicitly declare it.

Type assertions let you extract the concrete value from an interface value. A type assertion takes an interface value and extracts from it the concrete value of a specified type.

## Task

In this exercise, you'll implement functions that demonstrate interfaces and type assertions:
1. Create an interface for types that can be "printed"
2. Implement a function to print values using type assertions
3. Write a function that uses interface{} and type assertions to determine and handle different value types

## Run

```bash
cd exercises/01_core/09_interfaces_and_type_assertions/09_interfaces_and_type_assertions
go test -v
```

## Further practice

- Try creating additional interfaces with different method sets
- Explore how interfaces can be composed
- Implement a function that returns an interface value and use type assertions to handle it