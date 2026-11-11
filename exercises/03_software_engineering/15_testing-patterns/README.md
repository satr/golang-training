# Testing Patterns

## Theory

Table-driven tests iterate over a slice of inputs and expected outputs. `testing.TB` lets you share logic between tests and benchmarks. Subtests via `t.Run` produce readable output. `TestMain` runs once for the package and can set up or tear down shared resources.

## Task

Implement a calculator with table-driven tests, a shared test helper using `testing.TB`, subtests for edge cases, and optional `TestMain` for setup/teardown logging.

## Run

```bash
make 03_software_engineering_15_testing_patterns
```

## Further practice

Add golden-file testing, mock implementations, and fuzz tests.
