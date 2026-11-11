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
	if id, ok := ctx.Value(TraceHeader).(string); ok {
		return id
	}
	return fmt.Sprintf("req-%d", idGen.Add(1))
}

// WithRequestID returns a context carrying a request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, TraceHeader, id)
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
	return &Pipeline{steps: steps}
}

// Run executes every step in order, passing the output of one step as
// input to the next. If ctx is cancelled mid-pipeline, Run stops
// immediately and returns context.Canceled.
func (p *Pipeline) Run(ctx context.Context, input []byte) ([]byte, error) {
	data := input
	for _, step := range p.steps {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		out, err := step(ctx, data)
		if err != nil {
			return nil, err
		}
		data = out
	}
	return data, nil
}
