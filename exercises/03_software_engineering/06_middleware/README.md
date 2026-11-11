# HTTP Middleware

## Theory

Middleware in Go wraps HTTP handlers to add cross-cutting behavior. Each middleware takes the form `func(http.Handler) http.Handler`. Chain middleware by composing them, where the outermost middleware runs first and the innermost runs last.

## Task

Implement RequestID, Logging, Recovery, and Chain middleware functions.

## Run

```bash
make 03_software_engineering_06_middleware
```

## Further practice

Add authentication and rate-limiting middleware.