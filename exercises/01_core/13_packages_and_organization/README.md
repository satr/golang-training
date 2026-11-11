# Packages and standard-library organization

## Theory

Go packages are the primary unit of code organization. A well-designed package groups related types and functions behind a cohesive exported API while keeping implementation details hidden using unexported (lowercase) names.

Guidelines for package design:
- Keep the package focused on a single responsibility
- Export only what callers need; keep internals unexported
- When a zero value is not useful, provide an exported constructor function (e.g., `NewX`)
- Place methods on the type that owns the data
- Document exported identifiers; documentation is what shows up in `go doc`

The standard library demonstrates these patterns consistently: `net/http`, `bufio`, `fmt`, and `os` all export a minimal set of types and functions while keeping internals private.

## Task

Implement a `Counter` type that encapsulates counting behavior. The Counter should:
- Be created only through `NewCounter(string)` — the zero value is invalid
- Track a name and a mutable count
- Expose `Increment()`, `Decrement(int)`, and `Count() int` methods
- Expose `Name() string` method
- Start at a count of zero

The internal state (count) must be unexported and only accessible through methods.

## Run

```bash
make 01_core_13_packages_and_organization
```

## Further practice

- Add a `Reset()` method that sets the count back to zero
- Experiment with `NewCounterWithInitial(string, int)` that validates the initial value is non-negative
- Split the Counter into its own package to see how imports change