# Rate Limiting

## Theory

Rate limiting protects services from overload by bounding the number of requests processed in a time window. The token bucket algorithm is simple: start with N tokens, consume one per request, refuse when empty.

**Key concepts**:
- Token bucket offers burst then steady state
- Atomic consumption prevents race conditions
- Rate middleware is transparent to downstream handlers

## Task

Implement a token-bucket limiter and an HTTP middleware that enforces it.

1. Implement `Limiter.Allow` — consume one token atomically
2. Implement `Limiter.AllowN` — atomic batch consume
3. Wire `Middleware` to return 429 when tokens are exhausted

## Run

```bash
make 05_cs_rate-limiting
```

## Further practice
- Add periodic refill via ticker
- Implement per-client rate limits using a map
- Add rate reset at period boundaries (sliding window)
