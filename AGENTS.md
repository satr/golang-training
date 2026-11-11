# Repository instructions

## Purpose

This repository teaches idiomatic Go through progressive, test-driven
exercises. Prefer explicit errors, useful interfaces, ownership clarity,
composition, context cancellation, deterministic concurrency, and the standard
library over translations of patterns from other languages.

## Agent workflow

Select the `orchestrator` profile for repository changes. It coordinates the
mandatory `planner` -> `implementer` -> `verifier` -> `validator` ->
`security-auditor` sequence. Exercise changes also receive a read-only
`curriculum` review after planning. Subagents are deliberately unavailable for
automatic model selection and must be delegated by the orchestrator.

Every handoff must stay scoped to the user goal and include constraints,
affected files, and existing evidence. Failed gates return actionable findings
to the implementer, after which verification, validation, and security repeat.
Documentation-only requests may be handled directly when they do not change
code, configuration, dependencies, workflows, or the exercise learning
contract. Changes to exercise documentation still require the orchestrated
workflow and curriculum review.

## Exercise contract

Exercises live at `exercises/NN_name/NN_name/`, start at `01`, and
remain contiguous. Every exercise contains:

- `README.md` with `Theory`, `Task`, `Run`, and `Further practice` headings in
  that order. Theory must appear before the task and must not reveal the final
  implementation.
- `solution.go`, the active implementation.
- `TODO.go`, a `//go:build todo` reset template with matching declarations,
  clear TODOs, compilable stubs, and no hidden finished implementation.
- `solution_test.go`, containing deterministic behavior-focused tests.
- `doc.go`, documenting the package purpose and learning objective.

Use table-driven tests where appropriate. Cover normal behavior, boundaries,
and errors. Never depend on sleeps, timing luck, public networks, credentials,
paid APIs, or large downloads. Mark external integration tests with the
`integration` build tag.

## Authoring rules

1. Explain the idea before asking the learner to implement it.
2. Use standard-library-first solutions. Add a dependency only when it teaches
   an important concept, and document why.
3. Keep packages focused and errors explicit. Wrap errors when adding context.
4. Keep `TODO.go` resettable and free of final logic.
5. Add each exercise to `exercises/README.md` and the Makefile help and focused
   targets.
6. After changing an exercise, run its focused target, formatting, layout
   verification, and relevant static checks.

Repository changes must pass `go test ./...`, `go test -race ./...`,
`go vet ./...`, formatting, `golangci-lint`, and the layout verifier.
Integration tests must be explicitly marked and run separately.

Use the focused Make target plus `make layout` for an exercise change. Use
`make check` for cross-cutting changes, `make race` for concurrency or
synchronization changes, and `make vuln` for dependency or security-sensitive
changes. Verifiers prove requested behavior; validators own repository quality
and regression checks and must not substitute one responsibility for the other.

Never put secrets, credentials, tokens, private data, private endpoints, or
sensitive values in source, tests, fixtures, logs, prompts, or documentation.
Do not read unrelated credential stores or sensitive files. In agent handoffs
and output, identify only the affected location and category and replace values
with `[REDACTED]`.
