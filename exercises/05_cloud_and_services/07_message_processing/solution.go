//go:build !todo

package messageprocessing

import "context"

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
func NewQueue(size int, proc Processor) *Queue {
	return &Queue{
		size:      size,
		processor: proc,
		results:   make(chan Result, size),
		in:        make(chan Message, size),
	}
}

// Start launches `size` worker goroutines. Each worker ranges over the input
// channel, calls the processor, and sends the result to the results channel.
func (q *Queue) Start(ctx context.Context) {
	for range q.size {
		go func() {
			for msg := range q.in {
				err := q.processor.Process(ctx, msg)
				q.results <- Result{
					MessageID: msg.ID,
					Err:       err,
				}
			}
		}()
	}
}

// Enqueue sends a message into the queue. It blocks when the buffer is full
// and returns an error when the context is cancelled or the queue has been
// stopped (the input channel is closed for send).
func (q *Queue) Enqueue(ctx context.Context, msg Message) error {
	select {
	case q.in <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Results returns a receive-only channel that yields processed results.
func (q *Queue) Results() <-chan Result {
	return q.results
}

// Stop closes the input channel so that all workers finish their range loops
// and eventually exit. After Stop, no more messages may be enqueued.
func (q *Queue) Stop() {
	close(q.in)
}
