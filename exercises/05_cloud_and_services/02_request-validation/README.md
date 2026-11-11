# Request validation

## Theory

Request validation is a gate between unstructured input and business logic.
Rules are small, composable functions that take a string value and return
either an empty string (valid) or a message (invalid). Collecting all errors
into a single `ValidationError` lets callers fix multiple fields in one pass.

The `errors.Is` interface supports sentinel-pattern matching. `ValidationError`
implements `Is` so that framework-level error handling can recognize
validation failures programmatically.

Field names are resolved by convention so rules stay declarative. This
avoids reflection while keeping the API tight.

## Task

Copy `solution.go` from `TODO.go`, then implement:

1. `FieldError.Error` — return `"Field: Issue"`.
2. `ValidationError.Error` — join all field errors with `"; "`.
3. `ValidationError.Is` — return `true` when `target` is a `*ValidationError`.
4. `Validate` — run every rule against the matching request field, accumulate
   `FieldError`s, and return `nil` only when no errors exist.

## Run

```bash
make 05_cloud_and_services_03_request-validation
```

## Further practice

Add a `Required` helper that returns a Rule enforcing non-empty values. Add a
`MinLen(n int)` helper rule.
