# Control flow

## Theory

Go uses `if`, `else if`, and `else` without parentheses around the condition.
Boolean expressions can combine comparisons with `&&`, `||`, and `!`.
Validate invalid input before classifying valid values so each branch has one
clear responsibility.

```go
if value < minimum || value > maximum {
    return ErrOutOfRange
}
```

Go encourages explicit error returns rather than exceptional control flow.
Callers can use `errors.Is` to recognize a documented sentinel error.

## Task

Implement `Grade` for scores from 0 through 100 using the boundaries in the
tests, returning `ErrInvalidScore` outside that range. Implement
`QualifiesForFreeShipping`: members qualify at 25 units and non-members at 50.
Negative order totals never qualify.

## Run

```bash
make 01_core_02_control_flow
```

## Further practice

Represent grade boundaries as data and compare the trade-off with direct
branches.

