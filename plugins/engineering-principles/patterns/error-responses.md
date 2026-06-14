# Pattern — Error responses (RFC 9457 Problem Details)

Tech Stack T3 in practice. One error shape across all services so consumers parse once.

## Shape

Every error response is `Content-Type: application/problem+json` and follows [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html):

```json
{
  "type": "https://errors.example.com/invoicing/invoice-not-found",
  "title": "Invoice not found",
  "status": 404,
  "detail": "No invoice with id 'abc-123' exists in tenant 'acme'.",
  "instance": "/v2/customers/acme/invoices/abc-123",
  "trace_id": "01HW8KX9...",
  "errors": [
    { "field": "invoice_id", "code": "not_found" }
  ]
}
```

Required fields per RFC 9457: `type`, `title`, `status`, `detail`, `instance`.

Extension fields the chapter adds:

- `trace_id` — the OpenTelemetry trace id of the failing request. Lets support correlate to logs without asking the customer.
- `errors[]` — structured per-field validation errors for `400 Bad Request`. Each entry has `field` (JSON path) and `code` (machine-readable, e.g. `required`, `out_of_range`, `not_found`).

## `type` URIs

- Always an absolute URL. Convention: `https://errors.example.com/<service>/<slug>`.
- The URL doesn't have to resolve today, but it should be stable — consumers may match on it.
- Don't reuse a `type` for a different meaning. Coining a new one is cheap.

## HTTP status code mapping

| Situation | Status |
| --- | --- |
| Invalid input (shape, type, range) | `400` |
| Missing or invalid auth | `401` |
| Authed but not authorised | `403` |
| Resource doesn't exist | `404` |
| Method not allowed on resource | `405` |
| Idempotency-Key conflict (same key, different body) | `409` |
| Domain rule rejected the request | `422` |
| Rate-limited | `429` |
| Internal failure | `500` |
| Upstream dependency failed | `502` / `503` / `504` |

Don't use `200` with an error body. Use the right status.

## JVM example (Spring 6+ has `ProblemDetail`)

```java
@ControllerAdvice
class ApiExceptionHandler {

    @ExceptionHandler(InvoiceNotFound.class)
    ProblemDetail handleNotFound(InvoiceNotFound ex, HttpServletRequest req) {
        var pd = ProblemDetail.forStatusAndDetail(HttpStatus.NOT_FOUND, ex.getMessage());
        pd.setType(URI.create("https://errors.example.com/invoicing/invoice-not-found"));
        pd.setTitle("Invoice not found");
        pd.setInstance(URI.create(req.getRequestURI()));
        pd.setProperty("trace_id", Span.current().getSpanContext().getTraceId());
        return pd;
    }
}
```

`ProblemDetail` serialises with the correct content type and shape automatically.

## What this gives you

- Consumers parse one shape, not N. Stripe, GitHub, modern AWS APIs all use this pattern.
- Tracing is built in — every error response carries its `trace_id`.
- The OpenAPI spec documents the shape once (a shared `Problem` component).
- A breaking diff on the error contract is caught by CI per P11 + P12.

## References

- Constitution P11 (non-breaking changes), P12 (machine-readable contracts), P5 (observability)
- Tech Stack T3 (API design and specs)
- [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457.html)
