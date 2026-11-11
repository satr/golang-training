//go:build todo

package httpclientstimeouts

import (
	"context"
	"io"
	"time"
)

// SafeClient wraps http.Client and ensures every request has a timeout.
type SafeClient struct {
	client         *http.Client
	defaultTimeout time.Duration
}

// NewSafeClient creates a client that applies defaultTimeout to every request
// when the provided context has no deadline.
func NewSafeClient(timeout time.Duration) *SafeClient {
	// TODO: Create a SafeClient with the given default timeout.
	return &SafeClient{client: &http.Client{}, defaultTimeout: timeout}
}

// Get performs an HTTP GET with a bounded context.
func (c *SafeClient) Get(ctx context.Context, url string) (*http.Response, error) {
	// TODO: Ensure context has a deadline, create request, call client.Do.
	return nil, context.DeadlineExceeded
}

// Post performs an HTTP POST with body and content type.
func (c *SafeClient) Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	// TODO: Ensure context has a deadline, create request with body/contentType.
	return nil, context.DeadlineExceeded
}
