# Slices and arrays

## Theory

An array's length is part of its type: `[5]int` is copied when passed by
value. A slice is a small descriptor over an underlying array. Appending to a
slice may reuse that array, so an API that promises not to mutate caller-owned
data should allocate or copy deliberately.

```go
copyOfValues := append([]int(nil), values...)
```

A nil slice is distinct from an allocated empty slice, although both have zero
length. Preserve nil when it conveys "not provided."

## Task

Implement `ReverseFive` to return a reversed array while leaving its input
unchanged. Implement `AppendCopy` so it appends values without changing the
input slice or sharing writable storage with it; preserve a nil result when
both inputs are nil.

## Run

```bash
make 01_core_05_slices_and_arrays
```

## Further practice

Write an in-place slice reversal and document how its ownership contract
differs from `ReverseFive`.

