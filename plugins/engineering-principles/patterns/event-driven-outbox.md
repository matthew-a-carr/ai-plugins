# Pattern — Event-driven architecture with the outbox pattern

Make state changes and emitted events atomic. Constitution P10 in practice.

## Problem

A service writes to its database and emits a domain event. Two write targets, two failure modes:

- Commit DB, fail to publish → event lost.
- Publish, fail to commit DB → ghost event.

Distributed transactions across DB + broker are operationally toxic. Don't.

## Shape

1. In the same DB transaction as the state change, insert a row into an `outbox` table: `(id, aggregate_id, type, payload, created_at, dispatched_at NULL)`.
2. A separate dispatcher reads undispatched rows and publishes to the broker.
3. On successful publish, mark `dispatched_at`.
4. The broker is configured for at-least-once delivery. Consumers are idempotent (see `idempotency.md`).

## Event shape

- Use [CloudEvents](https://cloudevents.io/) (CNCF graduated) as the envelope: `id`, `source`, `type`, `time`, `subject`, `data`.
- The `id` doubles as the consumer's idempotency key.
- Schema the payload. Publish the schema alongside the AsyncAPI document (see Tech Stack T3 — API specs).

## When to use event-driven

- The producer doesn't need to know who the consumers are.
- The consumers don't need a synchronous answer.
- The domain has natural decoupling boundaries (e.g. "an order was placed" → fulfilment, billing, analytics each react independently).

## When NOT to use event-driven

- A synchronous answer is needed for the caller to proceed (use HTTP).
- The producer is one of N consumers in a tight loop (use direct method call).
- The boundary doesn't exist yet — don't introduce events to architect for hypothetical future consumers (Rule 2 — simplicity first).

## Outbox table DDL (Postgres)

```sql
CREATE TABLE outbox (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_id   VARCHAR(255) NOT NULL,
    event_type     VARCHAR(255) NOT NULL,
    payload        JSONB NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at   TIMESTAMPTZ
);

CREATE INDEX idx_outbox_unpublished ON outbox (created_at) WHERE published_at IS NULL;
```

The partial index on `published_at IS NULL` keeps the polling query fast — it only scans unpublished rows. Once `published_at` is set, the row drops out of the index.

Retention: delete or archive published rows on a schedule (e.g. older than 7 days). The outbox table is a transit buffer, not an event store.

## Polling publisher (Spring Boot)

```java
@Component
@RequiredArgsConstructor
class OutboxPublisher {
    private final JdbcTemplate jdbc;
    private final MessagePublisher publisher;

    @Scheduled(fixedDelay = 500)  // poll every 500ms
    @Transactional
    public void publishPending() {
        List<OutboxRow> rows = jdbc.query(
            """
            SELECT id, aggregate_type, aggregate_id, event_type, payload, created_at
            FROM outbox
            WHERE published_at IS NULL
            ORDER BY created_at
            LIMIT 100
            FOR UPDATE SKIP LOCKED
            """,
            (rs, i) -> new OutboxRow(
                rs.getObject("id", UUID.class),
                rs.getString("aggregate_type"),
                rs.getString("aggregate_id"),
                rs.getString("event_type"),
                rs.getString("payload"),
                rs.getTimestamp("created_at").toInstant()
            )
        );

        for (OutboxRow row : rows) {
            publisher.publish(row.toCloudEvent());
            jdbc.update("UPDATE outbox SET published_at = now() WHERE id = ?", row.id());
        }
    }
}
```

`FOR UPDATE SKIP LOCKED` prevents multiple instances from processing the same row. Each instance picks up a different batch.

Alternative to polling: Debezium CDC reads the Postgres WAL and publishes outbox rows as they are committed. Higher throughput, lower latency, but adds an operational dependency. Choose polling for simplicity first; move to CDC when polling latency or volume becomes a problem.

## CloudEvents envelope example

```json
{
  "specversion": "1.0",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "source": "/invoicing-service",
  "type": "com.example.invoicing.InvoiceCreated",
  "subject": "INV-2026-001",
  "time": "2026-06-14T13:00:00.000Z",
  "datacontenttype": "application/json",
  "data": {
    "invoiceId": "INV-2026-001",
    "customerId": "CUST-42",
    "totalAmount": 1500.00,
    "currency": "GBP",
    "lineItems": [
      { "description": "Consulting — June 2026", "amount": 1500.00 }
    ]
  }
}
```

Key fields:

- `id` — the outbox row ID. Consumers use this as the idempotency key (P6).
- `source` — identifies the producing service.
- `type` — reverse-DNS event type. Consumers filter on this.
- `subject` — the business identifier (invoice ID, order ID).
- `data` — the domain payload. Schema this and publish it alongside the AsyncAPI document (T3).

## References

- Constitution P6 (idempotency), P10 (event-driven where the domain warrants it)
- `sync-vs-async-decision.md` — when to use events vs sync HTTP
- `idempotency.md` — consumer deduplication
- CloudEvents: https://cloudevents.io/
- AsyncAPI: https://www.asyncapi.com/
