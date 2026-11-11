//go:build todo

package channelsfaninout

// FanOut distributes source values across n workers. Each value must be
// processed by exactly one worker, doubled, and sent to the returned channel.
//
// Requirements:
//   - Panic when n <= 0.
//   - Treat a nil source as empty and return a closed channel.
//   - Start exactly n workers for a non-nil source.
//   - Close the output only after every worker has stopped sending.
func FanOut(source <-chan int, n int) <-chan int {
	panic("TODO: implement FanOut")
}

// FanIn merges all non-nil input channels into one output channel. The output
// must close after every input has been closed and drained.
//
// Requirements:
//   - Return an already-closing output when there are no non-nil inputs.
//   - Start one forwarding goroutine for each non-nil input.
//   - Preserve ordering within each individual input.
//   - Close the output only after every forwarder has stopped sending.
func FanIn(channels ...<-chan int) <-chan int {
	panic("TODO: implement FanIn")
}
