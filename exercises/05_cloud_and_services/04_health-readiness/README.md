# Health and Readiness Probes

## Theory

Cloud platforms (Kubernetes, Cloud Run) use health and readiness probes to determine if a service is alive and accepting traffic. A health probe checks that the process is functional. A readiness probe additionally verifies the service can serve requests — it blocks traffic during startup or maintenance.

**Key concepts**:
- Health vs readiness: a service can be healthy but not yet ready
- Minimum ready time prevents flapping during startup
- Probes are composable: DB, cache, downstream HTTP endpoints all report independently

## Task

Implement a health checker that runs multiple probes, returns the worst status, and separately gates readiness behind a ready flag and minimum ready time.

1. Add `Health()` that runs all probes and returns the worst status
2. Add `Ready()` that blocks traffic until the service is marked ready AND min ready time elapsed
3. Add `ServeHTTP()` that exposes `/health` and `/ready` as HTTP endpoints

## Run

```bash
make 05_cs_health-readiness
```

## Further practice
- Add a timeout to individual probe execution
- Implement a circuit breaker that skips unhealthy probes
- Export probe results as Prometheus metrics
