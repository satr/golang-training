//go:build !todo

package channelsfaninout

import "sync"

// FanOut distributes source values across n workers. Each value is processed
// by exactly one worker, doubled, and sent to the returned channel.
//
// FanOut panics when n <= 0. A nil source is treated as an empty source.
func FanOut(source <-chan int, n int) <-chan int {
	if n <= 0 {
		panic("FanOut: worker count must be positive")
	}

	out := make(chan int)
	if source == nil {
		close(out)
		return out
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			for value := range source {
				out <- value * 2
			}
		}()
	}

	// Only the coordinator closes out, after every sender has finished.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FanIn merges all non-nil input channels into one output channel. It closes
// the output after every input is closed and drained.
func FanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	forward := func(input <-chan int) {
		defer wg.Done()
		for value := range input {
			out <- value
		}
	}

	for _, channel := range channels {
		if channel == nil {
			continue
		}
		wg.Add(1)
		go forward(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
