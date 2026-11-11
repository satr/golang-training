# Errors and Wrapped Errors

## Theory

In Go, errors are values that implement the `error` interface. An error is typically created using `errors.New()` or `fmt.Errorf()`. When you need to provide more context about an error, you should wrap it to preserve the original error information.

The `fmt.Errorf()` function provides a way to create formatted errors with additional context by using the `%w` verb to wrap the underlying error. This allows for error chaining where you can understand what happened at each level of the call stack.

Example:
```go
err := fmt.Errorf("failed to read file: %w", os.ErrNotExist)
```

## Task

Implement a function that returns an error with context information about what happened in case of failure. Use the `fmt.Errorf` and `%w` verb to wrap errors properly.

## Run

```bash
cd exercises/01_core/10_errors_and_wrapped_errors
go test -v
```

## Further practice

- Try creating a scenario where you have nested function calls that each wrap errors
- Experiment with `errors.Is()` and `errors.As()` for error type checking and unwrapping