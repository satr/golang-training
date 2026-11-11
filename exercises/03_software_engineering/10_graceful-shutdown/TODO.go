//go:build todo

package gracefulshutdown

import (
	"context"
	"net/http"
	"time"
)

// Server wraps an http.Server and provides graceful shutdown with hooks.
type Server struct {
	srv             *http.Server
	shutdownTimeout time.Duration
	hooks           []func(context.Context) error
	started         chan struct{}
}

// Option configures a Server.
type Option func(*Server)

// NewServer returns a new Server listening on addr with the given handler.
func NewServer(addr string, handler http.Handler, opts ...Option) *Server {
	// TODO: Create and return *Server with reasonable defaults.
	return nil
}

// WithShutdownTimeout sets the grace period for shutdown.
func WithShutdownTimeout(d time.Duration) Option {
	// TODO: Return an Option that sets shutdownTimeout.
	return nil
}

// WithHook adds a lifecycle hook to run on shutdown.
func WithHook(hook func(context.Context) error) Option {
	// TODO: Return an Option that appends a hook.
	return nil
}

// ListenAndServe starts the HTTP server and waits for context cancellation.
func (s *Server) ListenAndServe(ctx context.Context) error {
	// TODO: Start srv in a goroutine, close started, wait ctx.Done().
	// On cancel: create shutdown context with timeout, srv.Shutdown, run hooks.
	return nil
}
