# Channels and Fan-in/Fan-out

This exercise practices channel communication, work distribution, result
merging, goroutine completion, and channel ownership.

## Channel behavior

A channel is a typed synchronization and communication mechanism:

- Sending to an unbuffered channel waits for a matching receiver.
- Receiving from an unbuffered channel waits for a matching sender.
- A buffered send waits only when the buffer is full.
- Receiving from a closed and drained channel returns the zero value with
  `ok == false`.
- Sending to a closed channel panics.
- Sending or receiving on a nil channel blocks forever.

The producer responsible for the complete stream normally owns closing its
channel. Receivers should not close a channel merely to make their receive loop
stop.

## Fan-out versus broadcast

In this exercise, fan-out means several workers compete to receive from one
source channel:

```text
source -> worker 1
       -> worker 2
       -> worker 3
```

Each source value goes to exactly one worker. It is not broadcast to every
worker. Because workers complete independently, result order is not guaranteed.

Fan-in performs the inverse shape: one forwarding goroutine reads each input,
and all forwarders send into a shared output:

```text
input 1 -> forwarder 1 -> output
input 2 -> forwarder 2 -> output
input 3 -> forwarder 3 -> output
```

Ordering is preserved within one input channel, but values from different
inputs can interleave nondeterministically.

## Task 1: `FanOut`

Implement:

```go
func FanOut(source <-chan int, n int) <-chan int
```

Requirements:

- Panic synchronously when `n <= 0`.
- Treat a nil source as empty and return a closed output channel.
- For a non-nil source, start exactly `n` workers.
- Have workers compete for source values.
- Double each value and send it to the output exactly once.
- Close the output only after the source is closed, drained, and all workers
  have stopped sending.
- Do not close the source channel; its producer owns it.

The caller must eventually close the source and continuously drain the returned
channel. Otherwise workers may remain blocked forever.

## Task 2: `FanIn`

Implement:

```go
func FanIn(channels ...<-chan int) <-chan int
```

Requirements:

- Start one forwarding goroutine for every non-nil input.
- Ignore nil inputs.
- Forward every value exactly once.
- Preserve ordering within each individual input.
- Close the output only after all non-nil inputs close and their forwarders
  finish.
- With zero non-nil inputs, return an output channel that closes promptly.
- Never close an input channel.

Each input producer must close its channel. One input that remains open—even if
it sends no values—keeps the output open.

## Coordinated output closure

Several goroutines send to the same output, so none of them can safely close it.
A coordinator waits for every sender:

```go
var wg sync.WaitGroup

// Add workers before starting them.

go func() {
	wg.Wait()
	close(out)
}()
```

Calling `WaitGroup.Add` inside a new goroutine is unsafe because `Wait` may
observe zero before that goroutine is scheduled. Closing the output before
`Wait` returns can cause `panic: send on closed channel`.

These APIs intentionally do not accept a context. A caller that stops consuming
early cannot cancel workers or forwarders; production pipelines normally add
context cancellation to every blocking receive and send.

## Exercise and solution build tags

The starter and reference implementation are mutually exclusive:

```go
// TODO.go
//go:build todo

// solution.go
//go:build !todo
```

Normal tests use the reference solution:

```bash
go test ./...
```

To work on the exercise, implement the functions in `TODO.go` and run:

```bash
go test -tags todo ./...
```

The initial TODO functions panic deliberately. Returning a nil channel would
make range-based tests block forever and hide the missing implementation.

## Common mistakes

- Treating fan-out as broadcast.
- Starting zero workers and leaving the source producer blocked.
- Ranging over a nil input channel.
- Closing an output from one of several sender goroutines.
- Closing an output before all senders finish.
- Calling `WaitGroup.Add` after starting a worker.
- Assuming concurrent result order is deterministic.
- Stopping output consumption before cancelling or otherwise stopping senders.
- Enabling the TODO build tag without excluding the reference solution.

## Tests

Tests should cover:

- Empty input and multiple worker counts.
- More workers than values.
- Negative, zero, and duplicate values.
- Invalid worker counts.
- Nil fan-out source handling.
- Zero, one, and several fan-in inputs.
- Empty and nil fan-in inputs.
- An idle input alongside a ready input.
- Timely output closure.

Run the normal suite and race detector:

```bash
go test ./...
go test -race ./...
```

Run the exercise implementation with:

```bash
go test -tags todo ./...
go test -race -tags todo ./...
```

If using the repository Make target:

```bash
make 02_concurrency_02_channels_and_fan_in_out
```

## Further practice

- Add `context.Context` and make every blocking operation cancellable.
- Add bounded output buffering and observe how it changes backpressure.
- Implement broadcast fan-out where every subscriber receives every value.
- Generalize the functions over any element type.
- Preserve global input order by attaching sequence numbers.
