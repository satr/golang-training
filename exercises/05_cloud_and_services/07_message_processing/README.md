# Message Processing Queue

This exercise builds a message queue that dispatches work to a fixed pool of
goroutine workers using typed channels.

## Theory

Go channels provide typed, synchronised communication between goroutines. A
buffered channel decouples producers from consumers up to its capacity. When the
buffer is full, a send blocks until a consumer receives.

A message queue pairs:

- an **input** channel that accepts messages from callers, and
- an **infinite** pool of workers that each range over the input channel,
  call a user-supplied processor, and forward results.

Closing the input channel is the standard mechanism for telling workers that no
more work is coming. Workers terminate their range loop, and a WaitGroup
coordinator can collect results.

`context.Context` lets callers cancel in-flight work. A blocking `Enqueue`
should honour context cancellation so that the sender does not block forever
when the queue buffer is full.

## Task

Implement a message-processing queue in `TODO.go`.

### Data types

```go
type Message struct {
	ID      string
	Type    string
	Payload []byte
}

type Processor interface {
	Process(ctx context.Context, msg Message) error
}

type Result struct {
	MessageID string
	Err       error
}

type Queue struct {
	size      int
	processor Processor
	results   chan Result
	in        chan Message
}
```

### Functions to implement

1. `NewQueue(size int, proc Processor) *Queue` — allocate the struct with
   buffered channels of length `size`.
2. `(q *Queue) Start(ctx context.Context)` — spawn `size` worker goroutines.
   Each worker ranges over `q.in`, calls `q.processor.Process(ctx, msg)`,
   and sends a `Result` to `q.results`.
3. `(q *Queue) Enqueue(ctx context.Context, msg Message) error` — send `msg`
   to `q.in`, honouring context cancellation.
4. `(q *Queue) Results() <-chan Result` — return the results channel as
   receive-only.
5. `(q *Queue) Stop()` — close `q.in` so workers finish their range loops.

Constraints:

- Standard library only.
- `Stop` does **not** close the results channel — that responsibility can be
  left to a wrapper or a caller-managed WaitGroup (exercise boundary).
- Workers must call the processor with the context supplied to `Start`.

## Run

Test the reference solution:

```bash
go test ./exercises/05_cloud_and_services/08_message_processing
go test -race ./exercises/05_cloud_and_services/08_message_processing
```

Or work on the exercise (implement in `TODO.go`):

```bash
go test -tags todo ./exercises/05_cloud_and_services/08_message_processing
go test -race -tags todo ./exercises/05_cloud_and_services/08_message_processing
```

Or with the repository Make target:

```bash
make 05_cloud_and_services_08_message_processing
```

## Further practice

- Close the results channel once all workers exit (use `sync.WaitGroup`).
- Add per-message deadlines inside the processor.
- Implement back-pressure: drop messages when the queue is saturated instead
  of blocking.
- Add metrics (processed, errors, latency) via a `Counter` interface.
- Support multiple queue instances sharing the same `Processor`.
