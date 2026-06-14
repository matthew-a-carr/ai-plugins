# Pattern — Observability baseline

Anchors: P5, P26

What to add for every new use case so you can tell whether it is working. Ship observability in the same PR as the behaviour (P5).

## Minimum metric set

Every use case gets three metrics:

| Metric | Type | Name convention | Example |
| --- | --- | --- | --- |
| Latency | Histogram | `<domain>.<operation>.duration` | `invoicing.create_invoice.duration` |
| Request rate | Counter | `<domain>.<operation>.total` | `invoicing.create_invoice.total` |
| Error rate | Counter | `<domain>.<operation>.errors` | `invoicing.create_invoice.errors` |

Tag with at least: `outcome` (success / failure), `error_type` (timeout, validation, upstream, etc.).

## Structured log fields

Use structured JSON logging. Every log line includes:

| Field | Source | Example |
| --- | --- | --- |
| `correlationId` | Propagated from request header or generated | `550e8400-e29b-41d4-a716-446655440000` |
| `businessId` | The domain identifier relevant to this operation | `customerId=42`, `invoiceId=INV-2026-001` |
| `timestamp` | ISO 8601 with timezone | `2026-06-14T13:00:00.000Z` |
| `level` | Standard levels | `INFO`, `WARN`, `ERROR` |
| `message` | What happened, in plain language | `Invoice created` |
| `serviceName` | The service emitting the log | `invoicing-service` |

Do not log sensitive data (PII, credentials, tokens). Mask or omit.

## Trace span naming

```text
<service>.<layer>.<operation>
```

Examples: `invoicing.delivery.createInvoice`, `invoicing.infrastructure.saveInvoice`.

Keep spans short and specific. One span per significant I/O call (DB query, HTTP call, message publish). Do not wrap pure domain logic in spans — it adds noise without value.

## Spring Boot / Micrometer setup

```java
@Component
@RequiredArgsConstructor
class CreateInvoiceMetrics {
    private final MeterRegistry registry;

    private final Timer duration = Timer.builder("invoicing.create_invoice.duration")
        .description("Time to create an invoice")
        .register(registry);

    private final Counter total = Counter.builder("invoicing.create_invoice.total")
        .description("Invoices created")
        .register(registry);

    private final Counter errors = Counter.builder("invoicing.create_invoice.errors")
        .description("Invoice creation failures")
        .register(registry);

    public <T> T record(Supplier<T> operation) {
        try {
            T result = duration.record(operation);
            total.increment();
            return result;
        } catch (Exception e) {
            errors.increment();
            throw e;
        }
    }
}
```

For tracing, add `io.micrometer:micrometer-tracing-bridge-otel` and configure the exporter in `application.yml`:

```yaml
management:
  tracing:
    sampling:
      probability: 1.0   # 100% in dev/staging; tune down in production
  otlp:
    tracing:
      endpoint: http://otel-collector:4318/v1/traces
```

## Dashboard checklist

Before shipping the feature, verify the dashboard shows:

- [ ] Latency p50, p95, p99 for the use case
- [ ] Request rate (requests per second)
- [ ] Error rate and error ratio (errors / total)
- [ ] Alerts configured for error rate exceeding the SLO threshold (see P26)
- [ ] Correlation ID links from logs to traces

If the service has an SLO (P26), the dashboard should show the error budget burn rate.

## References

- Constitution P5 (observability first), P26 (SLOs)
- OpenTelemetry: https://opentelemetry.io/
- Micrometer: https://micrometer.io/
