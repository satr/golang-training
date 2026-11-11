//go:build todo

package contextpropagation

import (
	"context"
	"fmt"
	"sync/atomic"
)

// TraceHeader is the context key for request ID.
const TraceHeader = "X-Request-ID"

// idGen generates deterministic request IDs for testing.
var idGen atomic.Int64

// RequestID extracts or generates a request ID from context.
func RequestID(ctx context.Context) string {
	// TODO: Check ctx.Value(TraceHeader) for an existing ID (string).
	// If present, return it. Otherwise, generate a new ID using
	// fmt.Sprintf("req-%d", idGen.Add(1)) and return it.
	return ""
}

// WithRequestID returns a context carrying a request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	// TODO: Return context.WithValue(ctx, TraceHeader, id).
	return ctx
}

// Pipeline represents a series of processing steps.
type Pipeline struct {
	steps []Step
}

// Step transforms input data. It must observe ctx.Done() to stop promptly
// when the upstream caller cancels the request.
type Step func(context.Context, []byte) ([]byte, error)

// NewPipeline creates a pipeline from the given steps.
func NewPipeline(steps ...Step) *Pipeline {
	// TODO: Return a *Pipeline containing the provided steps.
	return nil
}

// Run executes every step in order, passing the output of one step as
// input to the next. If ctx is cancelled mid-pipeline, Run stops
// immediately and returns context.Canceled.
func (p *Pipeline) Run(ctx context.Context, input []byte) ([]byte, error) {
	// TODO: Iterate over p.steps. Before each step, select on ctx.Done()
	// to detect cancellation. Call each step with (ctx, data) where data
	// starts as input and becomes the previous step's output.
	// Return the final output, or ctx.Err() if cancelled.
	return nil, nil
}
