# Backpressure

This exercise practices bounded buffering, producer throttling, channel
ownership, and cancellation-aware pipelines.

## What backpressure means

Backpressure occurs when downstream capacity limits upstream production. A Go
channel provides this naturally:

```go
ch := make(chan int, 2)

ch <- 1 // buffer: [1]
ch <- 2 // buffer: [1, 2]
ch <- 3 // blocks until a receiver frees one slot
```

An unbuffered channel has capacity zero, so every send waits for a matching
receive. A buffered channel permits a bounded number of sends to get ahead of
the consumer. Once full, it applies the same blocking pressure.

Buffering changes where waiting occurs; it does not increase the consumer's
processing rate. A permanently slower consumer will eventually fill any finite
buffer.

Buffer-size tradeoffs:

- A larger buffer absorbs larger bursts and decouples short scheduling delays.
- A smaller buffer limits memory and propagates downstream slowness sooner.
- An unbounded queue avoids producer blocking only by accepting unbounded memory
  growth and latency.

## Channel ownership

The goroutine or component responsible for producing the complete stream should
normally close the channel. Receivers should not close channels merely to stop
producers: a concurrent send would panic.

Closing a buffered channel does not discard buffered values. Receivers drain
the remaining values and then receive the zero value with `ok == false`.

## Task 1: `ProcessBatch`

Implement:

```go
func ProcessBatch(ch chan<- int, values []int)
```

Requirements:

- Send every value in order.
- Block naturally when the channel cannot accept the next value.
- Return only after every value has been sent.
- Do not close `ch`; the caller may have additional producers.

This API intentionally has no context, timeout, or error return. Sending to a
nil channel or to a full channel with no consumer therefore blocks forever.
A production API that must be cancellable needs a context-aware variant.

## Task 2: `MergeBuffered`

Implement:

```go
func MergeBuffered(
	ctx context.Context,
	ch1, ch2 <-chan int,
	bufSize int,
) <-chan int
```

Requirements:

- Return immediately with an output channel whose capacity is `bufSize`.
- Accept `bufSize == 0` as an unbuffered output.
- Reject a negative buffer size synchronously with a clear panic.
- Merge both inputs without promising inter-input ordering.
- Preserve ordering within each individual input.
- Treat nil inputs as absent.
- Stop listening to an input after it closes by setting its local variable to
  nil.
- Make output sends cancellation-aware because a full buffer can block.
- Close the output after both inputs close or cancellation is observed.
- Never close either input channel.

Do not count closure events while leaving a closed input in `select`. A closed
channel remains permanently ready, so the same channel could be counted more
than once. Disable it instead:

```go
case value, ok := <-ch1:
	if !ok {
		ch1 = nil
		continue
	}
	// Forward value.
```

Cancellation of the merge does not automatically stop upstream producers. If
they may block while sending, they should observe the same context.

The output may contain buffered values when it is closed. A consumer using an
active context can still drain those values; a consumer using the cancelled
context may stop immediately instead.

## Task 3: `Drain`

Implement:

```go
func Drain(ctx context.Context, ch <-chan int) []int
```

Requirements:

- Return all values in receive order when `ch` closes normally.
- Return promptly with partial results when `ctx.Done()` closes.
- Return no values if cancellation is already observable on entry.
- Do not close the input channel.

If a value and cancellation become ready simultaneously, `select` may choose
either. After the initial cancellation check, a value already in flight may
therefore still be appended. Cancellation is cooperative rather than a strict
transactional boundary.

## Common mistakes

- Assuming a buffer removes backpressure instead of delaying it.
- Closing a channel from the receiving side.
- Sending to a nil channel and expecting an error.
- Counting the same closed input repeatedly in `select`.
- Forgetting cancellation around a potentially blocked output send.
- Cancelling a downstream stage while leaving upstream producers blocked.
- Expecting a cancelled drain to return an exact number of partial values.
- Sorting results when a single producer's ordering should be verified.
- Using sleeps as the primary synchronization mechanism in tests.

## Tests

Tests should cover:

- Empty, buffered, unbuffered, and buffer-overflow batches.
- Preservation of `ProcessBatch` send order.
- Confirmation that `ProcessBatch` does not close its channel.
- Zero, positive, and negative merge buffer sizes.
- Empty, nil, single, dual, and duplicate-valued inputs.
- Merge cancellation while inputs remain open.
- Normal drain ordering.
- Cancellation before draining begins.
- Cancellation while a drain is blocked receiving.
- Timely channel closure without leaked test producers.

Run the suite normally and with the race detector:

```bash
go test ./...
go test -race ./...
```

If using the repository Make target:

```bash
make 02_concurrency_07_backpressure
```

## Further practice

- Add `ProcessBatchContext` with cancellation-aware sends.
- Benchmark latency and producer blocking across several buffer sizes.
- Add a slow consumer and record buffer occupancy over time.
- Implement a drop-newest or drop-oldest policy and document the data-loss
  tradeoff.
- Build a token-bucket rate limiter and distinguish rate limiting from
  backpressure.
