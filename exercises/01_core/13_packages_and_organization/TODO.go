//go:build todo

package packagesandorganization

// NewCounter creates a new Counter with the given name and an initial count of zero.
func NewCounter(name string) *Counter {
	// TODO: Allocate and return a Counter with the provided name.
	return nil
}

// Counter holds a name and a mutable count.
type Counter struct {
	// TODO: add fields
}

// Increment increases the count by one and returns the new count.
func (c *Counter) Increment() int {
	// TODO: Increment c.count and return it.
	return 0
}

// Decrement decreases the count by n steps.
// If the count would go below zero, it is clamped to zero.
func (c *Counter) Decrement(n int) {
	// TODO: Decrease c.count by n, clamping at zero.
}

// Count returns the current count.
func (c *Counter) Count() int {
	// TODO: Return c.count.
	return 0
}

// Name returns the counter's name.
func (c *Counter) Name() string {
	// TODO: Return c.name.
	return ""
}
