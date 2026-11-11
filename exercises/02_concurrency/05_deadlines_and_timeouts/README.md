# Deadlines and Timeouts

This exercise practices bounding operation duration with `context.Context` and
implementing cancellation-aware waits without leaking timer resources.

## Timeout versus deadline

Both functions create a child context and return a cancel function:

```go
ctx, cancel := context.WithTimeout(parent, 500*time.Millisecond)
defer cancel()
```

```go
ctx, cancel := context.WithDeadline(parent, time.Now().Add(500*time.Millisecond))
defer cancel()
```

- `WithTimeout` uses a relative duration from now.
- `WithDeadline` uses an absolute point in time.
- A child cannot extend its parent's lifetime. The earliest parent or child
  deadline wins.
- Calling `cancel()` before the deadline also closes `Done()`.

Always call the returned cancel function, even when the operation finishes
successfully. This releases resources associated with the context and its
timer.

## Observing completion

Cancellation and deadlines are exposed through the same channel:

```go
select {
case <-ctx.Done():
	return ctx.Err()
case result := <-results:
	return result
}
```

After `Done()` closes, `Err()` reports why:

- `context.Canceled`: a cancel function was called or an ancestor was
  cancelled.
- `context.DeadlineExceeded`: the effective deadline expired.

Use `errors.Is` when checking returned errors, particularly if another layer
may wrap the context error:

```go
if errors.Is(err, context.DeadlineExceeded) {
	// Handle timeout.
}
```

A deadline is cooperative. It does not terminate a goroutine or interrupt
arbitrary computation by itself; the operation must observe `ctx.Done()`.

The effective deadline can be inspected with:

```go
deadline, ok := ctx.Deadline()
if ok {
	remaining := time.Until(deadline)
	_ = remaining
}
```

## Task 1: `BoundedCompute`

Implement:

```go
func BoundedCompute(ctx context.Context, inputs []int) ([]int, error)
```

Requirements:

- Square each input in order.
- Return all results and `nil` when computation completes.
- Check for cancellation before allocating or starting work.
- Check `ctx.Done()` before each computation.
- On cancellation, return the completed prefix and `ctx.Err()`.
- Avoid eagerly reserving space for every input when cancellation may occur
  early.
- Do not accept a nil context; callers should use `context.Background()` or
  `context.TODO()` when no propagated context exists.

Cancellation may occur immediately after a check, so one additional operation
can complete concurrently with cancellation. A context deadline is not a
transactional boundary.

## Task 2: `SafeSleep`

Implement:

```go
func SafeSleep(ctx context.Context, d time.Duration) error
```

Requirements:

- Return `ctx.Err()` immediately if the context is already done.
- Otherwise return `nil` immediately when `d <= 0`.
- For a positive duration, wait until either the timer fires or the context is
  done.
- Return `nil` when the timer wins.
- Return `ctx.Err()` when cancellation or the deadline wins.
- Stop the timer before returning so its resources can be released promptly.

Use an explicitly owned timer:

```go
timer := time.NewTimer(d)
defer timer.Stop()

select {
case <-timer.C:
	return nil
case <-ctx.Done():
	return ctx.Err()
}
```

If the timer and context become ready simultaneously, `select` may choose
either case. Case order does not establish priority. The explicit check on
entry only gives precedence to cancellation that is already observable.

## Common mistakes

- Forgetting to call the cancel function returned by `WithTimeout` or
  `WithDeadline`.
- Assuming a deadline forcibly stops a goroutine.
- Comparing context errors only with `==` after another layer wraps them.
- Calling `time.Sleep`, which cannot observe context cancellation.
- Creating timers repeatedly in a loop without stopping or reusing them.
- Treating `Canceled` and `DeadlineExceeded` as the same error.
- Writing tests that depend on an exact scheduler race or a tiny sleep.
- Discarding useful partial results when the API promises to return them.

## Tests

Tests should cover:

- Empty, single, negative, duplicate, and multiple input values.
- A context cancelled before computation starts.
- An already-expired deadline.
- Deterministic mid-computation cancellation and partial results.
- Negative, zero, and positive sleep durations.
- Cancellation precedence on entry.
- Early return when a deadline expires before the requested sleep.
- Error identity using `errors.Is`.

Run the suite normally and with the race detector:

```bash
go test ./...
go test -race ./...
```

If using the repository Make target:

```bash
make 02_concurrency_05_deadlines_and_timeouts
```

## Further practice

- Divide one parent deadline across several dependent operations.
- Retry with backoff until success or cancellation.
- Use `context.WithTimeoutCause` and inspect `context.Cause(ctx)`.
- Compare context deadlines with HTTP client and transport-level timeouts.
- Reuse a timer safely inside a loop with `Stop`, draining, and `Reset`.
