package goroutinesworkers

import (
	"strings"
	"sync"
)

type itemResult struct {
	index int
	value string
}

// ProcessItems starts one goroutine per item, converts each string to uppercase,
// and returns results in input order.
func ProcessItems(items []string) []string {
	if len(items) == 0 {
		return nil
	}

	completed := make(chan itemResult, len(items))
	var wg sync.WaitGroup

	for index, item := range items {
		wg.Add(1)
		go func(index int, item string) {
			defer wg.Done()
			completed <- itemResult{
				index: index,
				value: strings.ToUpper(item),
			}
		}(index, item)
	}

	// No sender can remain after Wait returns, so this coordinator owns closing
	// the results channel.
	go func() {
		wg.Wait()
		close(completed)
	}()

	results := make([]string, len(items))
	for result := range completed {
		results[result.index] = result.value
	}
	return results
}

// GoroutineCount starts n workers. Each worker sends one increment event to a
// reducer goroutine, which is the sole owner of the counter.
func GoroutineCount(n int) int {
	if n <= 0 {
		return 0
	}

	increments := make(chan struct{})
	total := make(chan int, 1)

	go func() {
		count := 0
		for range increments {
			count++
		}
		total <- count
	}()

	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			increments <- struct{}{}
		}()
	}

	wg.Wait()
	close(increments)
	return <-total
}
