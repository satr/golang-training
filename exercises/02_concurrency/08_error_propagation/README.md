# Error Propagation from Concurrent Work

This exercise combines bounded workers, context cancellation, ordered result
collection, sentinel errors, and safe goroutine shutdown.

## Why concurrent errors need coordination

A goroutine cannot return an error directly to the function that launched it:

```go
go func() error { // The caller has no result handle.
	return doWork()
}()
```

Concurrent code normally sends a result envelope through a channel:

```go
type result struct {
	index int
	value int
	err   error
}
```

The coordinator receives results, records the first observed error, cancels
sibling work, and continues coordinating shutdown until every internal sender
has stopped. Returning immediately on the first error can leak goroutines that
are blocked on task receives or result sends.

## Sentinel errors and wrapping

The exercise defines stable error categories:

```go
var ErrNegativeValue = errors.New("negative value")
```

Add useful detail while preserving the category with `%w`:

```go
return 0, fmt.Errorf("%w: %d", ErrNegativeValue, value)
```

Callers inspect the chain with `errors.Is`:

```go
if errors.Is(err, ErrNegativeValue) {
	// Handle this category.
}
```

`ProcessConcurrent` wraps the processor error again with its input index.
`errors.Is` still reaches the original sentinel through both wrapping layers.

## Task 1: `Square`

Implement:

```go
func Square(ctx context.Context, value int) (int, error)
```

Requirements:

- Return `ctx.Err()` if cancellation is already observable.
- Reject negative values with an error wrapping `ErrNegativeValue`.
- Return the square and `nil` for non-negative values.

The context check intentionally has precedence over input validation. For an
already-cancelled call with a negative value, return the context error.

## Task 2: `ProcessConcurrent`

Implement:

```go
type Processor func(context.Context, int) (int, error)

func ProcessConcurrent(
	ctx context.Context,
	inputs []int,
	workers int,
	process Processor,
) ([]int, error)
```

Requirements:

- Return an error wrapping `ErrInvalidWorkerCount` when `workers <= 0`.
- Return `ErrNilProcessor` when `process` is nil.
- Return `ctx.Err()` when the parent context is already done.
- Return `nil, nil` for empty input after argument and context validation.
- Use at most one worker per input when `workers > len(inputs)`.
- Preserve input order on success even when completion order differs.
- On the first observed processor error:
  - wrap it with the failing input index;
  - cancel sibling work;
  - stop scheduling new tasks;
  - wait for internal workers to exit;
  - return `nil` results and the wrapped error.
- Propagate parent cancellation.
- Make task sends, task receives, and result sends cancellation-aware.
- Close each owned channel exactly once.

## Cancellation is cooperative

The worker pool can signal cancellation, but it cannot forcibly terminate a
processor. A `Processor` must observe its context during blocking or expensive
work:

```go
select {
case value := <-input:
	return value, nil
case <-ctx.Done():
	return 0, ctx.Err()
}
```

If a processor ignores the context and blocks forever, `ProcessConcurrent`
cannot safely wait for complete shutdown. Panics are also outside this
exercise's error contract; an unrecovered panic still terminates the program.

## First observed error

Workers finish nondeterministically. If several processors fail concurrently,
the API returns the first error observed by the coordinator, not necessarily
the error belonging to the lowest input index. Tests should use one independent
failure unless they intentionally accept any valid error.

The function returns no partial slice on failure. This avoids presenting
zero-value entries as successful results and gives the API a clear
all-or-nothing result contract. Internal successful results may still have been
computed before cancellation.

## Safe lifecycle

The pipeline has four roles:

1. A feeder sends indexed tasks and owns closing the task channel.
2. Workers call the processor and send indexed result envelopes.
3. A coordinator waits for workers and owns closing the result channel.
4. The collector stores successful values by index and triggers cancellation
   after the first processing error.

The collector continues ranging until the result channel closes. This is what
turns error handling into a clean shutdown rather than an early-return leak.

## Exercise template workflow

`TODO.go` is a reset template and is excluded from normal builds with:

```go
//go:build todo
```

Do not enable that tag directly while the active `solution.go` is present,
because both files declare the same API. Use the repository reset command,
which replaces the active solution from its template:

```bash
make reset EXERCISE=exercises/02_concurrency/08_error_propagation
```

Then edit `solution.go` and run the focused target:

```bash
make 02_concurrency_08_error_propagation
```

Run repository-wide tests normally:

```bash
go test ./...
go test -race ./...
```

## Common mistakes

- Trying to return an error directly from a `go` statement.
- Sending only values and losing the associated input index.
- Returning immediately after receiving an error.
- Cancelling before the failing error has been delivered to the coordinator.
- Forgetting cancellation around a blocked task or result channel operation.
- Closing the result channel from one worker while siblings may still send.
- Returning partially filled results without defining their meaning.
- Comparing wrapped errors with `==` instead of `errors.Is`.
- Assuming concurrent failures have deterministic order.
- Expecting context cancellation to stop a processor that ignores its context.
- Testing cancellation with sleeps instead of barriers and completion signals.

## Tests

Tests should verify:

- Successful squaring and negative-value errors.
- Context precedence in `Square`.
- Empty, single, multiple, and excess-worker inputs.
- Input-order preservation despite out-of-order completion.
- Invalid worker counts and a nil processor.
- Wrapped error identity and failing-index context.
- Cancellation of a known blocked sibling after another worker fails.
- Parent cancellation propagation and complete worker shutdown.

## Further practice

- Return partial successes with explicit per-index status.
- Collect all errors instead of stopping after the first.
- Compare this implementation with `errgroup.Group`.
- Generalize `Processor` over input and output types.
- Add panic recovery at the worker boundary and decide how to represent it.
- Preserve the lowest-indexed error instead of the first observed error.
