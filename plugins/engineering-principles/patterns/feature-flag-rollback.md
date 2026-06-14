# Pattern — Feature flag + documented rollback

Constitution P9 in practice.

## Default shape

- Flag default = `off`.
- Code path for `off` = existing behaviour. Code path for `on` = new behaviour.
- Flag scoped at the most granular boundary that's safe — customer, tenant, or per-user.

## PR checklist

- [ ] Flag defined in the flag system with a clear name (`<context>-<feature>-v1`).
- [ ] Both code paths tested.
- [ ] Default off.
- [ ] Rollback plan in the PR description: "Turn flag off; the previous path runs unchanged. State any data implications."
- [ ] Observability shows which path served each request (`served_by` label on metrics).

## Expiry

Flags expire. After full rollout + soak (≥30 days at 100%), the flag and the old path are removed in a follow-up PR. Quarterly cleanup pass on orphans.

## When NOT to flag

- Bug fixes with no behavioural surface change (a tightening of validation that rejects what was previously broken).
- Internal refactors invisible to consumers.

If unsure, flag. Flags are cheap; outages aren't.

## References

- Constitution P3 (small slices), P9 (feature flags + rollback)
