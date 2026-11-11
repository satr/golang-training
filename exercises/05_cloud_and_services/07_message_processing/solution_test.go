package messageprocessing

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"sync"
	"testing"
)

// nopProcessor records every message ID it processes and never returns an
// error.
type nopProcessor struct {
	mu   sync.Mutex
	ids  []string
	fail map[string]struct{}
}

func (p *nopProcessor) Process(_ context.Context, msg Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.fail[msg.ID]; ok {
		return fmt.Errorf("%s failed", msg.Type)
	}
	p.ids = append(p.ids, msg.ID)
	return nil
}

type waitingProcessor struct {
	started chan struct{}
}

func (p *waitingProcessor) Process(ctx context.Context, _ Message) error {
	select {
	case p.started <- struct{}{}:
	default:
	}
	<-ctx.Done()
	return ctx.Err()
}

func TestQueueProcessesAllMessages(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		count   int
	}{
		{name: "single worker", workers: 1, count: 5},
		{name: "two workers", workers: 2, count: 10},
		{name: "more workers than messages", workers: 10, count: 3},
		{name: "single message single worker", workers: 1, count: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proc := &nopProcessor{}
			q := NewQueue(test.workers, proc)
			ctx := context.Background()
			q.Start(ctx)

			results := make(chan []Result, 1)
			go func() {
				got := make([]Result, test.count)
				for i := range got {
					got[i] = <-q.Results()
				}
				results <- got
			}()

			var wantIDs []string
			// Enqueue all messages.
			for i := 0; i < test.count; i++ {
				msg := Message{
					ID:      fmt.Sprintf("%s-msg-%d", test.name, i),
					Type:    "test",
					Payload: []byte{byte(i)},
				}
				if err := q.Enqueue(ctx, msg); err != nil {
					t.Fatalf("Enqueue: %v", err)
				}
				wantIDs = append(wantIDs, msg.ID)
			}

			q.Stop()
			gotResults := <-results
			if len(gotResults) != test.count {
				t.Fatalf("expected %d results; got %d", test.count, len(gotResults))
			}

			var resultIDs []string
			for _, result := range gotResults {
				if result.Err != nil {
					t.Errorf("processing %q: %v", result.MessageID, result.Err)
				}
				resultIDs = append(resultIDs, result.MessageID)
			}
			sort.Strings(resultIDs)
			sort.Strings(wantIDs)
			if !slices.Equal(resultIDs, wantIDs) {
				t.Errorf("result IDs = %v; want %v", resultIDs, wantIDs)
			}

			proc.mu.Lock()
			gotIDs := append([]string(nil), proc.ids...)
			proc.mu.Unlock()
			sort.Strings(gotIDs)
			if !slices.Equal(gotIDs, wantIDs) {
				t.Errorf("processed IDs = %v; want %v", gotIDs, wantIDs)
			}
		})
	}
}

func TestQueuePropagatesProcessorErrors(t *testing.T) {
	proc := &nopProcessor{
		fail: map[string]struct{}{"bad": {}},
	}
	q := NewQueue(2, proc)
	ctx := context.Background()
	q.Start(ctx)

	for _, id := range []string{"good-1", "bad", "good-2"} {
		if err := q.Enqueue(ctx, Message{ID: id, Type: "echo"}); err != nil {
			t.Fatalf("Enqueue: %v", err)
		}
	}

	q.Stop()
	results := make([]Result, 3)
	for i := range results {
		results[i] = <-q.Results()
	}
	seen := make(map[string]bool, len(results))
	var errorCount int
	for _, result := range results {
		if result.MessageID != "good-1" && result.MessageID != "bad" && result.MessageID != "good-2" {
			t.Errorf("unexpected result for message %q", result.MessageID)
		}
		if seen[result.MessageID] {
			t.Errorf("duplicate result for message %q", result.MessageID)
		}
		seen[result.MessageID] = true
		if result.MessageID == "bad" && result.Err == nil {
			t.Error("expected an error result for bad message")
		}
		if result.MessageID != "bad" && result.Err != nil {
			t.Errorf("unexpected error for %q: %v", result.MessageID, result.Err)
		}
		if result.Err != nil {
			errorCount++
		}
	}
	if errorCount != 1 {
		t.Fatalf("expected 1 error result; got %d", errorCount)
	}
	for _, id := range []string{"good-1", "bad", "good-2"} {
		if !seen[id] {
			t.Errorf("missing result for message %q", id)
		}
	}
}

func TestQueueEnqueueContextCancellation(t *testing.T) {
	proc := &waitingProcessor{started: make(chan struct{}, 1)}
	q := NewQueue(1, proc)
	ctx, cancel := context.WithCancel(context.Background())
	q.Start(ctx)

	if err := q.Enqueue(ctx, Message{ID: "1", Type: "fill"}); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}
	<-proc.started
	if err := q.Enqueue(ctx, Message{ID: "2", Type: "fill"}); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	cancel()
	err := q.Enqueue(ctx, Message{ID: "blocked", Type: "block"})
	if err != context.Canceled {
		t.Fatalf("Enqueue error = %v; want %v", err, context.Canceled)
	}
	q.Stop()
	for range 2 {
		<-q.Results()
	}
}
