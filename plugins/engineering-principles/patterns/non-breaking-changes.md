# Pattern — Non-breaking changes (expand–contract)

Constitution P11 in practice. Every change to a contract — API, event schema, database column — ships without breaking existing consumers.

## Expand–contract for APIs

1. **Expand.** Add the new field / endpoint / version alongside the old. Old consumers see the old shape, new consumers see the new.
2. **Migrate.** Move consumers one at a time. Track adoption with a metric tagged by consumer or version.
3. **Contract.** Remove the old surface only after adoption is 100% and a soak period has elapsed (≥30 days at zero traffic on the old).

Never remove + add in the same release.

## Expand–contract for databases

1. **Expand.** Add the new column nullable, or the new table empty. Write to both old and new.
2. **Backfill.** Migrate historical data to the new shape.
3. **Switch reads.** Move read paths to the new shape behind a flag.
4. **Contract.** Stop writing the old, drop it.

Each step is a separate deployment. Each is independently rollback-safe.

## API versioning

- Prefer additive changes; reserve major-version bumps for genuinely incompatible shape changes.
- When versioning, use URI versioning (`/v2/...`) for clarity; header versioning is acceptable when consumers are internal and tooling supports it.
- The OpenAPI / AsyncAPI document is the contract. CI breaks on breaking diff.

## Event schemas

- Add fields, never rename or remove in a minor version.
- Consumers must tolerate unknown fields.
- Use schema registry (Confluent / Apicurio) with backward-compatibility enforcement in CI.

## References

- Constitution P3 (small slices), P9 (feature flags), P11 (non-breaking changes)
- `patterns/feature-flag-rollback.md` for the flagged rollout half
