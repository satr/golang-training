# Context Propagation

## Theory

`context.Context` carries deadlines, cancellation signals, and request-scoped
values down a call chain. In a service pipeline, the same context flows through
every processing step so that:

- A cancelled upstream request stops all in-flight work.
- Tracing IDs (e.g. `X-Request-ID`) are attached once and visible everywhere.
- Deadlines propagate: no step can outlive the originating request timeout.

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// Propagate ctx to every step.
result, err := pipeline.Run(ctx, payload)
```

`context.WithValue` is for request-scoped data — not for passing optional
arguments, mutable state, or configuration. Prefer it for tracing IDs and
metadata that every step may need.

## Task

Implement the following API in `solution.go`:

- `RequestID(ctx context.Context) string` — extract an existing request ID
  from `ctx.Value(TraceHeader)`, or generate a new one (`req-1`, `req-2`, …)
  using a counter.
- `WithRequestID(ctx context.Context, id string) context.Context` — attach
  an ID to the context.
- `NewPipeline(steps ...Step) *Pipeline` — store steps.
- `(*Pipeline) Run(ctx context.Context, input []byte) ([]byte, error)` — run
  steps in order, passing each step's output to the next. Check `ctx.Done()`
  before every step to stop promptly on cancellation.

## Run

```bash
go test ./exercises/05_cloud_and_services/04_context-propagation/...
```

Or, with the Makefile:

```bash
make 05_cs_04_context_propagation
```

## Further practice

- Add step-level logging (step index, duration, request ID).
- Return per-step metadata alongside the final result.
- Support per-step timeouts in addition to the overall context deadline.
- Compare `context.WithCancelCause` for carrying a cancellation reason.
