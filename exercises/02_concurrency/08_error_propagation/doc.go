// Package errorpropagation demonstrates how concurrent workers report errors,
// cancel sibling work, preserve successful result ordering, and shut down
// without leaking goroutines.
package errorpropagation
