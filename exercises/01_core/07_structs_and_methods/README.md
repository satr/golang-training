# Structs and Methods

## Theory

Go has no classes. Instead, you can define structs that can have methods attached to them.

A struct is a collection of fields.
A method is a function with a special receiver argument.
The receiver appears in its own parameter list, between the func keyword and the method name.

Methods are functions with a receiver. The receiver can be:
- A value (struct, primitive type, etc.)
- A pointer to a value

When you call a method on a value, Go will automatically take a pointer if needed.
This is called "automatic dereferencing" or "automatic addressing".

## Task

Create a simple struct with methods attached to it. You should:
1. Define a struct with several fields
2. Create methods that operate on this struct
3. Understand and use both value and pointer receivers appropriately

## Run

```bash
cd exercises/01_core/07_structs_and_methods
go run solution.go
go test -v
```

## Further practice

1. Add more fields to the struct and corresponding methods
2. Experiment with different receiver types (value vs pointer)
3. Create multiple structs with methods