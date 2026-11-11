# Context Cancellation

This exercise practices cooperative cancellation across goroutines and channel
operations using Go's `context` package.

## Core idea

`context.Context` carries cancellation, deadlines, and request-scoped values
down a call chain. Cancellation is exposed through `Done()`:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

select {
case <-ctx.Done():
	return ctx.Err()
case value := <-work:
	return process(value)
}
```

Calling `cancel()` closes `ctx.Done()`. It does not terminate goroutines,
interrupt arbitrary code, close application channels, or cancel a blocked
channel operation automatically. Each operation must cooperate by selecting on
`ctx.Done()`.

After cancellation, `ctx.Err()` reports:

- `context.Canceled` when a cancel function was called.
- `context.DeadlineExceeded` when a deadline expired.

Always call the returned cancel function, normally with `defer`, so the context
can release its resources:

```go
ctx, cancel := context.WithTimeout(parent, 500*time.Millisecond)
defer cancel()
```

Do not pass a nil `Context`; use `context.Background()` or `context.TODO()` when
no caller context is available.

## Task 1: `LongRunningWork`

Implement:

```go
func LongRunningWork(ctx context.Context, n int) []int
```

Requirements:

- Produce the sequence `[0, 1, ..., n-1]`.
- Return `nil` when `n <= 0`.
- Check `ctx.Done()` during every iteration.
- Stop cooperatively and return the completed prefix after cancellation.
- Avoid eagerly allocating memory for all `n` values when cancellation may
  occur early.

Cancellation can happen immediately after a check, so one additional iteration
may complete concurrently with cancellation. Context cancellation is a request
to stop promptly, not a transactional boundary.

## Task 2: `FanInWithContext`

Implement:

```go
func FanInWithContext(
	ctx context.Context,
	in1, in2 <-chan int,
) <-chan int
```

Requirements:

- Merge values from both input channels into one output channel.
- Preserve every value while the context remains active, without promising
  ordering between inputs.
- Stop listening to an input after that input closes.
- Treat a nil input channel as absent.
- Stop promptly when `ctx.Done()` closes.
- Make both input receives and output sends cancellation-aware.
- Close the output exactly once, after both inputs close or cancellation wins.
- Never close either input channel; they belong to their producers.

The output sender owns the output channel and should close it:

```go
out := make(chan int)
go func() {
	defer close(out)
	// Merge values or return after cancellation.
}()
```

When an input closes, receiving from it remains immediately ready and returns
the element type's zero value. Set the local channel variable to `nil` to
disable that `select` case:

```go
case value, ok := <-in1:
	if !ok {
		in1 = nil
		continue
	}
	// Forward value.
```

A send must also cooperate with cancellation:

```go
select {
case out <- value:
	// Delivered.
case <-ctx.Done():
	return
}
```

Closing another channel does not cancel a blocked send. Closing `out` while a
goroutine may still send to it would instead cause `panic: send on closed
channel`.

## `select` and cancellation races

Case order does not define priority. If a channel operation and
`ctx.Done()` are both ready, `select` chooses one ready case pseudo-randomly.
The fan-in may therefore forward one final value while cancellation is racing
with normal work. Tests must not assume an exact partial-result count.

## Common mistakes

- Assuming `cancel()` forcibly stops a goroutine.
- Checking cancellation only before a blocking receive or send.
- Using `for value := range input` when the receive must be cancellable.
- Closing input channels from the consumer side.
- Reading repeatedly from a closed channel without checking `ok`.
- Leaving a closed input in a `select`, creating a busy loop.
- Returning without closing the owned output channel.
- Testing concurrent cancellation with sleeps and exact result counts.
- Using `context.WithValue` for optional function arguments or mutable state.

## Tests

Tests should verify:

- Non-positive and normal work sizes.
- Cancellation before work begins.
- Deterministic cancellation during work.
- Empty, nil, single, and dual input channels.
- Duplicate values and nondeterministic merge order.
- Cancellation while waiting for input.
- Cancellation while an output send is pending.
- Timely output closure without goroutine leaks or channel panics.

Run both the normal test suite and the race detector:

```bash
go test ./...
go test -race ./...
```

If using the repository Make target:

```bash
make 02_concurrency_04_context_cancellation
```

## Further practice

- Generalize fan-in to accept any number of channels.
- Make the fan-in generic over the channel element type.
- Return `context.Cause(ctx)` to explain why processing stopped.
- Give upstream producers the same context so they do not remain blocked after
  the fan-in stops receiving.
- Compare `WithCancel`, `WithTimeout`, `WithDeadline`, and `WithCancelCause`.
