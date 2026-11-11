# Maps and Deterministic Ordering

## Theory

In Go, maps are unordered collections of key-value pairs. The language specification does not guarantee any particular order when iterating over map keys. This means that if you iterate over a map multiple times, you might get the keys in different orders.

This behavior exists for performance reasons, as Go's map implementation doesn't maintain insertion order or sort keys. However, this also means that code relying on map iteration order is not deterministic and can lead to subtle bugs.

The Go runtime does randomize the starting point of map iteration to prevent algorithms from accidentally depending on the iteration order. This makes maps "unpredictable" but more secure against certain types of attacks.

When building indexes or other data structures where order matters, you need to be aware of this behavior and take steps to maintain determinism through explicit ordering.

## Task

Implement a function that builds an index using a map, but ensures deterministic output by explicitly sorting keys. The function should:
1. Take a slice of strings as input
2. Build a map where keys are the unique strings and values are their counts
3. Return a sorted slice of key-value pairs (as tuples) to ensure deterministic output

## Run

```
go test ./exercises/01_core/06_maps_and_deterministic_ordering
```

## Further practice

- Modify your implementation to handle more complex data types as keys
- Implement different sorting strategies for the index
- Explore how Go's map implementation works internally and why it doesn't preserve order