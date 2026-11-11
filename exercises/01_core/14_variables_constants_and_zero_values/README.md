# Variables, constants, and zero values

## Theory

Go gives every variable a zero value on declaration:
- Numeric types default to `0`
- Booleans default to `false`
- Strings default to `""`
- Pointers, slices, maps, and channels default to `nil`
- Structs default to the zero value of each field

Constants are evaluated at compile time and cannot be reassigned. The `iota` keyword auto-increments through a `const` block, making it ideal for enumerations.

Using meaningful zero values reduces initialization noise. For example, a `sync.Mutex` is safe to use at its zero value, and an empty slice behaves the same as `nil` for most operations.

## Task

Implement the `LogLevel` type using `iota` for an enumeration of log severity levels: `LevelDebug`, `LevelInfo`, `LevelWarn`, `LevelError`. Implement the `Label()` method on `LogLevel` that returns the human-readable label for each level (e.g., "DEBUG", "INFO", "WARN", "ERROR").

Implement `GetDefaultConfig` that returns a `Config` struct at its zero value. The Config should have these fields: `Timeout int`, `RetryLimit int`, `Verbose bool`, `Labels []string`, `LogLevel LogLevel`. Each field must be at its zero value (0, 0, false, nil, LevelDebug).

Implement `ValidateConfig` that checks if the config's Timeout is greater than zero and RetryLimit is between 0 and 100 inclusive. Returns nil if valid, or an error describing what is wrong.

## Run

```bash
make 01_core_14_variables_constants_and_zero_values
```

## Further practice

- Add a `String()` method to `LogLevel` for pretty printing
- Implement `ParseLevel(string) (LogLevel, error)` to parse level labels from text
- Explore typed constants vs untyped constants