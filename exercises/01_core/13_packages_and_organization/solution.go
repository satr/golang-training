package packagesandorganization

// Counter holds a name and a mutable count.
type Counter struct {
	name  string
	count int
}

// NewCounter creates a new Counter with the given name and an initial count of zero.
func NewCounter(name string) *Counter {
	return &Counter{name: name}
}

// Increment increases the count by one and returns the new count.
func (c *Counter) Increment() int {
	c.count++
	return c.count
}

// Decrement decreases the count by n steps.
// If the count would go below zero, it is clamped to zero.
func (c *Counter) Decrement(n int) {
	c.count -= n
	if c.count < 0 {
		c.count = 0
	}
}

// Count returns the current count.
func (c *Counter) Count() int {
	return c.count
}

// Name returns the counter's name.
func (c *Counter) Name() string {
	return c.name
}
