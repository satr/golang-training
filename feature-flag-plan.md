You are working in the `golang-training` repository.

First inspect the existing repository structure, especially several completed exercises. Follow their naming, package, README, test, build-tag, Makefile, and reset-template conventions exactly.

Create a new exercise about **feature flags and deterministic A/B testing**. Select the next available numeric directory under:

`exercises/03_software_engineering/`

Do not overwrite or renumber existing exercises.

Create exactly these five files:

* `doc.go`
* `README.md`
* `solution.go`
* `solution_test.go`
* `TODO.go`

Update the relevant exercise index and Makefile only if existing repository conventions require it.

## Exercise objective

Implement an in-memory, concurrency-safe feature-flag evaluator that:

* Supports a control variant and a treatment variant.
* Supports globally disabled flags.
* Supports rollout percentages from 0 through 100.
* Assigns users deterministically using a stable hash of:

    * flag key;
    * user ID.
* Always assigns the same user to the same variant for the same flag configuration.
* Allows explicit per-user overrides.
* Returns a useful evaluation reason.
* Supports safe runtime configuration replacement.
* Does not perform a remote network request during evaluation.
* Does not expose internal mutable maps to callers.

Use only the Go standard library.

## Suggested public API

You may adapt names slightly if repository conventions require it, but preserve the behavior:

```go
type Variant string

const (
    Control   Variant = "control"
    Treatment Variant = "treatment"
)

type Reason string

const (
    ReasonDisabled Reason = "disabled"
    ReasonOverride Reason = "override"
    ReasonRollout  Reason = "rollout"
)

type Flag struct {
    Enabled    bool
    Percentage int
    Overrides  map[string]Variant
}

type Evaluation struct {
    Variant Reason-compatible variant type
    Reason  Reason
}

type Evaluator struct {
    // concurrency-safe internal state
}

func NewEvaluator(flags map[string]Flag) (*Evaluator, error)

func (e *Evaluator) Evaluate(
    flagKey string,
    userID string,
) (Evaluation, error)

func (e *Evaluator) Replace(flags map[string]Flag) error
```

Correct the `Evaluation` declaration into valid idiomatic Go:

```go
type Evaluation struct {
    Variant Variant
    Reason  Reason
}
```

Define sentinel errors for at least:

* Empty flag key
* Empty user ID
* Unknown flag
* Percentage outside `0..100`
* Invalid override variant

Wrap errors with useful context and make them inspectable through `errors.Is`.

## Evaluation rules

Use this precedence:

1. Validate the flag key and user ID.
2. Return an unknown-flag error if the flag does not exist.
3. If the flag is disabled, return `Control` with `ReasonDisabled`.
4. If the user has an explicit override, return it with `ReasonOverride`.
5. Otherwise, perform stable percentage assignment and return `ReasonRollout`.

Percentage behavior:

* `0` assigns every non-overridden user to `Control`.
* `100` assigns every non-overridden user to `Treatment`.
* Other percentages use a deterministic bucket from `0` through `99`.
* Do not use `math/rand`.
* Use a documented stable standard-library hash, such as FNV-1a.
* Include both flag key and user ID in the hash input.
* Avoid ambiguous concatenation by using an explicit separator or length encoding.

## Configuration ownership

`NewEvaluator` and `Replace` must defensively copy:

* The top-level flags map.
* Every nested overrides map.

Changes made by the caller after construction or replacement must not affect the evaluator.

Validate the complete replacement configuration before publishing it. If validation fails, preserve the previous valid configuration.

Runtime reads and replacements must be safe when called concurrently. Keep evaluation lock scope small and do not return references to internal mutable data.

## Tests

Write deterministic, table-driven tests covering:

* Disabled flag.
* Zero-percent rollout.
* Hundred-percent rollout.
* Stable assignment for repeated evaluations.
* Explicit control override.
* Explicit treatment override.
* Override precedence over rollout percentage.
* Unknown flag.
* Empty flag key.
* Empty user ID.
* Invalid negative percentage.
* Percentage greater than 100.
* Invalid override variant.
* Constructor defensive copying.
* Replacement defensive copying.
* Failed replacement preserving the old configuration.
* Runtime configuration replacement.
* Concurrent evaluation.
* Concurrent evaluation and replacement, suitable for `go test -race`.

Do not use sleeps to create concurrency ordering. Use channels, barriers, and `sync.WaitGroup`. Timeouts may be used only as generous protection against a hanging test.

Do not write a statistically flaky test that expects an exact 10% distribution. Test deterministic properties and boundary behavior instead.

## README requirements

Explain:

* What a feature flag is.
* Release flags, operational flags, kill switches, and experiment flags.
* How A/B testing differs from simply enabling a feature.
* Control versus treatment groups.
* Why user assignment must be stable.
* Why random selection on every request is incorrect.
* Why the flag key must be included in the hash.
* Evaluation precedence.
* Defensive copying and configuration ownership.
* Local evaluation from cached configuration.
* What “calling the flag server” means.
* Why making a network call for every evaluation adds latency and creates a failure dependency.
* Polling or streaming configuration into a local evaluator.
* Safe fallback behavior when configuration is unavailable.
* Why feature flags are not authorization.
* Flag ownership, expiry, cleanup, and technical debt.
* Observability without placing high-cardinality user IDs in metric labels.

Also clearly distinguish these concepts:

* **Feature-flag rollout:** application code chooses behavior.
* **A/B testing:** user groups receive different variants and business outcomes are measured.
* **Canary release:** a small percentage of infrastructure traffic receives a new application version.
* **Blue-green deployment:** traffic switches between two complete environments.

Mention that this exercise implements feature-flag-based A/B assignment, not infrastructure-level blue-green deployment.

## Template requirements

`TODO.go` must:

* Begin with `//go:build todo`.
* Declare the same public API and sentinel errors as `solution.go`.
* Contain clear TODO comments and compiling stubs consistent with repository reset behavior.

`solution.go` must contain the complete reference implementation.

`doc.go` must contain a proper package comment.

## Verification

Run all applicable checks, including:

```bash
gofmt
go vet
go test
go test -race
```

Run the repository’s focused Make target and general validation target if present.

Do not claim checks passed unless they were actually executed. At the end, report:

* Files created or modified.
* Public API implemented.
* Important design decisions.
* Exact verification commands and results.
* Any unresolved limitation.
