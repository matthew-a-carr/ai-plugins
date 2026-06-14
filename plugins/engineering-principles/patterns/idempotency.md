# Pattern — Idempotency at boundaries

Constitution P6 in practice. Language- and framework-agnostic.

## Shape

For every external write or emitted side effect:

1. Derive a stable **idempotency key** from the natural-key fields of the operation (or accept a client-supplied one).
2. Persist `(key, result)` atomically with the side effect.
3. On retry, look up the key first — return the original result without re-doing the side effect.

## Pseudocode

```text
mint(key):
    existing = repo.find(key)
    if existing: return existing
    result = perform_side_effect()
    repo.save(key, result)   # in same transaction as side effect
    return result
```

- The key is a value object; equality is structural over all key fields.
- The persistence layer enforces a unique constraint on the key — the race lost by a concurrent writer is caught and re-read, not allowed to duplicate.
- Where the side effect is in an external system, use an outbox (see `event-driven-outbox.md`) so the persist + emit pair is atomic.

## Metrics

Emit at least:

```text
operation_minted_total{tenant, operation}
operation_reused_total{tenant, operation}
```

Steady-state reuse is high for retry-heavy workloads (≥99%). A drop in reuse rate is a signal something upstream changed.

## Tests

- Unit: same key → same result, no I/O.
- Integration (real DB): two concurrent attempts on the same key → one mint, one reuse, one row in the store.
- Contract (where the boundary is an external system): the consumer's retry produces the same result.

## References

- Constitution P6 (idempotency), P7 (test pyramid)
- `patterns/event-driven-outbox.md` for the persist+emit case
