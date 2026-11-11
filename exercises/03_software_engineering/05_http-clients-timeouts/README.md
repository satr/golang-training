# Wrapping HTTP clients with timeouts

## Theory

HTTP clients must bound request latency with timeouts. In Go, wrap `http.Client` with `context.Context` and `time.Duration`. Use `context.WithTimeout` or `context.WithDeadline` to create bounded contexts. Always defer cancel. If the caller passes a context without a deadline, apply a client-level default timeout.

This exercise demonstrates defensive client construction: every outgoing request respects a timeout, so a misbehaving server cannot hold resources forever.

## Task

Implement `SafeClient` that wraps `http.Client`. `NewSafeClient` accepts a default timeout. Provide `Get` and `Post` methods that accept a `context.Context`, apply a bounded context when none exists, and delegate to the underlying client. Use `httptest` for testing.

## Run

```
make 03_software_engineering_05_http_clients_timeouts
```

## Further practice

- Add retry logic with exponential backoff.
- Add connection pool configuration.
