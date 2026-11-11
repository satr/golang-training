# Graceful Shutdown

## Theory

Graceful shutdown ensures that in-flight requests complete before the server
exits. Use `context.Context` to signal shutdown, `http.Server.Shutdown()` for
clean HTTP termination, and `os.Signal` to catch SIGINT/SIGTERM.

A lifecycle wrapper accepts:
- a `shutdownTimeout` — how long to wait before forcing shutdown
- `hooks` — ordered callbacks for cleanup (close DB pools, flush logs, etc.)

The `started` channel signals when the listener is bound so tests can reach
the server programmatically.

## Task

Copy `solution.go` from `TODO.go`, then implement:

1. `NewServer` — initialize defaults (5 s shutdown timeout, empty hooks, created `started` chan).
2. `WithShutdownTimeout` — return a functional Option.
3. `WithHook` — return a functional Option that appends to the hook slice.
4. `ListenAndServe` — start the server in a goroutine, signal `started`, wait on `ctx.Done()`, then shutdown and run hooks in order.

## Run

```bash
make 03_se_10_graceful_shutdown
```

## Further practice

Add drain hooks, a `CloseIdleConnections` call before `Shutdown`, and
`os/signal` handling for SIGINT/SIGTERM in production wrappers. Measure how
quickly in-flight requests complete with varying timeouts.
