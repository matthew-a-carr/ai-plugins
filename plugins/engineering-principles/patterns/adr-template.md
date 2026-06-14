# Pattern — Architecture Decision Record (ADR) template

Anchors: P28, V3

Record decisions that shape the system so future engineers (and agents) know what was decided, why, and what was considered.

## Filename convention

Store ADRs in `docs/decisions/` using zero-padded sequence numbers:

```text
docs/decisions/
  0001-choose-message-broker.md
  0002-use-cursor-pagination.md
```

## Status values

`proposed` → `accepted` → optionally `superseded` or `deprecated`.

A superseded ADR links forward to the replacement. Never delete an ADR — the history matters.

## Template

```markdown
# NNNN — <Title>

**Status:** proposed | accepted | superseded by NNNN | deprecated
**Date:** YYYY-MM-DD
**Anchors cited:** P10, C10 (list the principle anchors this decision relates to)

## Context

What is the situation? What forces are at play? Keep it to 3–5 sentences.

## Decision

State the decision in one sentence. Then explain how it will be realised.

## Consequences

- What becomes easier or possible?
- What becomes harder or impossible?
- What operational burden does this add?

## Alternatives considered

| Option | Pros | Cons | Why rejected |
| --- | --- | --- | --- |
| Option A | ... | ... | ... |
| Option B | ... | ... | ... |
```

## Example — choosing a message broker

```markdown
# 0003 — Use SQS as the default message broker

**Status:** accepted
**Date:** 2026-03-12
**Anchors cited:** P10, C10, P28

## Context

The team needs an event broker for the outbox pattern (P10). The platform
already runs on AWS. The team has no Kafka operational experience.

## Decision

Use Amazon SQS with SNS fan-out as the default broker for chapter services.

## Consequences

- SQS is managed — no broker ops burden.
- No log-compaction or replay from an offset (unlike Kafka).
- Consumer ordering is best-effort (FIFO queues available per-partition).

## Alternatives considered

| Option | Pros | Cons | Why rejected |
| --- | --- | --- | --- |
| Apache Kafka (MSK) | Replay, ordering, CNCF ecosystem | Ops burden, cost at low volume | Team lacks Kafka experience; volume doesn't justify it |
| RabbitMQ | Lightweight, familiar | Not managed on this platform | Ops burden |
```

## When to write an ADR

- Novel technology choice not covered by the default stack (T1).
- Divergence from a Tier 3 default (the agent or reviewer flagged an unrecorded divergence).
- Recurring trade-off the team keeps debating — record it once.
- Significant new dependency (see `new-dependency-adr.md`).

## When NOT to write an ADR

- Following the default stack — the default is already documented in `tech-stack.md`.
- Trivial choices (library minor version, formatting rule) — these don't warrant the overhead.

## References

- Constitution P28 (simplicity — justify every moving part)
- Values V3 (prefer boring technology)
- [Documenting Architecture Decisions — Michael Nygard](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
