package boundedworkerpools

import (
	"context"
	"sync"
)

// WorkerPool squares inputs using at most workers concurrent workers.
// Results are returned in completion order, which may differ from input order.
func WorkerPool(inputs []int, workers int) []int {
	if len(inputs) == 0 {
		return nil
	}
	workers = normalizedWorkerCount(workers, len(inputs))

	tasks := make(chan int)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for value := range tasks {
				results <- value * value
			}
		}()
	}

	go func() {
		defer close(tasks)
		for _, value := range inputs {
			tasks <- value
		}
	}()

	// results has multiple senders, so a coordinator closes it only after every
	// worker has stopped sending.
	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]int, 0, len(inputs))
	for value := range results {
		out = append(out, value)
	}
	return out
}

// WorkerPoolWithTimeout squares inputs using at most workers concurrent workers.
// It stops scheduling and processing work when ctx is cancelled and returns the
// results delivered before the pipeline stops. Cancellation races with channel
// operations, so a value already in flight may still be included.
func WorkerPoolWithTimeout(ctx context.Context, inputs []int, workers int) []int {
	if len(inputs) == 0 || ctx.Err() != nil {
		return nil
	}
	workers = normalizedWorkerCount(workers, len(inputs))

	tasks := make(chan int)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-tasks:
					if !ok {
						return
					}

					select {
					case results <- value * value:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	// The producer is the sole owner responsible for closing tasks.
	go func() {
		defer close(tasks)
		for _, value := range inputs {
			select {
			case tasks <- value:
			case <-ctx.Done():
				return
			}
		}
	}()

	// No worker can send after Wait returns, so this is the safe place to close
	// results. Ranging below also ensures the function waits for clean shutdown.
	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]int, 0, len(inputs))
	for value := range results {
		out = append(out, value)
	}
	return out
}

func normalizedWorkerCount(workers, inputCount int) int {
	if workers < 1 {
		return 1
	}
	if workers > inputCount {
		return inputCount
	}
	return workers
}
