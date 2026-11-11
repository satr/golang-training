# Command-line argument parsing

## Theory

The `flag` package from the standard library provides a simple and idiomatic way to parse command-line flags in Go. Each flag is defined with a name, default value, and usage string. The `flag.Parse()` function processes `os.Args` and populates the registered flags.

Unlike some other languages, Go does not require a full framework for CLI parsing. The `flag` package handles:
- Required and optional flags with defaults
- Type-safe flag definitions (string, int, bool, float, duration)
- Help text via `flag.Usage`
- Errors during parsing

For testability, avoid calling `os.Exit` in library code. Instead, separate flag parsing from business logic.

## Task

Implement a `Runner` type that holds CLI configuration: a `Host` string, a `Port` int, a `Verbose` bool, and a `Timeout` duration. Implement `NewRunner` that parses flag values (simulated — accepts a `[][]string` of flag arguments like `[["-host", "localhost"], ["-port", "8080"]]`) and populates the Runner. Also implement `Run` that returns a formatted summary of the configuration.

The Runner should use default values when flags are not provided: host defaults to "localhost", port defaults to 8080, verbose defaults to false, timeout defaults to 30 seconds.

Implement `ParseArgs` as a standalone function that accepts raw flag arguments as `[]string` and returns a populated Runner and any error.

## Run

```bash
make 01_core_15_command_line_argument_parsing
```

## Further practice

- Add an `--config` flag that reads from a file using `flag.Lookup`
- Implement subcommand parsing by manually scanning `os.Args`
- Explore `flag.Visit` and `flag.VisitAll` for validation