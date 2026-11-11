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
	s := &Server{
		srv:             &http.Server{Addr: addr, Handler: handler},
		shutdownTimeout: 5 * time.Second,
		started:         make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithShutdownTimeout sets the grace period for shutdown.
func WithShutdownTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = d
	}
}

// WithHook adds a lifecycle hook to run on shutdown.
func WithHook(hook func(context.Context) error) Option {
	return func(s *Server) {
		s.hooks = append(s.hooks, hook)
	}
}

// ListenAndServe starts the HTTP server and waits for context cancellation.
func (s *Server) ListenAndServe(ctx context.Context) error {
	go func() {
		_ = s.srv.ListenAndServe()
	}()

	close(s.started)
	<-ctx.Done()

	shCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	_ = s.srv.Shutdown(shCtx)
	for _, hook := range s.hooks {
		_ = hook(shCtx)
	}
	return ctx.Err()
}
