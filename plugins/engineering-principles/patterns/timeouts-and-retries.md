# Pattern — Timeouts, retries, and backoff

Constitution P22, P23, P24 in practice.

## Timeouts

Every external call has:

- **Connect timeout** — how long to wait to establish the connection. Short (1–3s) for in-cluster calls; longer (5–10s) for external internet calls.
- **Read timeout** — how long to wait for a response after sending the request. Sized to the upstream's p99 latency plus headroom, not the maximum imaginable.
- **Overall deadline** — how long this call may take, total. Propagated as a header (`Deadline` / `X-Request-Deadline`) so the upstream knows when to give up.

If the framework defaults to no timeout, set one. "No timeout" is a bug.

## Retry policy

```text
attempt n: delay = min(cap, base * 2^n) + random(0, jitter)
```

- `base = 100ms`, `cap = 5s`, `jitter = 100ms` is a reasonable starting point.
- Cap the number of attempts (3 is often enough). After that, surface the failure.
- Cap the total elapsed time. A 3-attempt retry that takes 30s is worse than failing in 3s.
- Retry only on retryable failures: timeouts, 5xx, connection resets. Never retry on 4xx (the caller is wrong, retrying won't help).
- **Only retry idempotent operations**. Use the operation's idempotency key (Constitution P6) so the upstream can dedupe.

## Circuit breaker

State machine: **Closed → Open → Half-Open → Closed**.

- **Closed**: calls flow through. Count failures.
- Trip to **Open** when failure rate over a sliding window breaches a threshold (e.g. ≥50% failures over 20 calls).
- **Open**: calls fail fast without hitting the downstream.
- After a cool-down (e.g. 30s), move to **Half-Open**: allow a small probe of calls through.
- If probes succeed, back to Closed. If they fail, back to Open.

Libraries: Resilience4j (JVM), failsafe-go, polly (.NET), Envoy/Linkerd at the proxy layer (no application code).

## Bulkheads

Isolate independent dependencies so one slow one cannot starve the others:

- Separate HTTP client / connection pools per downstream.
- Separate thread pools (or async semaphores) for unrelated workloads.
- Separate queue partitions or consumer groups per source.

## What this all gives you

A service whose worst-case behaviour is bounded:

- A dependency outage causes localised, fast failures — not cascading timeouts.
- The service shed load when overwhelmed instead of falling over.
- Retries help with transient blips without making outages worse.

## References

- Constitution P6 (idempotency), P22 (timeouts), P23 (retries), P24 (circuit breakers)
- Google SRE book — chapter 22 "Addressing Cascading Failures"
- Release It! (Michael Nygard) — the canonical write-up of these patterns
