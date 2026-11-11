# Functions and multiple return values

## Theory

Go functions can return several values. A common API returns a useful result
and an `error`; callers must check the error before using the result. Named
sentinel errors let callers use `errors.Is` without comparing error text.

```go
value, err := operation()
if err != nil {
    return err
}
```

Multiple ordinary results are also useful when they belong together, such as
the minimum and maximum of one input.

## Task

Implement `Divide`, returning `ErrDivisionByZero` when the divisor is zero.
Implement `MinMax` in one traversal, returning `ErrEmptyInput` for nil or empty
slices.

## Run

```bash
make 01_core_04_functions_and_multiple_returns
```

## Further practice

Add a function returning minimum, maximum, and average, then consider whether
a small result struct would make the API clearer.

