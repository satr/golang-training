# Exercise index

Each exercise is one focused learning session. Read its theory and tests,
reset `solution.go` from `TODO.go`, implement the behavior, and use the listed
Make target.

## 01: Core Go

| Exercise | Description | Status |
| --- | --- | --- |
| [01 Values and basic types](01_core/01_values_and_types) | Format typed values and perform floating-point arithmetic. | Implemented |
| [02 Control flow](01_core/02_control_flow) | Express validation and business branches with `if` and boolean operators. | Implemented |
| [03 Loops and range](01_core/03_loops_and_range) | Traverse slices and Unicode strings without index mistakes. | Implemented |
| [04 Functions and multiple returns](01_core/04_functions_and_multiple_returns) | Return useful results together with explicit errors. | Implemented |
| [05 Slices and arrays](01_core/05_slices_and_arrays) | Distinguish fixed-size values from slice views and preserve ownership. | Implemented |
| 06 Maps and deterministic ordering | Build map-backed indexes with stable output. | Planned |
| 07 Structs and methods | Model data and attach behavior to useful receiver types. | Planned |
| 08 Pointers and mutation | Choose pointer receivers and make mutation visible. | Planned |
| 09 Interfaces and type assertions | Compose small capabilities and inspect dynamic values safely. | Implemented |
| 10 Errors and wrapped errors | Define, wrap, inspect, and propagate meaningful errors. | Planned |
| 11 Files and `io.Reader` / `io.Writer` | Stream data through standard I/O abstractions. | Planned |
| 12 JSON encoding and decoding | Validate and transform structured JSON. | Planned |
| 13 Packages and standard-library organization | Design cohesive packages with clear exported APIs. | Planned |
| 14 Variables, constants, and zero values | Use declarations and useful zero-value types deliberately. | Planned |
| 15 Command-line argument parsing | Build a testable CLI with `flag` and explicit exit behavior. | Planned |

## 02: Concurrency

1. Goroutines and worker functions
2. Channels and fan-in/fan-out
3. Bounded worker pools
4. Context cancellation
5. Deadlines and timeouts
6. Select statements
7. Backpressure
8. Error propagation from concurrent work
9. Safe shared state and mutexes
10. Race detection and deterministic concurrency tests

## 03: Software engineering

1. Configuration parsing
2. Dependency injection
3. Interfaces and test doubles
4. HTTP handlers
5. HTTP clients with timeouts
6. Middleware
7. Structured logging
8. Metrics boundaries
9. Package architecture
10. Persistence abstractions
11. Transaction-like workflows
12. CLI application design
13. Validation and domain errors
14. Fuzz testing
15. Benchmarking and profiling basics
16. Feature flags and A/B assignment

## 04: Data engineering

1. CSV reading and writing
2. JSON Lines processing
3. Streaming transformations
4. Schema validation
5. Deduplication
6. Aggregation
7. SQLite access with parameterized queries
8. Transactions and idempotent loading
9. Pagination and batching
10. Data quality reports
11. File-based pipelines
12. Incremental processing

## 05: Cloud and services

1. REST API design
2. Graceful HTTP server shutdown
3. Request validation
4. Context propagation
5. Health and readiness endpoints
6. Retries and retry budgets
7. Rate limiting
8. Message-processing abstractions
9. Service configuration
10. Secure secret handling
11. OpenTelemetry boundaries
12. Docker-friendly service layout
13. Kubernetes-oriented health behavior
14. Contract testing
15. Resilient service composition

## Roadmap

Add the remaining Core Go sessions first, then implement each later in
order. Concurrency tests will use channels, barriers, contexts, and explicit
completion signals instead of sleeps. Data exercises will use small local
fixtures. Cloud exercises will keep external systems behind interfaces and
mark optional integration tests with an `integration` build tag.
