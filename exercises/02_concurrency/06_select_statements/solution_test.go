package selectstatements

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestMultiplexer(t *testing.T) {
	tests := []struct {
		name    string
		sources [][]int
	}{
		{name: "no sources"},
		{name: "one empty source", sources: [][]int{{}}},
		{name: "one source", sources: [][]int{{1, 2, 3}}},
		{name: "two sources", sources: [][]int{{1, 2}, {3, 4, 5}}},
		{name: "three sources", sources: [][]int{{1}, {2}, {3}}},
		{name: "empty and non-empty", sources: [][]int{{}, {1, 2}}},
		{name: "duplicates", sources: [][]int{{2, 2}, {2, 3}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sources := make([]<-chan int, 0, len(test.sources))
			var want []int
			for _, values := range test.sources {
				sources = append(sources, channelOf(values...))
				want = append(want, values...)
			}

			got := collectUntilClosed(t, Multiplexer(context.Background(), sources))
			slices.Sort(got)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Fatalf("Multiplexer() = %v; want %v", got, want)
			}
		})
	}
}

func TestMultiplexerDoesNotLetIdleSourceBlockReadySource(t *testing.T) {
	idle := make(chan int)
	ready := channelOf(7, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := Multiplexer(ctx, []<-chan int{idle, ready})

	for _, want := range []int{7, 8} {
		select {
		case got := <-out:
			if got != want {
				t.Fatalf("received %d; want %d", got, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("ready source was blocked while waiting for %d", want)
		}
	}

	cancel()
	collectUntilClosed(t, out)
}

func TestMultiplexerIgnoresNilSources(t *testing.T) {
	sources := []<-chan int{nil, channelOf(1, 2), nil}

	got := collectUntilClosed(t, Multiplexer(context.Background(), sources))
	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Fatalf("Multiplexer() = %v; want %v", got, want)
	}
}

func TestMultiplexerClosesOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	out := Multiplexer(ctx, []<-chan int{make(chan int), make(chan int)})

	cancel()

	if got := collectUntilClosed(t, out); len(got) != 0 {
		t.Fatalf("Multiplexer() = %v after cancellation; want no values", got)
	}
}

func TestTickerMergeForwardsInputAndCloses(t *testing.T) {
	in := channelOf(1, 2, 3)

	// The long interval makes this test about input forwarding rather than a
	// race between input closure and a tick.
	got := collectUntilClosed(t, TickerMerge(context.Background(), in, time.Hour))
	want := []int{1, 2, 3}
	if !slices.Equal(got, want) {
		t.Fatalf("TickerMerge() = %v; want %v", got, want)
	}
}

func TestTickerMergeEmitsPulse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	out := TickerMerge(ctx, nil, time.Millisecond)

	select {
	case got := <-out:
		if got != TickerPulse {
			t.Fatalf("TickerMerge() pulse = %d; want %d", got, TickerPulse)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ticker pulse")
	}

	cancel()
	for _, value := range collectUntilClosed(t, out) {
		if value != TickerPulse {
			t.Fatalf("TickerMerge() emitted %d; want only %d pulses", value, TickerPulse)
		}
	}
}

func TestTickerMergeClosesOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	out := TickerMerge(ctx, make(chan int), time.Hour)

	cancel()

	if got := collectUntilClosed(t, out); len(got) != 0 {
		t.Fatalf("TickerMerge() = %v after cancellation; want no values", got)
	}
}

func TestTickerMergeRejectsNonPositiveInterval(t *testing.T) {
	for _, interval := range []time.Duration{0, -time.Second} {
		t.Run(interval.String(), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("TickerMerge did not panic for a non-positive interval")
				}
			}()

			TickerMerge(context.Background(), nil, interval)
		})
	}
}

func channelOf(values ...int) <-chan int {
	ch := make(chan int, len(values))
	for _, value := range values {
		ch <- value
	}
	close(ch)
	return ch
}

func collectUntilClosed(t *testing.T, ch <-chan int) []int {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	var values []int
	for {
		select {
		case value, ok := <-ch:
			if !ok {
				return values
			}
			values = append(values, value)
		case <-timer.C:
			t.Fatal("timed out waiting for output channel to close")
		}
	}
}
