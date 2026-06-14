# Pattern — Graceful degradation

Anchors: P25, P24

Define and test how the service behaves when a dependency is unavailable. The goal: the service stays useful, and the customer-visible impact is documented before the incident, not discovered during it.

## Degraded-modes table

Every service with external dependencies maintains a `docs/degraded-modes.md` table:

```markdown
| Dependency | When unavailable | Service behaviour | Customer-visible impact | Tested? |
| --- | --- | --- | --- | --- |
| Pricing API | Circuit open | Return cached price (max 5 min stale) | Prices may be slightly outdated | Yes — chaos test weekly |
| Payment gateway | Circuit open | Queue payment for retry, confirm order as "pending" | Customer sees "payment processing" instead of instant confirmation | Yes — integration test |
| Email service | Timeout / error | Suppress notification, log for retry | Customer does not receive confirmation email immediately | No — manual test only |
| Postgres (primary) | Connection refused | Service returns 503 on all write paths | Full outage on writes; reads via replica if configured | Yes — failover drill quarterly |
```

Update this table whenever a dependency is added or a fallback strategy changes. Review it during incident retrospectives.

## Fallback strategies

| Strategy | Use when | Example |
| --- | --- | --- |
| Cached result | A slightly stale answer is acceptable | Return last-known price from local cache |
| Default response | A safe default exists | Show "standard shipping" when the shipping calculator is down |
| Feature hidden | The feature is non-critical | Hide "recommended products" panel when the recommendation engine is unavailable |
| Queue for retry | The operation must complete but not immediately | Queue the email and retry when the mail service recovers |

## Spring Boot example

```java
@CircuitBreaker(name = "pricing", fallbackMethod = "cachedPrice")
public Money getPrice(ProductId productId) {
    return pricingClient.fetchPrice(productId);
}

private Money cachedPrice(ProductId productId, Throwable t) {
    log.warn("Pricing API unavailable, returning cached price for {}", productId);
    return priceCache.getLatest(productId)
        .orElseThrow(() -> new ServiceUnavailableException("No cached price available"));
}
```

The fallback method receives the same parameters plus the exception. If the fallback itself cannot serve a result, let the failure propagate — do not return invented data.

## Fail open vs fail closed

| Criteria | Fail open (degrade gracefully) | Fail closed (reject the request) |
| --- | --- | --- |
| A safe default or cached value exists | ✓ | |
| The operation is non-critical to the user's primary task | ✓ | |
| Incorrect data is worse than no data | | ✓ |
| The operation involves money, compliance, or safety | | ✓ |
| The dependency outage is expected to be transient | ✓ | |

When in doubt, fail closed and discuss with the team. It is easier to relax a strict stance than to recover from serving bad data.

## Anti-patterns

- **Discovering degraded behaviour during an incident.** If the first time you learn what happens when the payment gateway is down is during a production outage, the degraded-modes table was not maintained or tested.
- **Fallback returns invented data.** A fallback that makes up a price or fabricates a response is worse than an error.
- **No circuit breaker on the fallback path.** If the fallback calls another external service, that service needs its own circuit breaker.

## References

- Constitution P25 (graceful degradation), P24 (circuit breakers / bulkheads)
- `circuit-breaker-setup.md` — default circuit breaker configuration
- `timeouts-and-retries.md` — timeout and retry defaults
