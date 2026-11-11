//go:build todo

package goroutinesworkers

// ProcessItems starts one goroutine per item, converts each string to uppercase,
// and returns results in input order.
func ProcessItems(items []string) []string {
	// TODO: Create a results slice, spawn one goroutine per item,
	// uppercase each in its goroutine, and collect results in order
	// using a done channel.
	return nil
}

// GoroutineCount starts n workers. Each worker sends one increment event to a
// reducer goroutine, which is the sole owner of the counter.
func GoroutineCount(n int) int {
	// TODO: Spawn N goroutines that each send 1 to a channel,
	// then receive N values and sum them.
	return 0
}
