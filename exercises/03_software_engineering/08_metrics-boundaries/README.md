# Metrics and Circuit Breakers

## Theory

Service boundaries need metrics to track latency and error rates. A circuit breaker prevents cascading failures by tripping to "open" after enough consecutive errors, then probing with a half-open request to see if the upstream recovered.

## Task

Implement a Counter that tracks total requests, errors, and provides error rate. Implement a CircuitBreaker that trips open after a configurable number of consecutive failures, stays open for a configurable duration, then allows one half-open probe.

## Run

```bash
make 03_software_engineering_08_metrics_boundaries
```

## Further practice

Add histogram buckets for latency tracking, exponential backoff for recovery.
