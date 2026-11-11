# Feature Flags and A/B Assignment

## Theory

Feature flags let you toggle behavior without deploying new code. A flag key
maps to configuration — enabled status, rollout percentage, and user overrides
— that an *evaluator* consults at runtime.

### Kinds of flags

- **Release flags** control visibility of unfinished features. They are removed
  once the feature ships.
- **Operational flags** guard risky runtime paths (e.g. enabling a new
  downstream) and can stay in the codebase long-term.
- **Kill switches** disable specific functionality entirely, usually for safety.
- **Experiment flags** partition users into groups so you can measure the
  impact of a change.

### A/B testing vs. feature toggle

A feature toggle decides *whether* code runs. A/B testing takes that further:
users are assigned to a **control** group (existing behavior) or a **treatment**
group (new behavior) so business outcomes can be compared statistically.

### Control vs. treatment groups

The control group receives the known-baseline behavior. The treatment group
receives the change under investigation. Assignments must be deterministic:
the same user must always land in the same group for a given flag.

### Stable user assignment

A per-request random choice is wrong: the same user would flip between control
and treatment over successive requests, contaminating every metric. Instead,
hash a stable identifier (user ID) together with the flag key to produce a
deterministic bucket. The flag key must be part of the hash so that assignments
for one flag are independent of another — user "alice" can be treatment for
flag-1 and control for flag-2 simultaneously.

### Evaluation precedence

1. **Validate** inputs — reject empty flag keys and user IDs.
2. **Unknown flag** — an unrecognized key returns an error (safe default).
3. **Disabled** — an explicitly disabled flag always returns control.
4. **Override** — a per-user override forces a specific variant.
5. **Rollout** — the deterministic bucket is compared against the target
   percentage.

### Defensive copying and config ownership

A flag server hands the application a new configuration block. The evaluator
takes ownership of a *copy*; the caller may freely reuse or discard the original
slice. Defensive copying prevents callers from mutating shared state that
corrupts concurrent readers.

### Local evaluation from cached config

A well-designed evaluator reads only local memory after a config refresh. It
does *not* call back to the flag server during evaluation, so a single
downstream failure cannot fan out into every request.

### What "calling the flag server" means

Fetching flags means sending a network request to a configuration service. That
adds latency to the critical path and creates a failure mode: if the service is
unavailable, every request might fail.

### Polling vs. streaming refresh

Applications poll a known endpoint on a timer, or open a long-lived stream
(SSE, gRPC) and update when new config arrives. Polling is simpler; streaming
is more responsive.

### Safe fallback

When configuration refresh fails, the evaluator keeps the previous valid
configuration. Evaluation continues against stale data rather than failing
open.

### Feature flags are not authorization

Flags steer behavior; they do not enforce access control. Never use feature
flags as a security boundary.

### Governance: ownership, expiry, cleanup

Every flag must have an owner and an expiry date. Expired flags become
technical debt: they accumulate, slow evaluation, obscure real behavior, and
make audits harder. Treat flag cleanup as a recurring chore, not an afterthought.

### Observability

Log evaluations and track rollout coverage, but avoid high-cardinality user IDs
in metric labels. Aggregate by flag key and variant; store per-user results
where full fidelity is available.

---

Related concepts that are *not* implemented here:

- **Feature-flag rollout** — application code chooses behavior based on flag
  evaluation. This is what this exercise implements.
- **A/B testing** — user groups receive different variants so business metrics
  can be compared. Feature-flag-based assignment powers A/B testing.
- **Canary release** — a small percentage of infrastructure traffic is routed
  to a new version. This is network-level, not application-level.
- **Blue-green deployment** — traffic switches between two complete environments.
  Fully infrastructure-driven.

This exercise implements feature-flag-based A/B assignment, not
infrastructure-level blue-green deployment.

## Task

Implement an in-memory feature-flag evaluator.

Provide the following types:

- `Variant` — a `string`-based type with `Control` and `Treatment` constants
  and an `IsValid() bool` method.
- `Reason` — a `string`-based type with `ReasonDisabled`, `ReasonOverride`,
  and `ReasonRollout` constants.
- `Flag` — a struct with `Enabled bool`, `Percentage int`, and
  `Overrides map[string]Variant`.
- `Evaluation` — a struct with `Variant` and `Reason` fields.
- `Evaluator` — holds flags under an `sync.RWMutex`.

Implement:

- `NewEvaluator(flags map[string]Flag) (*Evaluator, error)` — validates
  percentage (0–100) and override variants, then defensively deep-copies the
  flags map.
- `(e *Evaluator) Evaluate(flagKey, userID string) (Evaluation, error)` —
  follows the precedence: validate inputs, unknown flag, disabled, override,
  rollout. Use FNV-1a (`crypto/fnv`) to hash `flagKey+"|"+userID` into a
  bucket 0–99.
- `(e *Evaluator) Replace(flags map[string]Flag) error` — validates the new
  map first; only on success does it swap the internal map under write lock.
  A failed replace leaves the old configuration intact.

## Run

```bash
make 03_se_16_feature_flags
```

## Further practice

- Persistent storage for overrides and flag definitions.
- HTTP management API to create, update, and expire flags.
- Evaluation metrics (count, latency) without leaking user IDs.
- User segmentation on attributes (region, plan, cohort).
- Multi-variant experiments (control / treatment-A / treatment-B).
- Flag dependencies (flag B only evaluates when flag A is active).
- Gradual rollout schedules (0% → 10% → 50% → 100%).
