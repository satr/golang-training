# Select Statements

This exercise practices channel multiplexing, cancellation-aware sends and
receives, disabled `select` cases, and ticker lifecycle management.

## How `select` works

`select` waits until at least one channel operation can proceed:

```go
select {
case value, ok := <-ch1:
	// Receive from ch1.
case ch2 <- value:
	// Send to ch2.
case <-ctx.Done():
	// Stop cooperatively.
}
```

Important rules:

- With no ready case, `select` blocks.
- If several cases are ready, one is chosen pseudo-randomly. Source order does
  not define priority, and strict fairness is not guaranteed.
- A `default` case runs immediately when no channel operation is ready. Using
  it inside a loop can create a CPU-intensive busy loop.
- Receiving from a closed channel is always ready and returns the element
  type's zero value with `ok == false`.
- Sending to a closed channel panics.
- Sending or receiving on a nil channel blocks forever. In a `select`, a nil
  channel case is disabled.

Setting a local channel variable to nil is therefore a useful way to remove a
closed channel from future selections:

```go
case value, ok := <-input:
	if !ok {
		input = nil
		continue
	}
	use(value)
```

## Task 1: `Multiplexer`

Implement:

```go
func Multiplexer(ctx context.Context, sources []<-chan int) <-chan int
```

Requirements:

- Merge an arbitrary number of input channels into one output channel.
- Forward each received value exactly once.
- Do not promise ordering between different sources.
- Preserve ordering within an individual source.
- Ignore nil source channels.
- Close the output after every non-nil source closes.
- Close the output promptly when `ctx.Done()` closes.
- Make both input receives and output sends cancellation-aware.
- Never close an input channel; input channels belong to their producers.
- Ensure only one goroutine closes the output.

The number of cases in a normal `select` is fixed at compile time. For a
dynamic source list, common approaches are:

1. Start one forwarding goroutine per source and coordinate completion with a
   `sync.WaitGroup`.
2. Build cases dynamically with `reflect.Select` when a single goroutine is a
   strict requirement.

This exercise uses the first approach because it is type-safe and idiomatic.
Each forwarder must select on cancellation while receiving and while sending.

If a caller stops consuming the output early, it must also cancel the context;
otherwise forwarders can remain blocked trying to deliver values. Upstream
producers should usually observe the same context so they do not remain blocked
after the multiplexer stops receiving.

## Task 2: `TickerMerge`

Implement:

```go
func TickerMerge(
	ctx context.Context,
	in <-chan int,
	interval time.Duration,
) <-chan int
```

Requirements:

- Forward values received from `in`.
- Emit `TickerPulse` for every ticker event the goroutine observes.
- Treat a nil input as ticker-only mode.
- Close the output when `ctx` is cancelled or `in` closes.
- Make output sends cancellation-aware.
- Stop the ticker before the goroutine exits.
- Reject `interval <= 0` synchronously with a clear panic.

The exercise uses a sentinel because its output type is fixed to `int`:

```go
const TickerPulse = -1
```

This means an input value of `-1` cannot be distinguished from a pulse. A
production API should normally return a tagged event type instead:

```go
type Event struct {
	Value int
	Pulse bool
}
```

`time.Ticker` does not queue an unlimited number of ticks. If the receiver is
slow, ticks may be dropped or coalesced. The function therefore promises to
forward ticker events it observes, not every theoretical interval.

Input closure can race with a ticker event. If both are ready, `select` may emit
one final pulse before observing that the input is closed.

## Common mistakes

- Iterating sources and performing a blocking receive from each one in turn.
- Assuming the first textual `select` case has priority.
- Claiming `select` provides a strict starvation guarantee.
- Leaving a closed channel in a `select`, causing a busy loop.
- Using a `default` case when blocking is actually required.
- Forgetting cancellation around an output send.
- Closing the output while another goroutine may still send to it.
- Forgetting to call `Ticker.Stop`.
- Generating pulse values from wall-clock timestamps, making tests unstable.
- Starting a ticker with a non-positive interval.

## Tests

Tests should cover:

- Zero, one, and several sources.
- Empty, nil, idle, and duplicate-valued sources.
- A ready source that appears after an idle source in the slice.
- Cancellation while forwarders wait for input.
- Input forwarding without an intervening tick.
- Ticker pulse emission.
- Cancellation while the ticker/input loop is blocked.
- Non-positive interval validation.
- Timely output closure.

Run the normal suite and race detector:

```bash
go test ./...
go test -race ./...
```

If using the repository Make target:

```bash
make 02_concurrency_06_select_statements
```

## Further practice

- Implement the multiplexer with `reflect.Select` and compare tradeoffs.
- Replace the ticker sentinel with a tagged generic event type.
- Implement a non-blocking send and measure how many values are dropped.
- Add per-source buffering and discuss backpressure.
- Build a broadcast fan-out and define its slow-consumer policy.
