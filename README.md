# Go Training

An exercise-driven path from Go fundamentals to concurrent programs,
production software engineering, data pipelines, and cloud-native services.
Each exercise is a small package with concise theory, a resettable template,
deterministic tests, and a focused command.

## Requirements

[Go 1.27.1](https://go.dev/dl/) or newer and `make` are required. The module
path is `github.com/satr/golang-training`.

## Installation

```bash
git clone https://github.com/satr/golang-training.git
cd golang-training
make install
```

`make install` places pinned development tools in `.bin/`; it does not modify
your system-wide Go installation.

## Learning workflow

1. Choose an exercise from the [exercise index](exercises/README.md).
2. Read its `Theory` and tests before changing code. Tests define observable
   behavior, including boundaries and errors.
3. Reset the active implementation when you want a clean attempt:

   ```bash
   make reset EXERCISE=exercises/01_core/01_values_and_types
   ```

4. Implement the TODOs in `solution.go`, then run only that exercise:

   ```bash
   make 01_core_01_values_and_types
   ```

5. Run the repository checks before moving on:

   ```bash
   make check
   ```

Attempt `solution.go` yourself before consulting hints or comparing against a
working revision. Prefer small changes and use test failures as feedback.

To deliberately reset every exercise, run `make reset-all` and type the
requested confirmation. This operation overwrites all active solutions.

## AI agent workflow

For repository changes, select the `orchestrator` profile in `.github/agents/`.
It delegates the required planning, implementation, behavior verification,
repository validation, and security gates. Exercise work additionally uses the
`curriculum` specialist to preserve the progressive Go teaching contract.

Subagent profiles deliberately omit a fixed model and disable automatic model
invocation so the orchestrator owns every handoff. A failed gate returns scoped
findings to the implementer and all three independent gates run again. Agent
prompts, logs, and handoffs must never contain secrets or private data; report
only the location and category with values replaced by `[REDACTED]`.

## Commands

```text
make install
make test
make lint
make fmt
make vet
make race
make check
make reset EXERCISE=exercises/01_core/01_values_and_types
make reset-all
```

Additional targets include `make coverage`, `make vuln`, `make layout`, and
`make help`. Focused targets run exactly one exercise package.

## Track overview

| Track | Curriculum |
| --- | --- |
| 01 Core Go | Values and basic types; variables, constants, and zero values; control flow and boolean expressions; loops and range; functions and multiple returns; slices and arrays; maps and deterministic ordering; structs and methods; pointers and mutation; interfaces and type assertions; errors and wrapped errors; files and I/O interfaces; JSON; packages; CLI parsing |
| 02 Concurrency | Goroutines; channels and fan-in/fan-out; bounded pools; cancellation; deadlines; select; backpressure; concurrent errors; mutexes; race detection |
| 03 Software engineering | Configuration; dependency injection; test doubles; HTTP handlers and clients; middleware; logging; metrics; architecture; persistence; workflows; CLI design; validation; fuzzing; benchmarking |
| 04 Data engineering | CSV; JSON Lines; streams; schemas; deduplication; aggregation; SQLite; transactions and idempotency; pagination; quality reports; file pipelines; incremental processing |
| 05 Cloud and services | REST; graceful shutdown; validation; context propagation; health; retries; rate limits; messages; configuration; secrets; telemetry; containers; Kubernetes health; contracts; resilient composition |

The initial release implements the first five Core Go sessions. The complete
planned sequence and future roadmap are in
[`exercises/README.md`](exercises/README.md).

## Quality checks

`make check` verifies layout, formatting, static analysis, lint rules, and unit
tests. Run `make race` for race-enabled tests and `make vuln` to scan code and
its actual dependency graph with `govulncheck`. Standard exercises require no
network, credentials, external service, or large download.
