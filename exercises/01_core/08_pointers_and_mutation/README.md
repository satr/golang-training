# Pointers and Mutation

## Theory

In Go, understanding when to use value vs. pointer receivers is crucial for effective mutation and performance.

When you have a method that needs to modify the receiver's state, you must use a pointer receiver. Value receivers create a copy of the struct, so any changes made inside the method won't affect the original struct.

Pointer receivers also provide performance benefits when working with large structs - instead of copying the entire struct, Go passes only the memory address.

## Task

Implement the `Person` struct with methods to modify person's name and age. Create two implementations:
1. One using value receiver (which will not modify the original)
2. One using pointer receiver (which will modify the original)
3. Demonstrate the difference between these approaches in a test

## Run

```bash
go test ./exercises/01_core/08_pointers_and_mutation
```

## Further practice

1. Try to implement a method to increase the person's age by a given amount using a pointer receiver.
2. Create a method that returns a copy of the person with modified fields using value receivers.
3. Explore cases where you might want to use value receivers even for methods that modify state (e.g., when copying is cheaper than referencing).