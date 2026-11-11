package channelsfaninout

import (
	"slices"
	"testing"
	"time"
)

func TestFanOut(t *testing.T) {
	tests := []struct {
		name    string
		items   []int
		workers int
	}{
		{name: "empty source", workers: 1},
		{name: "single worker", items: []int{1, 2, 3}, workers: 1},
		{name: "two workers", items: []int{1, 2, 3, 4}, workers: 2},
		{name: "more workers than values", items: []int{1, 2}, workers: 5},
		{name: "negative zero and duplicates", items: []int{-2, 0, 3, -2}, workers: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := collectUntilClosed(t, FanOut(channelOf(test.items...), test.workers))
			want := make([]int, len(test.items))
			for index, value := range test.items {
				want[index] = value * 2
			}

			slices.Sort(got)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Fatalf("FanOut() = %v; want %v", got, want)
			}
		})
	}
}

func TestFanOutNilSourceClosesOutput(t *testing.T) {
	got := collectUntilClosed(t, FanOut(nil, 2))
	if len(got) != 0 {
		t.Fatalf("FanOut(nil, 2) = %v; want no values", got)
	}
}

func TestFanOutRejectsNonPositiveWorkerCount(t *testing.T) {
	for _, workers := range []int{0, -1} {
		t.Run(workerTestName(workers), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("FanOut did not panic for %d workers", workers)
				}
			}()

			FanOut(channelOf(), workers)
		})
	}
}

func TestFanIn(t *testing.T) {
	tests := []struct {
		name   string
		inputs [][]int
	}{
		{name: "no channels"},
		{name: "one empty channel", inputs: [][]int{{}}},
		{name: "one channel", inputs: [][]int{{1, 2, 3}}},
		{name: "two channels", inputs: [][]int{{1, 2}, {3, 4}}},
		{name: "three channels", inputs: [][]int{{1}, {2}, {3}}},
		{name: "duplicates", inputs: [][]int{{2, 2}, {2, 3}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			channels := make([]<-chan int, 0, len(test.inputs))
			var want []int
			for _, values := range test.inputs {
				channels = append(channels, channelOf(values...))
				want = append(want, values...)
			}

			got := collectUntilClosed(t, FanIn(channels...))
			slices.Sort(got)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Fatalf("FanIn() = %v; want %v", got, want)
			}
		})
	}
}

func TestFanInIgnoresNilChannels(t *testing.T) {
	got := collectUntilClosed(t, FanIn(nil, channelOf(1, 2), nil))
	want := []int{1, 2}
	if !slices.Equal(got, want) {
		t.Fatalf("FanIn() = %v; want %v", got, want)
	}
}

func TestFanInIdleInputDoesNotBlockReadyInput(t *testing.T) {
	idle := make(chan int)
	out := FanIn(idle, channelOf(7, 8))

	for _, want := range []int{7, 8} {
		select {
		case got, ok := <-out:
			if !ok {
				t.Fatal("FanIn closed before the ready input was drained")
			}
			if got != want {
				t.Fatalf("FanIn() produced %d; want %d", got, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("idle input blocked ready value %d", want)
		}
	}

	close(idle)
	if remaining := collectUntilClosed(t, out); len(remaining) != 0 {
		t.Fatalf("FanIn() produced unexpected values %v", remaining)
	}
}

func channelOf(values ...int) <-chan int {
	channel := make(chan int, len(values))
	for _, value := range values {
		channel <- value
	}
	close(channel)
	return channel
}

func collectUntilClosed(t *testing.T, channel <-chan int) []int {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	var values []int
	for {
		select {
		case value, ok := <-channel:
			if !ok {
				return values
			}
			values = append(values, value)
		case <-timer.C:
			t.Fatal("timed out waiting for channel to close")
		}
	}
}

func workerTestName(workers int) string {
	if workers == 0 {
		return "zero"
	}
	return "negative"
}
