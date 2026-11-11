# Retries and Retry Budgets

## Theory

Network calls fail. Transient errors require retries, but blind retries cause amplified load (thundering herd). A retry budget limits total retries within a time window regardless of individual requests, preventing cascade failures.

**Key concepts**:
- Exponential backoff space out retries
- Context limits prevent infinite waiting
- Budgets protect downstream services at system level

## Task

Implement `Execute` with exponential backoff and budget checks. Implement `Budget` that tracks system-wide retry attempts within a sliding window.

1. Implement `Execute` — retry failed calls with exponential backoff
2. Implement `Budget.Allow` — cap retries within a time window
3. Wire budget into Execute so system-wide limit stops retrying

## Run

```bash
make 05_cs_retries
```

## Further practice
- Add jitter to backoff to prevent synchronized hammering
- Implement per-endpoint budgets
- Track retry histograms as observability signals
