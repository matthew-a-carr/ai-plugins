# Pattern — Circuit breaker setup

Anchors: P24, P22, P23

Default circuit breaker and bulkhead configuration for a Spring Boot service using Resilience4j.

## Default thresholds

| Parameter | Default | Rationale |
| --- | --- | --- |
| Failure rate threshold | 50% | Trip when half the calls fail — sensitive enough to catch real outages, not so sensitive that transient blips trip it |
| Sliding window type | Count-based | Simpler to reason about than time-based |
| Sliding window size | 10 calls | Small enough to react quickly, large enough to avoid noise |
| Wait duration in open state | 30 seconds | Gives the downstream time to recover before probing |
| Permitted calls in half-open | 5 | Enough probes to confirm recovery |
| Minimum number of calls | 5 | Don't trip on the first few calls at startup |

Tune these per dependency based on its traffic volume and expected failure profile. These are starting points, not universal answers.

## Resilience4j Spring Boot configuration

```yaml
resilience4j:
  circuitbreaker:
    configs:
      default:
        slidingWindowType: COUNT_BASED
        slidingWindowSize: 10
        failureRateThreshold: 50
        waitDurationInOpenState: 30s
        permittedNumberOfCallsInHalfOpenState: 5
        minimumNumberOfCalls: 5
        registerHealthIndicator: true
    instances:
      pricing:
        baseConfig: default
      payment:
        baseConfig: default
        failureRateThreshold: 30    # stricter for payment — fail fast
        waitDurationInOpenState: 60s

  bulkhead:
    instances:
      pricing:
        maxConcurrentCalls: 20
      payment:
        maxConcurrentCalls: 10

  timelimiter:
    instances:
      pricing:
        timeoutDuration: 3s
      payment:
        timeoutDuration: 5s
```

## Bulkhead configuration

Each downstream dependency gets its own bulkhead (separate concurrency limit or thread pool). This prevents a slow dependency from exhausting the service's capacity and starving other call paths.

```java
@CircuitBreaker(name = "pricing", fallbackMethod = "cachedPrice")
@Bulkhead(name = "pricing")
@TimeLimiter(name = "pricing")
public CompletableFuture<Money> getPrice(ProductId productId) {
    return CompletableFuture.supplyAsync(() -> pricingClient.fetchPrice(productId));
}
```

## Health indicator integration

With `registerHealthIndicator: true`, the circuit breaker state is exposed on the `/actuator/health` endpoint:

```json
{
  "circuitBreakers": {
    "pricing": { "state": "CLOSED", "failureRate": "12.5%" },
    "payment": { "state": "OPEN", "failureRate": "60.0%" }
  }
}
```

Use this in readiness probes only if the service should stop receiving traffic when a critical circuit is open. For non-critical dependencies, report the state but do not affect readiness.

## Metrics to watch

| Metric | What it tells you |
| --- | --- |
| `resilience4j.circuitbreaker.state` | Current state (closed / open / half-open) — alert on open |
| `resilience4j.circuitbreaker.failure.rate` | Failure percentage over the window |
| `resilience4j.circuitbreaker.calls` (tagged by kind) | Successful, failed, not-permitted call counts |
| `resilience4j.bulkhead.available.concurrent.calls` | Remaining capacity — alert when near zero |
| `resilience4j.timelimiter.calls` (tagged by kind) | Successful vs timed-out calls |

## Fail open vs fail closed

| Criteria | Fail open | Fail closed |
| --- | --- | --- |
| A cached or default response is acceptable | ✓ | |
| The dependency is non-critical to the primary flow | ✓ | |
| Serving stale or incorrect data has financial or safety consequences | | ✓ |
| The dependency is required for correctness (e.g., payment authorisation) | | ✓ |

See `graceful-degradation.md` for the fallback implementation.

## References

- Constitution P22 (timeouts), P23 (retries), P24 (circuit breakers / bulkheads)
- `graceful-degradation.md` — fallback strategies
- `timeouts-and-retries.md` — timeout and retry defaults
- Resilience4j: https://resilience4j.readme.io/
