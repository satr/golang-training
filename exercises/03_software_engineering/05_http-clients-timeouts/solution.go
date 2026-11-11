package httpclientstimeouts

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	return &SafeClient{
		client:         &http.Client{Timeout: timeout},
		defaultTimeout: timeout,
	}
}

// Get performs an HTTP GET with a bounded context.
func (c *SafeClient) Get(ctx context.Context, url string) (*http.Response, error) {
	ctx = c.withTimeout(ctx)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	return c.client.Do(req)
}

// Post performs an HTTP POST with body and content type.
func (c *SafeClient) Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	ctx = c.withTimeout(ctx)
	if body == nil {
		body = http.NoBody
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return c.client.Do(req)
}

func (c *SafeClient) withTimeout(ctx context.Context) context.Context {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		ctx, _ = context.WithTimeout(ctx, c.defaultTimeout)
	}
	return ctx
}
