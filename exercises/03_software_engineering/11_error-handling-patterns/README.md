# Error Handling Patterns

## Theory

Go errors are values. Use sentinel errors (`var ErrNotFound = ...`) for known failure modes, `fmt.Errorf("context: %w", err)` to wrap, and `errors.Is`/`errors.As` to check. Custom error types carry structured context for callers that need it.

## Task

Implement a validator with sentinel errors, a wrapper function that adds context, and a custom error type for validation failures.

## Run

```bash
make 03_software_engineering_11_error_handling_patterns
```

## Further practice

Add a multi-error aggregator that combines several errors.
