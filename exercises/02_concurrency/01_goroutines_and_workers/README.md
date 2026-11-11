# Goroutines and Workers

This exercise practices goroutine lifecycle management, result collection,
ordering, and single-owner state.

## Goroutine fundamentals

The `go` statement starts a function concurrently and returns immediately:

```go
go doWork()
```

A goroutine has no built-in join handle. If its caller needs to know when work
has finished, the program must synchronize explicitly using a channel,
`sync.WaitGroup`, or a higher-level abstraction.

```go
var wg sync.WaitGroup
wg.Add(1) // Add before starting the goroutine.

go func() {
	defer wg.Done()
	doWork()
}()

wg.Wait()
```

Do not call `Add` from inside the new goroutine. `Wait` could observe a zero
counter before that goroutine gets scheduled.

Channels can carry both data and synchronization. A successful send happens
before the corresponding receive completes, so receiving a result also makes
the sender's preceding work visible to the receiver.

## Task 1: `ProcessItems`

Implement:

```go
func ProcessItems(items []string) []string
```

Requirements:

- Return `nil` for an empty input.
- Start one goroutine for each item.
- Convert every string with Unicode-aware uppercase rules.
- Collect worker results through a channel.
- Preserve the original input order even when workers finish out of order.
- Wait for every worker before returning.
- Do not modify the input slice.
- Close the results channel only after every possible sender has finished.

Completion order and input order are different concepts. Include the input
index in each result message:

```go
type itemResult struct {
	index int
	value string
}
```

The collector places each result into its indexed position. This preserves
ordering without requiring workers to finish sequentially.

Use `strings.ToUpper` instead of manually modifying bytes. Go strings contain
UTF-8, and a Unicode character may occupy multiple bytes. Converting a ranged
rune to `byte` corrupts non-ASCII text.

One goroutine per item is intentional for this introductory exercise. It is not
appropriate for an unbounded or extremely large input; a bounded worker pool is
the usual production alternative.

## Task 2: `GoroutineCount`

Implement:

```go
func GoroutineCount(n int) int
```

Requirements:

- Return zero when `n <= 0`.
- Start exactly `n` worker goroutines.
- Have each worker send one increment event through a channel.
- Use one reducer goroutine as the sole owner of the counter.
- Wait until all workers have sent their event.
- Close the increment channel only after all workers finish.
- Return the reducer's final count.

The workers do not directly mutate shared state:

```text
workers -> increment channel -> reducer-owned counter
```

Because only the reducer accesses the counter, no mutex or atomic operation is
needed. This is similar to the actor model: state belongs to one execution
context and other goroutines communicate with it using messages.

## Channel ownership

For both tasks:

- Workers send but do not close a shared results channel.
- A coordinator closes the channel after all senders have stopped.
- Receivers range until closure.
- Sending on a closed channel panics.
- A channel should not be closed merely to signal that one sender is finished.

## Common mistakes

- Returning before workers finish.
- Calling `WaitGroup.Add` inside the worker.
- Capturing loop variables incorrectly in code targeting pre-Go 1.22 language
  semantics; passing values as function parameters is explicit and portable.
- Appending concurrent results and assuming their order matches the input.
- Letting several workers close the same channel.
- Writing to one shared counter without synchronization.
- Converting UTF-8 text byte-by-byte.
- Using `time.Sleep` as synchronization in tests.
- Starting one goroutine per item for an unbounded input.

## Tests

Tests should cover:

- Nil, empty, single, and multiple inputs.
- Preservation of input order.
- Unicode uppercase conversion.
- Confirmation that the input slice is unchanged.
- A larger batch, especially under the race detector.
- Negative, zero, one, and many counter workers.

Run the tests normally and with race detection:

```bash
go test ./...
go test -race ./...
```

If using the repository Make target:

```bash
make 02_concurrency_01_goroutines_and_workers
```

## Further practice

- Replace one-goroutine-per-item with a bounded worker pool.
- Add context cancellation and decide how partial results should be represented.
- Return worker errors alongside successful values.
- Compare the reducer pattern with `sync.Mutex` and `atomic.Int64`.
- Benchmark goroutine-per-item versus a fixed-size pool.
