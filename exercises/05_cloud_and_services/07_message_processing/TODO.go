//go:build todo

package messageprocessing

// Message represents a single unit of work that flows through the queue.
type Message struct {
	ID      string
	Type    string
	Payload []byte
}

// Processor processes a single message. Implementations must honour context
// cancellation.
type Processor interface {
	Process(ctx context.Context, msg Message) error
}

// Result holds the outcome of processing one message.
type Result struct {
	MessageID string
	Err       error
}

// Queue is a bounded message queue that dispatches messages to workers.
type Queue struct {
	size      int
	processor Processor
	results   chan Result
	in        chan Message
}

// NewQueue creates a new queue with the given buffer size and processor.
//
// TODO: allocate and return a *Queue.
func NewQueue(size int, proc Processor) *Queue {
	panic("TODO: implement NewQueue")
}

// Start launches the worker goroutines. Each worker pulls messages from the
// input channel, calls the processor, and sends the result downstream.
// The workers exit when `Stop` closes the input channel.
//
// Start must be called exactly once before Enqueue.
//
// TODO: spawn `size` workers.
func (q *Queue) Start(ctx context.Context) {
	panic("TODO: implement Start")
}

// Enqueue sends a message into the queue. It blocks when the queue buffer is
// full. Returns an error derived from `ctx` when the context is cancelled or
// the queue has been stopped.
//
// TODO: send `msg` to the internal input channel, respecting context.
func (q *Queue) Enqueue(ctx context.Context, msg Message) error {
	panic("TODO: implement Enqueue")
}

// Results returns a receive-only channel that yields processed results.
// The channel is closed once every worker has exited.
//
// TODO: return q.results as a receive-only channel.
func (q *Queue) Results() <-chan Result {
	panic("TODO: implement Results")
}

// Stop signals all workers to drain the input channel and exit. After Stop,
// no more messages may be enqueued.
//
// TODO: close the input channel so workers finish their range loops.
func (q *Queue) Stop() {
	panic("TODO: implement Stop")
}
